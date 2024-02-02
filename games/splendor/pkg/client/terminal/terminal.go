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
	splendor "github.com/mat285/boardgames/games/splendor/pkg/game"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
	client "github.com/mat285/boardgames/pkg/client/v1alpha1"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	messages "github.com/mat285/boardgames/pkg/messages/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

var _ connection.Client = new(TerminalPlayer)

type TerminalPlayer struct {
	*client.Player

	State        splendor.State
	NeedMove     bool
	NeedMoveLock sync.Mutex
	MoveChan     chan game.Move

	LogMessages chan string

	logPause bool
	logLock  sync.Mutex
}

func NewTerminalPlayer(username string, g game.Game, conn connection.Client) *TerminalPlayer {
	tp := &TerminalPlayer{
		Player:      client.NewPlayer(username, g, conn),
		NeedMove:    false,
		MoveChan:    make(chan game.Move),
		LogMessages: make(chan string, 100),
	}
	return tp
}

func (p *TerminalPlayer) Start(ctx context.Context) error {
	go p.Client.Listen(ctx, p.Handle)
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
			player, err := p.State.GetCurrentPlayer()
			if err != nil {
				p.Println(err)
			}
			hand := player.Hand
			p.Println(prettyJSON(hand))
		case "gems":
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}
			player, err := p.State.GetCurrentPlayer()
			if err != nil {
				p.Println(err)
			}
			hand := player.Hand
			p.Println(prettyJSON(hand.Gems.Add(hand.CardCounts())))
		case "cards":
			if p.State.Players == nil {
				p.Println("No active state")
				continue
			}
			player, err := p.State.GetCurrentPlayer()
			if err != nil {
				p.Println(err)
			}
			hand := player.Hand
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
			cm := &splendor.CardMove{Card: card}
			mv := &splendor.Move{}
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

			p.MoveChan <- &splendor.Move{Collect: &splendor.CollectMove{Take: take, Return: ret}}
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

func (p *TerminalPlayer) Handle(ctx context.Context, packet wire.Packet) error {
	switch packet.Type {
	case messages.PacketTypeRequestMove:
		state, err := p.Message.ExtractState(packet)
		if err != nil {
			return err
		}
		go p.Request(ctx, state, packet.ID)
		return nil
	case messages.PacketTypePlayerMoveInfo:
		var so game.SerializedObject
		json.Unmarshal(packet.Payload, &so)
		fmt.Println(packet.Type, string(so.Data))
		info, err := p.Message.ExtractPlayerMoveInfo(packet)
		if err != nil {
			return err
		}
		fmt.Println("Got player move info", string(info.Move.Data))
	}
	p.Println(fmt.Sprintf("\nMessage from game server:%s\n", string(packet.Payload)))
	return nil
}

func (p *TerminalPlayer) Request(ctx context.Context, state game.StateData, req uuid.UUID) error {
	typed, ok := state.(splendor.State)
	if !ok {
		return fmt.Errorf("Wrong game")
	}
	p.NeedMoveLock.Lock()
	p.State = typed
	p.NeedMove = true
	p.NeedMoveLock.Unlock()

	p.Println("\nRequest for next move from game server!\n")

	move := <-p.MoveChan
	p.NeedMoveLock.Lock()
	p.NeedMove = false
	p.NeedMoveLock.Unlock()
	packet, err := p.Message.MessagePlayerMove(move, req)
	if err != nil {
		return err
	}
	return p.Client.Send(ctx, *packet)
}

func (p *TerminalPlayer) PrintMove(i int, move game.Move) error {
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
	msg := strings.Join(ss, " ") + "\n"
	fmt.Printf(msg)
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
