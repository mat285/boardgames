package v1alpha1

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
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type TerminalPlayer struct {
	State        game.State
	NeedMove     bool
	NeedMoveLock sync.Mutex
	MoveChan     chan v1alpha1.Move
}

func NewTerminalPlayer() *TerminalPlayer {
	return &TerminalPlayer{
		NeedMove: false,
		MoveChan: make(chan v1alpha1.Move),
	}
}

func (p *TerminalPlayer) Run() {
	for {
		entry := Prompt(">")
		parts := strings.Split(entry, " ")
		cmd := parts[0]
		switch cmd {
		case "exit":
			os.Exit(0)
		case "board":
			if p.State.Players == nil {
				fmt.Println("No active state")
				continue
			}
			for _, card := range p.State.Board.AvailableCards() {
				fmt.Println(jsonMarshal(card))
			}
		case "hand":
			if p.State.Players == nil {
				fmt.Println("No active state")
				continue
			}
			hand := p.State.Players[p.State.CurrentIndex].Hand
			fmt.Println(prettyJSON(hand))
		case "gems":
			if p.State.Players == nil {
				fmt.Println("No active state")
				continue
			}
			hand := p.State.Players[p.State.CurrentIndex].Hand
			fmt.Println(prettyJSON(hand.Gems.Add(hand.CardCounts())))
		case "cards":
			if p.State.Players == nil {
				fmt.Println("No active state")
				continue
			}
			hand := p.State.Players[p.State.CurrentIndex].Hand
			fmt.Println(prettyJSON(hand.Cards))
		case "moves":
			if p.State.Players == nil {
				fmt.Println("No active state")
				continue
			}

			moves, err := p.State.ValidMoves()
			if err != nil {
				fmt.Println(err)
				continue
			}
			for i, move := range moves {
				err = p.PrintMove(i, move)
				if err != nil {
					fmt.Println(err)
				}
			}
			if !p.NeedMove {
				fmt.Println("not waiting for move")
				continue
			}
			num, err := p.PollForInput()
			if err != nil {
				fmt.Println(err)
				continue
			}
			p.MoveChan <- moves[num]

		default:
			fmt.Println("Commands: board hand gems cards moves exit")
		}
	}

}

func (p *TerminalPlayer) RequestMove(ctx context.Context, req v1alpha1.MoveRequest) (v1alpha1.Move, error) {
	typed, ok := req.State.(game.State)
	if !ok {
		return nil, fmt.Errorf("Wrong game")
	}
	p.NeedMoveLock.Lock()
	p.State = typed
	p.NeedMove = true
	p.NeedMoveLock.Unlock()

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
	fmt.Println(i, string(data))
	return nil
}

func (p *TerminalPlayer) PollForInput() (int, error) {
	str := Prompt("Choose a move number:")
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

// Prompt gives a prompt and reads input until newlines.
func Prompt(prompt string) string {
	return PromptFrom(os.Stdout, os.Stdin, prompt)
}

// Promptf gives a prompt of a given format and args and reads input until newlines.
func Promptf(format string, args ...interface{}) string {
	return PromptFrom(os.Stdout, os.Stdin, fmt.Sprintf(format, args...))
}

// PromptFrom gives a prompt and reads input until newlines from a given set of streams.
func PromptFrom(stdout io.Writer, stdin io.Reader, prompt string) string {
	fmt.Fprint(stdout, prompt)

	scanner := bufio.NewScanner(stdin)
	var output string
	if scanner.Scan() {
		output = scanner.Text()
	}
	return output
}
