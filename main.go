package main

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer screen.Fini()

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen.SetStyle(defStyle)
	screen.Clear()

	homeScreenText := `package main

import (
  "fmt"
  "net/http"
  "time"
)

func greet(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
  http.HandleFunc("/", greet)
  http.ListenAndServe(":8080", nil)
}`

	// Function to display the text
	displayText(screen, homeScreenText, defStyle)

	// Show the screen
	screen.Show()

	// Display exit instructions
	displayInstructions(screen, homeScreenText)

	// Wait for a key press before exiting
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}
		}
	}
}

// displayText displays the given text on the screen starting at the top-left corner
func displayText(screen tcell.Screen, text string, style tcell.Style) {
	lines := splitLines(text)
	for y, line := range lines {
		for x, r := range line {
			screen.SetContent(x, y, r, nil, style)
		}
	}
}

// displayInstructions displays instructions on how to exit the program
func displayInstructions(screen tcell.Screen, screenText string) {
	instructions := "Press ESC or Ctrl+C to exit"
	y := strings.Count(screenText, "\n")
	for i, r := range instructions {
		screen.SetContent(i, y+3, r, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	}
	screen.Show()
}

// splitLines splits a string into lines
func splitLines(text string) []string {
	lines := []string{}
	line := ""
	for _, r := range text {
		if r == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(r)
		}
	}
	lines = append(lines, line) // Append the last line
	return lines
}
