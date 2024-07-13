package terminal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/blend/go-sdk/uuid"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	splendorgame "github.com/mat285/boardgames/games/splendor"
	splendor "github.com/mat285/boardgames/games/splendor/pkg/game"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
	client "github.com/mat285/boardgames/pkg/client/core/v1alpha1"
	httpclient "github.com/mat285/boardgames/pkg/client/http/v1alpha1"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	messages "github.com/mat285/boardgames/pkg/messages/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

var _ connection.Client = new(Terminal)

type (
	errMsg error
)

type Model struct {
	title       viewport.Model
	view        viewport.Model
	history     []string
	index       int
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error

	Terminal *Terminal
}

func (p *Terminal) InitialModel() Model {
	ta := textarea.New()
	ta.Placeholder = "Enter a command..."
	ta.Focus()

	ta.Prompt = "> "
	ta.CharLimit = 280

	ta.SetWidth(90)
	ta.SetHeight(1)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	title := viewport.New(30, 2)
	title.SetContent("Pickl.io Terminal Client\n")

	view := viewport.New(90, 30)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Model{
		Terminal:    p,
		textarea:    ta,
		history:     []string{},
		title:       title,
		view:        view,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (p *Terminal) Start() error {
	go p.retryListen(p.Ctx)
	_, err := tea.NewProgram(p.InitialModel()).Run()
	return err
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	// m.view, vpCmd = m.view.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			// fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyUp:
			prevIdx := len(m.history) - m.index - 1
			if prevIdx >= 0 {
				m.index++
				m.textarea.SetValue(m.history[prevIdx])
			}
		case tea.KeyDown:
			nextIdx := len(m.history) - m.index
			if nextIdx < len(m.history) {
				m.index--
				m.textarea.SetValue(m.history[nextIdx])
			}
			if nextIdx == len(m.history) {
				m.index = 0
				m.textarea.SetValue("")
			}
		case tea.KeyEnter:
			message := m.textarea.Value()
			m.history = append(m.history, message)
			result := m.Terminal.HandleMessage(m.Terminal.Ctx, message)
			m.view.SetContent(result)
			m.textarea.Reset()
			m.view.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		m.title.View(),
		m.view.View(),
		m.textarea.View(),
	) + "\n\n"
}

type Terminal struct {
	Ctx context.Context

	*client.Player

	apiClient *httpclient.Client
	gid       uuid.UUID

	State        splendor.State
	NeedMove     bool
	NeedMoveLock sync.Mutex
	MoveChan     chan game.Move

	LogMessages chan string

	logPause bool
	logLock  sync.Mutex
}

func NewTerminal(ctx context.Context, username string, cli *httpclient.Client) *Terminal {
	tp := &Terminal{
		Ctx:         ctx,
		Player:      client.NewPlayer(username, nil, cli),
		apiClient:   cli,
		NeedMove:    false,
		MoveChan:    make(chan game.Move),
		LogMessages: make(chan string, 100),
	}
	return tp
}

func (p *Terminal) HandleMessage(ctx context.Context, entry string) (result string) {
	parts := strings.Split(entry, " ")
	cmd := parts[0]
	switch cmd {
	case "":
		return
	case "exit":
		os.Exit(0)
	case "newtest":
		if p.apiClient == nil {
			result += fmt.Sprintln("not allowed without api client")
			return
		}
		p.apiClient.Username = p.Username
		gcfg := splendor.Config{VictoryPoints: 7}
		gid, err := p.apiClient.NewGame(ctx, "splendor", gcfg)
		if err != nil {
			result += fmt.Sprintln("Error making new game", err.Error)
			return
		}
		result += fmt.Sprintln("Created new game with ID:", gid)
		p.gid = gid
		p.Game = splendorgame.NewGameWithConfig(gcfg)
		p.Player.Game = p.Game
		p.Player.Message = messages.NewProvider(p.Game)
		err = p.apiClient.Join(ctx, gid)
		if err != nil {
			result += fmt.Sprintln("Error joining game", gid, err.Error)
			return
		}
		// p.Player.Game =
		result += fmt.Sprintln("Joined game", gid)
		err = p.apiClient.Start(ctx, gid)
		if err != nil {
			result += fmt.Sprintln("Error starting game", gid, err.Error)
			return
		}

		result += fmt.Sprintln("Started game", gid)
		return

	case "new":
		if p.apiClient == nil {
			result += fmt.Sprintln("not allowed without api client")
			return
		}
		p.apiClient.Username = p.Username
		gcfg := splendor.Config{VictoryPoints: 7}
		gid, err := p.apiClient.NewGame(ctx, "splendor", gcfg)
		if err != nil {
			result += fmt.Sprintln("Error making new game", err.Error)
			return
		}
		result += fmt.Sprintln("Created new game with ID:", gid)
		return
	case "join":
		if p.apiClient == nil {
			result += fmt.Sprintln("not allowed without api client")
			return
		}
		if len(parts) < 2 {
			result += "Need game id\n"
			return
		}
		gid, err := uuid.Parse(parts[1])
		if err != nil {
			result += err.Error() + "\n"
			return
		}
		p.gid = gid
		gcfg := splendor.Config{VictoryPoints: 7}
		p.Game = splendorgame.NewGameWithConfig(gcfg)
		p.Player.Game = p.Game
		p.Player.Message = messages.NewProvider(p.Game)
		err = p.apiClient.Join(ctx, gid)
		if err != nil {
			result += fmt.Sprintln("Error joining game", gid, err.Error)
			return
		}
		// p.Player.Game =
		result += fmt.Sprintln("Joined game", gid)
		return
	case "start":
		if p.apiClient == nil {
			result += fmt.Sprintln("not allowed without api client")
			return
		}
		if p.gid.IsZero() {
			result += fmt.Sprintln("Need to join game first")
			return
		}
		err := p.apiClient.Start(ctx, p.gid)
		if err != nil {
			result += fmt.Sprintln("Error starting game", p.gid, err.Error)
			return
		}

		result += fmt.Sprintln("Started game", p.gid)
		return
	case "board":
		if err := p.maybeFetchState(ctx); err != nil {
			result += fmt.Sprintln("error fetching state", err)
			return
		}
		if p.State.Players == nil {
			result += fmt.Sprintln("No active state")
			return
		}
		for _, card := range p.State.Board.AvailableCards() {
			result += fmt.Sprintln(jsonMarshal(card))
		}
	case "players":
		if err := p.maybeFetchState(ctx); err != nil {
			result += fmt.Sprintln("error fetching state", err)
			return
		}
		if p.State.Players == nil {
			result += fmt.Sprintln("No active state")
			return
		}
		for i, player := range p.State.Players {
			result += fmt.Sprintln(i, player.ID, jsonMarshal(player.Hand))
		}
	case "player":
		if err := p.maybeFetchState(ctx); err != nil {
			result += fmt.Sprintln("error fetching state", err)
			return
		}
		if p.State.Players == nil {
			result += fmt.Sprintln("No active state")
			return
		}
		if len(parts) != 2 {
			result += fmt.Sprintln("Need player index")
			return
		}
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			result += fmt.Sprintln(err)
			return
		}
		player := p.State.Players[num]
		result += fmt.Sprintln(prettyJSON(player.Hand))
	case "hand":
		if err := p.maybeFetchState(ctx); err != nil {
			result += fmt.Sprintln("error fetching state", err)
			return
		}
		if p.State.Players == nil {
			result += fmt.Sprintln("No active state")
			return
		}
		player, err := p.State.GetCurrentPlayer()
		if err != nil {
			result += fmt.Sprintln(err)
		}
		hand := player.Hand
		result += fmt.Sprintln(prettyJSON(hand))
	case "gems":
		if err := p.maybeFetchState(ctx); err != nil {
			result += fmt.Sprintln("error fetching state", err)
			return
		}
		if p.State.Players == nil {
			result += fmt.Sprintln("No active state")
			return
		}
		player, err := p.State.GetCurrentPlayer()
		if err != nil {
			result += fmt.Sprintln(err)
		}
		hand := player.Hand
		result += fmt.Sprintln(prettyJSON(hand.Gems.Add(hand.CardCounts())))
	case "cards":
		if err := p.maybeFetchState(ctx); err != nil {
			result += fmt.Sprintln("error fetching state", err)
			return
		}
		if p.State.Players == nil {
			result += fmt.Sprintln("No active state")
			return
		}
		player, err := p.State.GetCurrentPlayer()
		if err != nil {
			result += fmt.Sprintln(err)
		}
		hand := player.Hand
		result += fmt.Sprintln(prettyJSON(hand.Cards))
	case "reserve", "purchase":
		if len(parts) != 2 {
			result += fmt.Sprintln("Need card id")
			return
		}
		if err := p.maybeFetchState(ctx); err != nil {
			result += fmt.Sprintln("error fetching state", err)
			return
		}
		if p.State.Players == nil {
			result += fmt.Sprintln("No active state")
			return
		}
		if !p.NeedMove {
			result += fmt.Sprintln("not waiting for move")
			return
		}
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			result += fmt.Sprintln(err)
			return
		}
		card := items.Cards()[num]
		cm := &splendor.CardMove{Card: card}
		mv := &splendor.Move{}
		if cmd == "reserve" {
			mv.Reserve = cm
		} else {
			mv.Purchase = cm
		}

		made, err := p.maybeSendMove(ctx, mv)
		if err != nil {
			result += fmt.Sprintln("Error sending move", err)
			return
		}
		if made {
			result += fmt.Sprintln("Made move")
			return
		}
		p.MoveChan <- mv

	case "collect":
		if err := p.maybeFetchState(ctx); err != nil {
			result += fmt.Sprintln("error fetching state", err)
			return
		}
		if p.State.Players == nil {
			result += fmt.Sprintln("No active state")
			return
		}
		if !p.NeedMove {
			result += fmt.Sprintln("not waiting for move")
			return
		}
		take, ret, err := parseGems(parts[1:])
		if err != nil {
			result += fmt.Sprintln(err)
			return
		}
		if take.Wild > 0 {
			result += fmt.Sprintln("Cannot take wilds")
			return
		}
		move := &splendor.Move{Collect: &splendor.CollectMove{Take: take, Return: ret}}

		made, err := p.maybeSendMove(ctx, move)
		if err != nil {
			result += fmt.Sprintln("Error sending move", err)
			return
		}
		if made {
			result += fmt.Sprintln("Made move")
			return
		}
		p.MoveChan <- move
		return
	case "moves":
		if err := p.maybeFetchState(ctx); err != nil {
			result += fmt.Sprintln("error fetching state", err)
			return
		}
		if p.State.Players == nil {
			result += fmt.Sprintln("No active state")
			return
		}

		moves, err := p.State.ValidMoves()
		if err != nil {
			result += fmt.Sprintln(err)
			return
		}
		for i, move := range moves {
			result += moveString(i, move)
		}
		if !p.NeedMove {
			result += fmt.Sprintln("not waiting for move")
			return
		}
	case "move":
		if err := p.maybeFetchState(ctx); err != nil {
			result += fmt.Sprintln("error fetching state", err)
			return
		}
		if p.State.Players == nil {
			result += fmt.Sprintln("No active state")
			return
		}
		if len(parts) < 2 {
			result += "Need move number\n"
			return
		}
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			result += fmt.Sprintln(err)
			return
		}
		moves, err := p.State.ValidMoves()
		if err != nil {
			result += fmt.Sprintln(err)
			return
		}

		move := moves[num]
		made, err := p.maybeSendMove(ctx, move)
		if err != nil {
			result += fmt.Sprintln("Error sending move", err)
			return
		}
		if made {
			result += fmt.Sprintln("Made move")
			return
		}
		p.MoveChan <- move

	default:
		result += fmt.Sprintln("Commands: board hand gems cards moves exit\n")
	}
	return

}

func (p *Terminal) maybeFetchState(ctx context.Context) error {
	if p.apiClient == nil {
		return nil
	}
	p.NeedMoveLock.Lock()
	defer p.NeedMoveLock.Unlock()
	packet, err := p.apiClient.GetState(ctx, p.gid)
	if err != nil {
		fmt.Println(string(packet.Payload))
		return err
	}
	// fmt.Println("fetch state packet", string(packet.MustJSON()))
	state, err := p.Game.DeserializeState(&game.SerializedObject{
		ID:   p.gid,
		Data: packet.Payload,
	})
	// state, err := p.Message.ExtractState(*packet)
	if err != nil {
		return err
	}
	typed, ok := state.(splendor.State)
	if !ok {
		return fmt.Errorf("Wrong game")
	}
	p.State = typed
	cid, _ := p.State.CurrentPlayer()
	if cid.Equal(p.apiClient.UserID) {
		p.NeedMove = true
	}
	return nil
}

func (p *Terminal) maybeSendMove(ctx context.Context, move game.Move) (bool, error) {
	if p.apiClient == nil {
		return false, nil
	}
	packet, err := p.Message.MessagePlayerMove(move, p.apiClient.UserID)
	if err != nil {
		return false, err
	}
	_, err = p.apiClient.MakeMove(ctx, p.gid, p.apiClient.UserID, *packet)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *Terminal) retryListen(ctx context.Context) {
	for {
		player := p.Player
		if player == nil {
			time.Sleep(time.Second * 2)
			continue
		}
		player.Listen(ctx, p.Handle)
	}
}

func (p *Terminal) Handle(ctx context.Context, packet wire.Packet) error {
	fmt.Println("got inbound packet from server")
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
	// p.Println(fmt.Sprintf("\nMessage from game server:%s\n", string(packet.Payload)))
	return nil
}

func (p *Terminal) Request(ctx context.Context, state game.StateData, req uuid.UUID) error {
	typed, ok := state.(splendor.State)
	if !ok {
		return fmt.Errorf("Wrong game")
	}
	p.NeedMoveLock.Lock()
	p.State = typed
	p.NeedMove = true
	p.NeedMoveLock.Unlock()

	// p.Println("\nRequest for next move from game server!\n")

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

func moveString(i int, move game.Move) string {
	data, err := json.Marshal(move)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintln(i, string(data))
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
