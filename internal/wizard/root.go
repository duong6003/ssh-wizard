package wizard

import tea "github.com/charmbracelet/bubbletea"

type Step int

const (
	StepWelcome Step = iota
	StepServerInfo
	StepAuth
	StepKeyInstall
	StepSSHConfig
	StepConnTest
	StepDone
)

type StepConstructors struct {
	Welcome       func() tea.Model
	ServerInfo    func(*State) tea.Model
	Auth          func(*State) tea.Model
	KeyInstall    func(*State) tea.Model
	SSHConfigStep func(*State) tea.Model
	ConnTest      func(*State) tea.Model
	Done          func(*State) tea.Model
}

var constructors StepConstructors

func RegisterStepConstructors(c StepConstructors) {
	constructors = c
}

type RootModel struct {
	step    Step
	state   *State
	current tea.Model
	width   int
	height  int
	err     error
}

func NewRootModel() RootModel {
	state := NewState()
	return RootModel{
		step:    StepWelcome,
		state:   state,
		current: constructors.Welcome(),
	}
}

func (m RootModel) Init() tea.Cmd {
	return m.current.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case AdvanceMsg:
		return m.advance()
	case ErrMsg:
		m.err = msg.Err
		return m, nil
	}

	var cmd tea.Cmd
	m.current, cmd = m.current.Update(msg)
	return m, cmd
}

func (m RootModel) View() string {
	if m.err != nil {
		return "\n  Error: " + m.err.Error() + "\n\n  Press Ctrl+C to exit.\n"
	}
	return m.current.View()
}

func (m RootModel) advance() (RootModel, tea.Cmd) {
	m.step++
	switch m.step {
	case StepServerInfo:
		m.current = constructors.ServerInfo(m.state)
	case StepAuth:
		m.current = constructors.Auth(m.state)
	case StepKeyInstall:
		m.current = constructors.KeyInstall(m.state)
	case StepSSHConfig:
		m.current = constructors.SSHConfigStep(m.state)
	case StepConnTest:
		m.current = constructors.ConnTest(m.state)
	case StepDone:
		m.current = constructors.Done(m.state)
	default:
		return m, tea.Quit
	}
	return m, m.current.Init()
}
