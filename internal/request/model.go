package request

import (
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Response struct {
	StatusCode int
	Body       string
}

type Model struct {
	methods []string
	method  int
	url     textinput.Model

	response         *Response
	responseViewport viewport.Model
	ready            bool
}

type responseMsg struct {
	Response Response
	Err      error
}

var panelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1)

var titleStyle = lipgloss.NewStyle().
	Bold(true)

func New() Model {
	url := textinput.New()
	url.SetValue("https://jsonplaceholder.typicode.com/users")
	url.Focus()
	url.CharLimit = 300
	url.Width = 80

	vp := viewport.New(76, 45)

	return Model{
		methods:          []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		method:           0,
		url:              url,
		responseViewport: vp,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			return m, sendRequest(
				m.methods[m.method],
				m.url.Value(),
			)

		case "left":
			if m.method > 0 {
				m.method--
			}

		case "right":
			if m.method < len(m.methods)-1 {
				m.method++
			}
		}
	case tea.WindowSizeMsg:
		leftWidth := msg.Width / 3
		rightWidth := msg.Width - leftWidth - 6

		m.responseViewport.Width = rightWidth - 4
		m.responseViewport.Height = msg.Height - 8

	case responseMsg:
		if msg.Err != nil {
			m.response = &Response{
				StatusCode: 0,
				Body:       msg.Err.Error(),
			}

			m.responseViewport.SetContent(msg.Err.Error())
			return m, nil
		}

		m.response = &msg.Response
		m.responseViewport.SetContent(msg.Response.Body)
		return m, nil
	}

	m.url, cmd = m.url.Update(msg)

	if m.response != nil {
		m.responseViewport, cmd = m.responseViewport.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	requestPanel := fmt.Sprintf(
		`%s

Method: %s

URL:
%s

←/→ change method
Ctrl+R send request
Esc back
	`,
		titleStyle.Render("Request"),
		m.methods[m.method],
		m.url.View(),
	)

	responsePanel := titleStyle.Render("Response")

	if m.response != nil {
		responsePanel = fmt.Sprintf(
			"\n%s\n\nStatus: %d\n\n%s\n",
			titleStyle.Render("Response"),
			m.response.StatusCode,
			m.responseViewport.View(),
		)
	}

	left := panelStyle.
		Width(50).
		Render(requestPanel)

	right := panelStyle.
		Width(80).
		Render(responsePanel)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		right,
	)
}

func (m Model) Method() string {
	return m.methods[m.method]
}

func (m Model) URL() string {
	return m.url.Value()
}

func sendRequest(method, url string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return responseMsg{Err: err}
		}

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return responseMsg{Err: err}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return responseMsg{Err: err}
		}

		return responseMsg{
			Response: Response{
				StatusCode: resp.StatusCode,
				Body:       string(body),
			},
		}
	}
}
