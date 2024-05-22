package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type gameSession = struct {
	gameType       int
	gameDifficulty int
	hasGameStarted bool
}

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

	currSession := gameSession{
		gameDifficulty: -1,
		gameType:       -1,
		hasGameStarted: false,
	}

	menuOptions := []string{
		"1: Individual word mode",
		"2: Paragraph mode",
		"3: Exit",
	}

	difficultyOptions := []string{
		"1. Easy: Top Hundred English words",
		"2. Medium: Top Thousand English words",
		"3. Hard: Top Ten Thousand English words",
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
			if currSession.hasGameStarted {
				// TODO: check for user inputs during game time
			} else {
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyCtrlC:
					return
				case tcell.KeyUp:
					if selectedOption > 0 {
						selectedOption--
					}
					if currSession.gameType == -1 {
						displayMenu(screen, "Welcome to not-my-type for practicing your poor typing skills", menuOptions, selectedOption)
					} else {
						displayMenu(screen, "These are the difficulties, choose how you want to lose", difficultyOptions, selectedOption)
					}
				case tcell.KeyDown:
					if selectedOption < len(menuOptions)-1 {
						selectedOption++
					}
					if currSession.gameType == -1 {
						displayMenu(screen, "Welcome to not-my-type for practicing your poor typing skills", menuOptions, selectedOption)
					} else {
						displayMenu(screen, "These are the difficulties, choose how you want to lose", difficultyOptions, selectedOption)
					}

				case tcell.KeyBackspace:
					if currSession.gameDifficulty > -1 {
						currSession.gameDifficulty = -1
						displayMenu(screen, "These are the difficulties, choose how you want to lose", difficultyOptions, selectedOption)
					} else {
						currSession.gameType = -1
						displayMenu(screen, "Welcome to not-my-type for practicing your poor typing skills", menuOptions, selectedOption)
					}
				case tcell.KeyEnter:
					if currSession.gameType == -1 {
						if selectedOption == len(menuOptions)-1 { // Exit option
							return
						}
						currSession.gameType = selectedOption
						displayMenu(screen, "These are the difficulties, choose how you want to lose", difficultyOptions, selectedOption)
					} else {
						currSession.gameDifficulty = selectedOption
						currSession.hasGameStarted = true
						showGameScreen(screen, fmt.Sprintf("let's start your game mode: %v, difficulty: %v. Press space to start", currSession.gameType, currSession.gameDifficulty))
					}
					selectedOption = 0
					continue
				}
			}
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

// func handleOption(screen tcell.Screen, option int) {
// 	screen.Clear()
// 	options := []string{
// 		"1. Easy: Top Hundred English words",
// 		"2. Medium: Top Thousand English words",
// 		"3. Hard: Top Ten Thousand English words",
// 	}
// 	switch option {
// 	case 0:
// 		displayMenu(screen, "Choose your difficulty for the individual mode", options, 0)
// 	case 1:
// 		displayMenu(screen, "Choose your difficulty for the paragraph mode", options, 0)
// 	case 2:
// 		showMessage(screen, "Exiting...")
// 	}
// }

// func showMessage(screen tcell.Screen, message string) {
// 	for i, r := range message {
// 		screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
// 	}
// 	screen.Show()
// }

// func showInstructions(screen tcell.Screen, text string) {
// 	_, y := screen.Size()
//
// 	for i, item := range text {
// 		screen.SetContent(i, y-10, item, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
// 	}
// 	screen.Show()
// }

func showGameScreen(screen tcell.Screen, text string) {
	screen.Clear()

	for i, item := range text {
		screen.SetContent(i, 0, item, nil, tcell.StyleDefault.Foreground(tcell.ColorDefault))
	}
	screen.Show()
}
