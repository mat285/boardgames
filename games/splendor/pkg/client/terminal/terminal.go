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

	"github.com/mat285/boardgames/games/splendor/pkg/game"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type TerminalPlayer struct {
	State        game.State
	NeedMove     bool
	NeedMoveLock sync.Mutex
	MoveChan     chan v1alpha1.Move

	LogMessages chan string
}

func NewTerminalPlayer() *TerminalPlayer {
	return &TerminalPlayer{
		NeedMove:    false,
		MoveChan:    make(chan v1alpha1.Move),
		LogMessages: make(chan string, 5),
	}
}

func (p *TerminalPlayer) Run(ctx context.Context) error {
	go p.writeLogs(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		entry := p.Prompt(">")
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
			p.Println("Commands: board hand gems cards moves exit")
		}
	}

}

func (p *TerminalPlayer) Accept(ctx context.Context, message v1alpha1.Message) error {
	p.Println(fmt.Sprintf("\nMessage from game server:%s\n", message))
	return nil
}

func (p *TerminalPlayer) Request(ctx context.Context, req v1alpha1.MoveRequest) (v1alpha1.Move, error) {
	typed, ok := req.State.(game.State)
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
	return move, nil
}

func (p *TerminalPlayer) PrintMove(i int, move v1alpha1.Move) error {
	data, err := move.Serialize()
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

func (p *TerminalPlayer) writeLogs(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-p.LogMessages:
			fmt.Printf(msg)
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
	p.Printf(prompt)

	scanner := bufio.NewScanner(stdin)
	var output string
	if scanner.Scan() {
		output = scanner.Text()
	}
	return output
}
