package main

import (
	"fmt"
	"os"

	"reqtea/internal/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(app.New(), tea.WithAltScreen())

	_, err := p.Run()

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}
