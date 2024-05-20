package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Starting program")
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

	menuOptions := []string{
		"1: Individual word mode",
		"2: Paragraph mode",
		"3: Exit",
	}

	// exitInstructions := "Press Escape or Ctrl + C to exit anytime"

	selectedOption := 0

	displayMenu(screen, "Welcome to not-my-type for practicing your poor typing skills", menuOptions, selectedOption)

	// showInstructions(screen, exitInstructions)

	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyUp:
				if selectedOption > 0 {
					selectedOption--
				}
			case tcell.KeyDown:
				if selectedOption < len(menuOptions)-1 {
					selectedOption++
				}
			case tcell.KeyEnter:
				handleOption(screen, selectedOption)
				if selectedOption == len(menuOptions)-1 { // Exit option
					return
				}
				continue
			}
			displayMenu(screen, "Welcome to not-my-type for practicing your poor typing skills", menuOptions, selectedOption)
		}
	}
}

func displayMenu(screen tcell.Screen, menuTitle string, menuOptions []string, selectedOption int) {
	log.Println(menuOptions)
	screen.Clear()
	yOffset := 0

	for i, item := range menuTitle {
		screen.SetContent(i, yOffset, item, nil, tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault))
	}
	yOffset = strings.Count(menuTitle, "\n") + 2
	for i, option := range menuOptions {
		style := tcell.StyleDefault
		item := option
		if i == selectedOption {
			style = style.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack).Underline(true)
			item += "  <-"
		}
		for j, r := range item {
			screen.SetContent(j, i+yOffset, r, nil, style)
		}
	}
	screen.Show()
}

func handleOption(screen tcell.Screen, option int) {
	screen.Clear()
	switch option {
	case 0:
		options := []string{
			"1. Easy: Top Hundred English words",
			"2. Medium: Top Thousand English words",
			"3. Hard: Top Ten Thousand English words",
		}
		symp := option == 0
		log.Println(symp, "on line 109")
		displayMenu(screen, "Choose your difficulty", options, 0)
	case 1:
		showMessage(screen, fmt.Sprintf("Current Time: %v", time.Now()))
	case 2:
		showMessage(screen, "Exiting...")
	}
}

func showMessage(screen tcell.Screen, message string) {
	for i, r := range message {
		screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
	}
	screen.Show()
}

func showInstructions(screen tcell.Screen, text string) {
	_, y := screen.Size()

	for i, item := range text {
		screen.SetContent(i, y-10, item, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	}
	screen.Show()
}
