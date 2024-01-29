package terminal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor/pkg/game"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
	"github.com/mat285/boardgames/games/splendor/serializer"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var _ engine.Player = new(TerminalPlayer)

type TerminalPlayer struct {
	engine.ConnectedPlayer

	Connection *engine.Connection

	serializer.Get

	State        game.State
	NeedMove     bool
	NeedMoveLock sync.Mutex
	MoveChan     chan v1alpha1.Move

	LogMessages chan string

	logPause bool
	logLock  sync.Mutex
}

func NewTerminalPlayer(g v1alpha1.Game, server connection.Server) *TerminalPlayer {
	tp := &TerminalPlayer{
		NeedMove:    false,
		MoveChan:    make(chan v1alpha1.Move),
		LogMessages: make(chan string, 100),
	}
	tp.ConnectedPlayer = engine.NewConnectedPlayer(uuid.V4(), "user", g, server)
	conn, _ := server.Connect(context.Background())
	tp.Connection = engine.NewConnection(g, conn)
	return tp
}

func (p *TerminalPlayer) Connect(ctx context.Context) error {
	go p.ConnectedPlayer.Receive(ctx, p.Handle)
	go p.Run(ctx)
	return nil
}

func (p *TerminalPlayer) Run(ctx context.Context) error {
	go p.writeLogs(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		entry := p.Prompt("\n>")
		parts := strings.Split(entry, " ")
		cmd := parts[0]
		switch cmd {
		case "":
			continue
		case "exit":
			os.Exit(0)
		case "board":
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}
			for _, card := range p.State.Board.AvailableCards() {
				p.Println(jsonMarshal(card))
			}
		case "players":
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}
			for i, player := range p.State.Players {
				p.Println(i, player.ID, jsonMarshal(player.Hand))
			}
		case "player":
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}
			if len(parts) != 2 {
				p.Println("Need player index")
				continue
			}
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				p.Println(err)
				continue
			}
			player := p.State.Players[num]
			p.Println(prettyJSON(player.Hand))
		case "hand":
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}
			hand := p.State.Players[p.State.CurrentIndex].Hand
			p.Println(prettyJSON(hand))
		case "gems":
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}
			hand := p.State.Players[p.State.CurrentIndex].Hand
			p.Println(prettyJSON(hand.Gems.Add(hand.CardCounts())))
		case "cards":
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}
			hand := p.State.Players[p.State.CurrentIndex].Hand
			p.Println(prettyJSON(hand.Cards))
		case "reserve", "purchase":
			if len(parts) != 2 {
				p.Println("Need card id")
				continue
			}
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}
			if !p.NeedMove {
				p.Println("not waiting for move")
				continue
			}
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				p.Println(err)
				continue
			}
			card := items.Cards()[num]
			cm := &game.CardMove{Card: card}
			mv := &game.Move{}
			if cmd == "reserve" {
				mv.Reserve = cm
			} else {
				mv.Purchase = cm
			}

			p.MoveChan <- mv
		case "collect":
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}
			if !p.NeedMove {
				p.Println("not waiting for move")
				continue
			}
			take, ret, err := parseGems(parts[1:])
			if err != nil {
				p.Println(err)
				continue
			}
			if take.Wild > 0 {
				p.Println("Cannot take wilds")
				continue
			}

			p.MoveChan <- &game.Move{Collect: &game.CollectMove{Take: take, Return: ret}}
		case "moves":
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}

			moves, err := p.State.ValidMoves()
			if err != nil {
				p.Println(err)
				continue
			}
			for i, move := range moves {
				err = p.PrintMove(i, move)
				if err != nil {
					p.Println(err)
				}
			}
			if !p.NeedMove {
				p.Println("not waiting for move")
				continue
			}
			num, err := p.PollForInput()
			if err != nil {
				p.Println(err)
				continue
			}
			p.MoveChan <- moves[num]

		default:
			p.Println("Commands: board hand gems cards moves exit\n")
		}
	}

}

func (p *TerminalPlayer) Handle(ctx context.Context, message v1alpha1.Message) (*v1alpha1.Message, error) {
	switch message.Type {
	case v1alpha1.MessageTypeRequestMove:
		so, err := message.DeserializeToObject()
		if err != nil {
			return nil, err
		}
		state, err := p.Serializer().DeserializeState(so)
		if err != nil {
			return nil, err
		}
		return p.Request(ctx, state)
	}
	p.Println(fmt.Sprintf("\nMessage from game server:%s\n", message))
	return nil, nil
}

func (p *TerminalPlayer) Request(ctx context.Context, state v1alpha1.StateData) (*v1alpha1.Message, error) {
	typed, ok := state.(game.State)
	if !ok {
		return nil, fmt.Errorf("Wrong game")
	}
	p.NeedMoveLock.Lock()
	p.State = typed
	p.NeedMove = true
	p.NeedMoveLock.Unlock()

	p.Println("\nRequest for next move from game server!\n")

	move := <-p.MoveChan
	p.NeedMoveLock.Lock()
	p.NeedMove = true
	p.NeedMoveLock.Unlock()
	data, err := p.Serializer().SerializeMove(move)
	if err != nil {
		return nil, err
	}
	return v1alpha1.NewMessage(v1alpha1.MessageTypePlayerMove, data)
}

func (p *TerminalPlayer) PrintMove(i int, move v1alpha1.Move) error {
	data, err := json.Marshal(move)
	if err != nil {
		return err
	}
	p.Println(i, string(data))
	return nil
}

func (p *TerminalPlayer) Println(vs ...interface{}) {
	if len(vs) == 0 {
		return
	}
	ss := make([]string, len(vs))
	for i := range vs {
		ss[i] = fmt.Sprintf("%v", vs[i])
	}
	p.LogMessages <- strings.Join(ss, " ") + "\n"
}

func (p *TerminalPlayer) Printf(format string, args ...interface{}) {
	p.LogMessages <- fmt.Sprintf(format, args...)
}

func (p *TerminalPlayer) drain() {
	for len(p.LogMessages) > 0 {
		msg := <-p.LogMessages
		fmt.Printf(msg)
	}
}

func (p *TerminalPlayer) pauseLogs() {
	if p.logPause {
		return
	}
	p.logLock.Lock()
	p.logPause = true
}

func (p *TerminalPlayer) resumeLogs() {
	if !p.logPause {
		return
	}
	p.logPause = false
	p.logLock.Unlock()
}

func (p *TerminalPlayer) writeLogs(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-p.LogMessages:
			if p.logPause {
				p.logLock.Lock()
				fmt.Printf(msg)
			} else {
				fmt.Printf(msg)
			}
		}
	}
}

func (p *TerminalPlayer) PollForInput() (int, error) {
	str := p.Prompt("Choose a move number:")
	return strconv.Atoi(str)
}

func prettyJSON(i interface{}) string {
	bytes, _ := json.MarshalIndent(i, "", "  ")
	return string(bytes)
}
func jsonMarshal(i interface{}) string {
	bytes, _ := json.Marshal(i)
	return string(bytes)
}

func parseGems(parts []string) (items.GemCount, items.GemCount, error) {
	take := items.GemCount{}
	returning := false
	count := 0
	gems := items.GemCount{}
	for _, part := range parts {
		if part == "return" {
			returning = true
			if count < 2 || count > 3 {
				return items.GemCount{}, items.GemCount{}, fmt.Errorf("Wrong number of gems")
			}
			take = gems
			gems = items.GemCount{}
			count = 0
			continue
		}
		gem, err := items.ParseGem(part)
		if err != nil {
			return items.GemCount{}, items.GemCount{}, err
		}
		gems = gems.AddGem(gem, 1)
		count++
	}
	if returning {
		return take, gems, nil
	}

	if count < 2 || count > 3 {
		return items.GemCount{}, items.GemCount{}, fmt.Errorf("Wrong number of gems")
	}

	return gems, items.GemCount{}, nil
}

// Prompt gives a prompt and reads input until newlines.
func (p *TerminalPlayer) Prompt(prompt string) string {
	return p.PromptFrom(os.Stdout, os.Stdin, prompt)
}

// PromptFrom gives a prompt and reads input until newlines from a given set of streams.
func (p *TerminalPlayer) PromptFrom(stdout io.Writer, stdin io.Reader, prompt string) string {
	// p.pauseLogs()
	// p.drain()
	// defer p.resumeLogs()
	fmt.Printf(prompt)

	scanner := bufio.NewScanner(stdin)
	var output string
	if scanner.Scan() {
		output = scanner.Text()
	}
	return output
}
