package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type GameSession struct {
	CurrentWord      string
	UserInput        string
	WordList         []string
	WordIndex        int
	GameType         int
	GameDifficulty   int
	HasGameStarted   bool
	HasGameCompleted bool
}

const (
	MenuMainTitle       = "Welcome to not-my-type for practicing your poor typing skills"
	MenuDifficultyTitle = "These are the difficulties, choose how you want to lose"
)

var (
	mainMenuOptions = []string{
		"1: Individual word mode",
		"2: Paragraph mode",
		"3: Exit",
	}

	difficultyOptions = []string{
		"1. Easy: Top Hundred English words",
		"2. Medium: Top Thousand English words",
		"3. Hard: Top Ten Thousand English words",
	}
)

func main() {
	setupLogging()

	screen := setupScreen()
	defer screen.Fini()

	session := &GameSession{}
	session.GameDifficulty = -1
	session.GameType = -1

	selection := 0
	pSelection := &selection
	displayMenu(screen, MenuMainTitle, mainMenuOptions, pSelection)

	for {
		ev := screen.PollEvent()
		code := handleEvent(screen, ev, session, pSelection)
		if code == -1 {
			break
		}
	}
}

func setupLogging() {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.Println("Starting program")
}

func setupScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))
	screen.Clear()
	return screen
}

func handleEvent(screen tcell.Screen, ev tcell.Event, session *GameSession, selectedOption *int) int {
	var code int
	switch ev := ev.(type) {
	case *tcell.EventResize:
		screen.Sync()
	case *tcell.EventKey:
		if session.HasGameStarted {
			code = handleGameInput(screen, ev, session)
		} else {
			code = handleMenuInput(screen, ev, session, selectedOption)
		}
	}
	return code
}

func handleMenuInput(screen tcell.Screen, ev *tcell.EventKey, session *GameSession, selectedOption *int) int {
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		return -1
	case tcell.KeyUp:
		*selectedOption = max(0, *selectedOption-1)
	case tcell.KeyDown:
		if session.GameType == -1 {
			*selectedOption = min(len(mainMenuOptions)-1, *selectedOption+1)
		} else {
			*selectedOption = min(len(difficultyOptions)-1, *selectedOption+1)
		}
	case tcell.KeyBackspace:
		if session.GameDifficulty > -1 {
			session.GameDifficulty = -1
		} else {
			session.GameType = -1
		}
	case tcell.KeyEnter:
		if session.HasGameCompleted {
			session.GameDifficulty = -1
			session.GameType = -1
			session.UserInput = ""
			session.CurrentWord = ""
			session.HasGameStarted = false
			session.HasGameCompleted = false
			session.WordList = []string{}
		} else {
			code := handleMenuSelection(screen, session, selectedOption)
			if code == -1 {
				return -1
			}
		}
	}
	updateMenuDisplay(screen, session, selectedOption)
	return 0
}

func handleMenuSelection(screen tcell.Screen, session *GameSession, selectedOption *int) int {
	if session.GameType == -1 {
		if *selectedOption == len(mainMenuOptions)-1 { // Exit option
			return -1
		}
		session.GameType = *selectedOption
		*selectedOption = 0
		displayMenu(screen, MenuDifficultyTitle, difficultyOptions, selectedOption)
	} else {
		*selectedOption = 0
		session.GameDifficulty = *selectedOption
		startGame(session)
	}
	return 0
}

func startGame(session *GameSession) {
	session.HasGameStarted = true
	session.WordList = getWords(session.GameDifficulty)
	session.WordIndex = 0
	if session.GameType == 0 {
		session.CurrentWord = session.WordList[session.WordIndex]
	}
	if session.GameType == 1 {
		wordString := strings.Join(session.WordList, " ")
		session.CurrentWord = wordString
	}
	session.UserInput = ""
}

func updateMenuDisplay(screen tcell.Screen, session *GameSession, selectedOption *int) {
	if session.HasGameCompleted {
		displayMenu(screen, "Congratulations! You completed the game. Press Enter to return to the main menu.", []string{}, selectedOption)
		return
	}
	if session.HasGameStarted {
		showGameScreen(screen, session.CurrentWord, "")
		return
	}
	if session.GameType == -1 {
		displayMenu(screen, MenuMainTitle, mainMenuOptions, selectedOption)
	} else {
		displayMenu(screen, MenuDifficultyTitle, difficultyOptions, selectedOption)
	}
}

func displayMenu(screen tcell.Screen, menuTitle string, menuOptions []string, selectedOption *int) {
	screen.Clear()
	yOffset := 0
	log.Println(menuTitle)

	// Display menu title
	for i, r := range menuTitle {
		screen.SetContent(i, yOffset, r, nil, tcell.StyleDefault)
	}
	yOffset = strings.Count(menuTitle, "\n") + 2
	// Display menu options
	for i, option := range menuOptions {
		style := tcell.StyleDefault
		item := option
		var target int
		if selectedOption == nil {
			target = 0
		} else {
			target = *selectedOption
		}
		if i == target {
			style = style.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack).Underline(true)
			item += "  <-"
		}
		for j, r := range item {
			screen.SetContent(j, i+yOffset, r, nil, style)
		}
	}
	screen.Show()
}

func showGameScreen(screen tcell.Screen, currentWord string, userInput string) {
	screen.Clear()
	for i, r := range currentWord {
		screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
	}

	for i, r := range "Press space to submit word" {
		screen.SetContent(i, 1, r, nil, tcell.StyleDefault)
	}

	var userInputStyle tcell.Style
	for i, r := range userInput {
		if i >= len(currentWord) || r != rune(currentWord[i]) {
			screen.SetContent(i, 4, r, nil, userInputStyle.Background(tcell.ColorRed).Foreground(tcell.ColorWhite))
		} else {
			screen.SetContent(i, 4, r, nil, userInputStyle.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack))
		}
	}
	screen.Show()
}

func handleGameInput(screen tcell.Screen, ev *tcell.EventKey, session *GameSession) int {
	switch ev.Key() {
	case tcell.KeyRune:
		session.UserInput += string(ev.Rune())
		switch session.GameType {
		case 0:
			if session.UserInput == session.CurrentWord+" " {
				session.WordIndex++
				if session.WordIndex < len(session.WordList) {
					session.CurrentWord = session.WordList[session.WordIndex]
					session.UserInput = ""
					showGameScreen(screen, session.CurrentWord, "")
				} else {
					session.HasGameCompleted = true
					return 0
				}
			}
		case 1:
			if session.CurrentWord+" " == session.UserInput {
				session.HasGameCompleted = true
				return 0
			}
		default:
			return -1
		}
	case tcell.KeyBackspace2:
		if len(session.UserInput) > 0 {
			session.UserInput = session.UserInput[:len(session.UserInput)-1]
		}
	case tcell.KeyCtrlC:
		return -1

	case tcell.KeyEnter:
		if session.HasGameCompleted {
			session.GameDifficulty = -1
			session.GameType = -1
			session.UserInput = ""
			session.CurrentWord = ""
			session.HasGameStarted = false
			session.HasGameCompleted = false
			session.WordList = []string{}
		}
		displayMenu(screen, MenuMainTitle, mainMenuOptions, nil)
		return 0
	}

	// Optionally show the current user input on the screen
	if session.HasGameCompleted {
		displayMenu(screen, "Congratulations! You completed the game. Press Enter to return to the main menu.", []string{}, nil)
	} else {
		showGameScreen(screen, session.CurrentWord, session.UserInput)
	}
	return 0
}

func getWords(difficulty int) []string {
	var fileName string
	difficultyFiles := map[int]string{
		0: "topHundred.txt",
		1: "topThousand.txt",
		2: "topTenThousand.txt",
	}

	fileName, ok := difficultyFiles[difficulty]
	if !ok {
		log.Println("Invalid difficulty level")
		return []string{}
	}

	file, err := os.Open(fmt.Sprintf("./lessons/%v", fileName))
	if err != nil {
		log.Println("something went wrong while opening the files")
		return []string{}
	}

	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	var result []string
	for i := 0; i < 10; i++ {
		num := rand.Intn(len(words))
		result = append(result, words[num])
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
