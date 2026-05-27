package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var p *tea.Program

func main() {
	p = tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
