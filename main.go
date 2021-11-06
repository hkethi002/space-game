package main

import (
	"fmt"
	"math"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type Point3 [3]float64
type Vector3 [3]float64

func ScalarProduct(vector Vector3, scalar float64) Vector3 {
	return Vector3{vector[0] * scalar, vector[1] * scalar, vector[2] * scalar}
}

func Magnitude(vector Vector3) float64 {
	return math.Sqrt(math.Pow(vector[0], 2) + math.Pow(vector[1], 2) + math.Pow(vector[2], 2))
}

func Normalize(vector Vector3) Vector3 {
	magnitude := Magnitude(vector)
	return ScalarProduct(vector, 1/magnitude)
}

func Subtract(vectorA, vectorB Vector3) Vector3 {
	return Vector3{vectorA[0] - vectorB[0], vectorA[1] - vectorB[1], vectorA[2] - vectorB[2]}
}

func Add(vectorA, vectorB Vector3) Vector3 {
	return Vector3{vectorA[0] + vectorB[0], vectorA[1] + vectorB[1], vectorA[2] + vectorB[2]}
}

type CelestialBody struct {
	Name                     string
	Mass                     float64
	internalAbsoluteLocation Point3
	internalAbsoluteVelocity Vector3
}

type Model struct {
	choices          []string
	cursor           int
	celestialBodies  []CelestialBody
	frameOfReference int
	selected         map[int]struct{}
	sensors          Sensors
}

type Sensors struct {
	fuel        float64
	deltaV      float64
	orientation Vector3
	velocity    Vector3
	location    Point3
	body        CelestialBody
}

func initialModel() Model {
	return Model{
		// Our shopping list is a grocery list
		choices: []string{"MainThrusters", "LeftThrusters", "RightThrusters", "TopThrusters", "BotThrusters"},
		celestialBodies: []CelestialBody{
			CelestialBody{
				Name:                     "Earth",
				Mass:                     45.0,
				internalAbsoluteLocation: Point3{100.0, 150.3, 56.6},
				internalAbsoluteVelocity: Vector3{0.0, 0.0, 0.0},
			},
		},

		selected: make(map[int]struct{}),
		sensors:  Sensors{fuel: 50.0, orientation: Vector3{1, 0, 0}},
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Duration(time.Second), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (model Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tick()
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return model, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if model.cursor > 0 {
				model.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if model.cursor < len(model.choices)-1 {
				model.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := model.selected[model.cursor]
			if ok {
				delete(model.selected, model.cursor)
			} else {
				model.selected[model.cursor] = struct{}{}
			}
		}
	case tickMsg:
		// Velocity
		model.sensors.velocity = Add(model.sensors.velocity, ScalarProduct(model.sensors.orientation, model.sensors.deltaV))

		// Fuel
		model.sensors.fuel = model.sensors.fuel - model.sensors.deltaV

		// Force
		if _, ok := model.selected[0]; ok {
			model.sensors.deltaV = math.Min(10.0, model.sensors.deltaV+1.0)
		} else {
			model.sensors.deltaV = math.Max(0, model.sensors.deltaV-1.0)
		}
		model.sensors.deltaV = math.Min(model.sensors.fuel, model.sensors.deltaV)

		// if _, lok := model.selected[1]; lok {
		// 	if _, rok := model.selected[2]; rok {
		// 	} else {

		return model, tick()
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return model, nil
}

func (model Model) View() string {
	// The header
	s := "Sensors\n\n"
	s += fmt.Sprintf("Speed (m/s): %f\n", Magnitude(model.sensors.velocity))
	s += fmt.Sprintf("∆V (m/s): %f\n", model.sensors.deltaV)
	s += fmt.Sprintf("Fuel: %f\n", model.sensors.fuel)

	if model.sensors.deltaV > 0.0 {
		s += `
       !
       !
       ^
      / \
     /___\
    |=   =|
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
   /|##!##|\
  / |##!##| \
 /  |##!##|  \
|  / ^ | ^ \  |
| /  ( | )  \ |
|/   ( | )   \|
    ((   ))
   ((  :  ))
   ((  :  ))
    ((   ))
     (( ))
      ( )
       .
       .
       .
`
	} else {
		s += `
       !
       !
       ^
      / \
     /___\
    |=   =|
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
   /|##!##|\
  / |##!##| \
 /  |##!##|  \
|  /       \  |
| /         \ |
|/           \|









`
	}

	s += "\n\nPilot Controls\n\n"

	// Iterate over our choices
	for i, choice := range model.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if model.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := model.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
