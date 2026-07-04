package request

import (
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Response struct {
	StatusCode int
	Body       string
}

type Model struct {
	methods []string
	method  int
	url     textinput.Model

	response *Response
}

type responseMsg struct {
	Response Response
	Err      error
}

func New() Model {
	url := textinput.New()
	url.SetValue("https://jsonplaceholder.typicode.com/users")
	url.Focus()
	url.CharLimit = 300
	url.Width = 80

	return Model{
		methods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		method:  0,
		url:     url,
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

	case responseMsg:
		if msg.Err != nil {
			m.response = &Response{
				StatusCode: 0,
				Body:       msg.Err.Error(),
			}
			return m, nil
		}

		m.response = &msg.Response
		return m, nil
	}

	m.url, cmd = m.url.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	view := fmt.Sprintf(
		`New Request

Method: %s

URL:
%s

←/→ change method
Ctrl+R send request
Esc back
`,
		m.methods[m.method],
		m.url.View(),
	)

	if m.response != nil {
		view += fmt.Sprintf(
			"\nStatus: %d\n\n%s",
			m.response.StatusCode,
			m.response.Body,
		)
	}

	return view
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
