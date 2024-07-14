package terminal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/uuid"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	splendorgame "github.com/mat285/boardgames/games/splendor"
	splendorhttpclient "github.com/mat285/boardgames/games/splendor/pkg/client/http"
	splendor "github.com/mat285/boardgames/games/splendor/pkg/game"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
	httpclient "github.com/mat285/boardgames/pkg/client/http/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	messages "github.com/mat285/boardgames/pkg/messages/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

// var _ connection.Client = new(Terminal)

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

	ta.SetWidth(140)
	ta.SetHeight(1)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	title := viewport.New(30, 2)
	title.SetContent("Splendor Terminal Client\n")

	view := viewport.New(140, 35)

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
			result = lipgloss.NewStyle().Width(120).Render(result)
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

	Game    game.Game
	Message messages.Provider

	SplendorClient *splendorhttpclient.Client
	CurrentGame    uuid.UUID

	State   splendor.State
	Packets chan wire.Packet
}

func NewTerminal(ctx context.Context, username string, cli *httpclient.Client) *Terminal {
	g := splendorgame.NewGameWithConfig(splendor.Config{VictoryPoints: 7})
	tp := &Terminal{
		Ctx:            ctx,
		Game:           g,
		Message:        messages.NewProvider(g),
		SplendorClient: splendorhttpclient.New(cli),
		Packets:        make(chan wire.Packet, 10),
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
	case "login":
		if len(parts) < 2 {
			result += fmt.Sprintln("need usernam")
			return
		}
		p.SplendorClient.Username = parts[1]
		p.SplendorClient.UserID = nil
		err := p.SplendorClient.Login(ctx)
		if err != nil {
			result += fmt.Sprintln("error logging in", err)
			return
		}
		err = p.SplendorClient.Connect(ctx, nil)
		if err != nil {
			result += fmt.Sprintln("error starting websocket", err)
			return
		}
		go p.retryListen(ctx)
		result += fmt.Sprintln("Successfully logged in as", p.SplendorClient.Username)
		return
	case "mine":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
			return
		}
		games, err := p.SplendorClient.GetUserGames(ctx)
		if err != nil {
			result += fmt.Sprintln("error getting user games", err)
			return
		}
		for _, g := range games {
			result += fmt.Sprintln("ID:", g.ID, "Game:", g.Game)
		}
		return
	case "new":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
			return
		}
		p.SplendorClient.Username = p.SplendorClient.Username
		gcfg := splendor.Config{VictoryPoints: 7}
		gid, err := p.SplendorClient.NewGame(ctx, "splendor", gcfg)
		if err != nil {
			result += fmt.Sprintln("Error making new game", err.Error)
			return
		}
		result += fmt.Sprintln("Created new game with ID:", gid)
		p.CurrentGame = gid
		err = p.SplendorClient.Join(ctx, gid)
		if err != nil {
			result += fmt.Sprintln("Error joining game", gid, err.Error)
			return
		}
		result += fmt.Sprintln("Joined game", gid)
		return
	case "join":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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
		p.CurrentGame = gid
		err = p.SplendorClient.Join(ctx, gid)
		if err != nil {
			result += fmt.Sprintln("Error joining game", gid, err.Error)
			return
		}
		result += fmt.Sprintln("Joined game", gid)
		return
	case "play":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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
		p.CurrentGame = gid
		result += fmt.Sprintln("switched to game room", gid)
		return
	case "start":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
			return
		}
		gid := p.CurrentGame
		if len(parts) > 1 {
			parsedID, err := uuid.Parse(parts[1])
			if err != nil {
				result += err.Error() + "\n"
				return
			}
			gid = parsedID
		}
		if gid.IsZero() {
			result += fmt.Sprintln("Need valid game id or current game")
			return
		}

		err := p.SplendorClient.Start(ctx, gid)
		if err != nil {
			result += fmt.Sprintln("Error starting game", gid, err.Error)
			return
		}

		result += fmt.Sprintln("Started game", gid)
		return
	case "unread":
		result += fmt.Sprintln("Currently have", len(p.Packets), "unread packets")
		return
	case "read":
		l := len(p.Packets)
		for i := 0; i < l; i++ {
			result += fmt.Sprintln("Packet", i, jsonMarshal(<-p.Packets))
		}
		return
	case "board":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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
		for _, card := range p.State.Board.AvailableCards() {
			result += fmt.Sprintln(jsonMarshal(card))
		}
		return
	case "players":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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
		for i, player := range p.State.Players {
			result += fmt.Sprintln(i, player.Username, player.ID, jsonMarshal(player.Hand))
		}
		return
	case "player":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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
		return
	case "hand":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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
		player, err := p.State.GetCurrentPlayer()
		if err != nil {
			result += fmt.Sprintln(err)
		}
		hand := player.Hand
		result += fmt.Sprintln(prettyJSON(hand))
		return
	case "gems":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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
		player, err := p.State.GetCurrentPlayer()
		if err != nil {
			result += fmt.Sprintln(err)
		}
		hand := player.Hand
		result += fmt.Sprintln(prettyJSON(hand.Gems.Add(hand.CardCounts())))
		return
	case "cards":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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
		player, err := p.State.GetCurrentPlayer()
		if err != nil {
			result += fmt.Sprintln(err)
		}
		hand := player.Hand
		result += fmt.Sprintln(prettyJSON(hand.Cards))
		return
	case "reserve", "purchase", "buy":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
			return
		}
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
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			result += fmt.Sprintln(err)
			return
		}
		card := items.Cards()[num]
		cm := &splendor.CardMove{Card: card}
		move := &splendor.Move{}
		if cmd == "reserve" {
			move.Reserve = cm
		} else {
			move.Purchase = cm
		}

		err = p.sendMove(ctx, move)
		if err != nil {
			result += fmt.Sprintln("Error sending move", err)
			return
		}

		result += fmt.Sprintln("Made move\n", prettyJSON(move))
		return

	case "collect", "take":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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

		err = p.sendMove(ctx, move)
		if err != nil {
			result += fmt.Sprintln("Error sending move", err)
			return
		}

		result += fmt.Sprintln("Made move\n", prettyJSON(move))
		return

	case "moves":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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

		moves, err := p.State.ValidMoves()
		if err != nil {
			result += fmt.Sprintln(err)
			return
		}
		for i, move := range moves {
			result += moveString(i, move)
		}
	case "move":
		if p.SplendorClient.UserID.IsZero() {
			result += fmt.Sprintln("please login first")
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
		err = p.sendMove(ctx, move)
		if err != nil {
			result += fmt.Sprintln("Error sending move", err)
			return
		}

		result += fmt.Sprintln("Made move\n", prettyJSON(move))
		return

	default:
		result += fmt.Sprintln("Commands: board hand gems cards moves exit\n")
	}
	return

}

func (p *Terminal) maybeFetchState(ctx context.Context) error {
	if p.SplendorClient == nil {
		return nil
	}
	state, err := p.SplendorClient.GetState(ctx, p.CurrentGame)
	if err != nil {
		return err
	}
	p.State = *state
	// state, err := p.Game.DeserializeState(&game.SerializedObject{
	// 	ID:   p.CurrentGame,
	// 	Data: packet.Payload,
	// })
	// if err != nil {
	// 	return err
	// }
	// typed, ok := state.(splendor.State)
	// if !ok {
	// 	return fmt.Errorf("Wrong game")
	// }
	// p.State = typed
	return nil
}

func (p *Terminal) sendMove(ctx context.Context, move game.Move) error {
	made, err := p.maybeSendMove(ctx, move)
	if err != nil {
		return err
	}
	if made {
		return nil
	}
	return nil
}

func (p *Terminal) maybeSendMove(ctx context.Context, move game.Move) (bool, error) {
	if p.SplendorClient == nil {
		return false, nil
	}
	packet, err := p.Message.MessagePlayerMove(move, p.SplendorClient.UserID)
	if err != nil {
		return false, err
	}
	_, err = p.SplendorClient.SendPacket(ctx, p.CurrentGame, p.SplendorClient.UserID, *packet)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *Terminal) retryListen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		p.SplendorClient.Connect(ctx, nil)
		err := p.SplendorClient.Listen(ctx, p.Handle)
		if err != nil {
			// fmt.Println("Error listening for websocket", err)
		}
		p.SplendorClient.Close(ctx)
		time.Sleep(time.Second * 10)
	}
}

func (p *Terminal) Handle(ctx context.Context, packet wire.Packet) error {
	p.Packets <- packet
	return nil
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
