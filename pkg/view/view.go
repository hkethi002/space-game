package view

import (
	"space-game/pkg/maths"
	"space-game/pkg/space"

	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

type tickMsg time.Time

type WindowInfo struct {
	width  int
	height int
	ready  bool
}

type ViewModel struct {
	windowInfo WindowInfo
	activePane int
	ship       space.Ship
}

func main() {
	fmt.Println("vim-go")
}

func tick() tea.Cmd {
	return tea.Tick(time.Duration(time.Second), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (model ViewModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tea.Batch(tick(), tea.EnterAltScreen)
}

func (model ViewModel) View() string {

	baseStyle := lipgloss.NewStyle().
		Width(model.windowInfo.width/2 - 2).
		Height(model.windowInfo.height/2 - 2).
		Align(lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))
	topRightStyle := lipgloss.NewStyle().PaddingLeft(model.windowInfo.width/4 - 8).Inherit(baseStyle)

	topLeft := "Sensors\n\n"
	topLeft += fmt.Sprintf("Speed (m/s): %f\n", maths.Magnitude(model.ship.Velocity))
	topLeft += fmt.Sprintf("∆V (m/s): %f\n", model.ship.DeltaV)
	topLeft += fmt.Sprintf("Fuel: %f\n", model.ship.Fuel)

	topRight := ""
	if model.ship.DeltaV > 0.0 {
		topRight += space.ThrusterOn
	} else {
		topRight += space.ThrusterOff
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		baseStyle.Render(topLeft),
		topRightStyle.Render(topRight),
	)
}
