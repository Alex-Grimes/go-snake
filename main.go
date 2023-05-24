package main


import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/encoding"
	"github.com/gdamore/tcell/v2"
)

type Coordinate struct {
	x, y int
}

type Snake struct {
	body                        []*Coordinate
	columnVelocity, rowVelocity int
	symbol                      rune
}

type Food struct {
	point  *Coordinate
	symbol rune
}

var snake *Snake
var food *Food
var CoordinateToClear []*Coordinate
var Screen tcell.Screen
var screenW, screenH int
var isGameOver, isGamePaused bool
var score int

const FRAME_WIDTH = 80
const FRAME_HEIGHT = 15
const FRAME_BORDER_THICKNESS = 1
const FRAME_BORDER_VERTICAL = '|'
const FRAME_BORDER_HORIZONTAL = '-'
const FRAME_BORDER_TOP_LEFT = '+'
const FRAME_BORDER_TOP_RIGHT = '+'
const FRAME_BORDER_BOTTOM_LEFT = '+'
const FRAME_BORDER_BOTTOM_RIGHT = '+'
const SNAKE_SYMBOL = 'O'
const FOOD_SYMBOL = '*'

func main() {
	initScreen()
	initializeGameObjects()
	displayFrame()
	displayGameScore()
	userInput := readUserInput()
	var key string
	for !isGameOver {
		if isGamePaused {
			displayGamePausedInfo()
		}
		key = getUserInput(userInput)
		handleUserInput(key)
		updateGameState()
		displayGameObjects()
		time.Sleep(75 * time.Millisecond)
	}

	displayGameOverInfo()
	time.Sleep(3 * time.Second)

}

func updateGameState() {
	if isGamePaused {
		return
	}
	clearScreen()
	updateSnake()
	updateFood()
}

func updateFood() {
	for isFoodOnSnake() {
		CoordinateToClear = append(CoordinateToClear, food.point)
		food.point.x, food.point.y = generateNewFoodCoordinate()
	}
}

func generateNewFoodCoordinate() (int, int) {
	rand.Seed(time.Now().UnixMicro())
	randomX := rand.Intn(FRAME_WIDTH - 2*FRAME_BORDER_THICKNESS)
	randomY := rand.Intn(FRAME_HEIGHT - 2*FRAME_BORDER_THICKNESS)

	newCoordinate := &Coordinate{
		randomX, randomY,
	}

	transformCoordinateInsideFrame(newCoordinate)

	return newCoordinate.x, newCoordinate.y
}

func transformCoordinateInsideFrame(newCoordinate *Coordinate) {
	leftX, rightX, topY, bottomY := getBoundaries()
	newCoordinate.x += leftX + FRAME_BORDER_THICKNESS
	newCoordinate.y += topY + FRAME_BORDER_THICKNESS
	for newCoordinate.x > rightX {
		newCoordinate.x--
	}
	for newCoordinate.y >= bottomY {
		newCoordinate.y--
	}
}

func isFoodOnSnake() bool {
	panic("unimplemented")
}

func updateSnake() {
	panic("unimplemented")
}

func handleUserInput(key string) {
	if key == "Rune[q]" {
		Screen.Fini()
		os.Exit(0)
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	} else if !isGamePaused {
		if key == "Up" && snake.rowVelocity == 0 {
			snake.rowVelocity = -1
			snake.columnVelocity = 0
		} else if key == "Down" && snake.rowVelocity == 0 {
			snake.rowVelocity = 1
			snake.columnVelocity = 0
		} else if key == "Left" && snake.columnVelocity == 0 {
			snake.rowVelocity = 0
			snake.columnVelocity = -1
		} else if key == "Right" && snake.columnVelocity == 0 {
			snake.rowVelocity = 0
			snake.columnVelocity = 1
		}
	}
}

func getUserInput(userInput chan string) string {
	var key string
	select {
	case key = <-userInput:
	default:
		key = ""
	}
	return key
}

func readUserInput() chan string {
	userInput := make(chan string)
	go func() {
		for {
			switch ev := Screen.PollEvent().(type) {
			case *tcell.EventKey:
				userInput <- ev.Name()
			}
		}
	}()
	return userInput
}

func initializeGameObjects() {
	panic("unimplemented")
}

func initScreen() {
	encoding.Register()
	var err error
	Screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprint(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = Screen.Init(); err != nil {
		fmt.Fprint(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	Screen.SetStyle(defStyle)
	screenW, screenH = Screen.Size()

	if screenW < FRAME_WIDTH || screenH < FRAME_HEIGHT {
		fmt.Fprint(os.Stderr, "Screen size too small\n")
		os.Exit(1)
	}
}

func print(x, y, w, h int, style tcell.Style, char rune) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			Screen.SetContent(x+i, y+j, char, nil, style)
		}
	}
}

func getFrameOrigin() (int, int) {
	return (screenW - FRAME_WIDTH) / 2, (screenH - FRAME_HEIGHT) / 2
}

func displayFrame() {
	frameOriginX, frameOriginY := getFrameOrigin()
	printUnfilledRectangle(frameOriginX, frameOriginY, FRAME_WIDTH, FRAME_HEIGHT, FRAME_BORDER_THICKNESS, FRAME_BORDER_VERTICAL, FRAME_BORDER_HORIZONTAL, FRAME_BORDER_TOP_LEFT, FRAME_BORDER_TOP_RIGHT, FRAME_BORDER_BOTTOM_LEFT, FRAME_BORDER_BOTTOM_RIGHT)
	Screen.Show()
}

func displayGameObjects() {
	displaySnake()
	displayFood()
	Screen.Show()
}

func displaySnake() {
	style := tcell.StyleDefault.Foreground(tcell.ColorDarkGreen.TrueColor())
	for _, snakeCoordinate := range snake.body {
		print(snakeCoordinate.x, snakeCoordinate.y, 1, 1, style, snake.symbol)
	}
}

func displayFood() {
	style := tcell.StyleDefault.Foreground(tcell.ColorDarkRed.TrueColor())
	print(food.point.x, food.point.y, 1, 1, style, food.symbol)
}

func displayGamePausedInfo() {
	_, frameY := getFrameOrigin()
	printAtCenter(frameY-2, "Game Paused", true)
	printAtCenter(frameY-1, "Press P to resume", true)
}

func displayGameOverInfo() {
	centerY := (screenH - FRAME_HEIGHT) / 2
	printAtCenter(centerY-2, "Game Over", true)
	printAtCenter(centerY-1, fmt.Sprintf("Your Score : %d", score), false)
}

func displayGameScore() {
	_, frameY := getFrameOrigin()
	printAtCenter(frameY+FRAME_HEIGHT+2, fmt.Sprintf("Score : %d", score), false)
}

func printAtCenter(startY int, content string, trackClear bool) {
	startX := (screenW - len(content)) / 2
	for i := 0; i < len(content); i++ {
		print(startX+i, startY, 1, 1, tcell.StyleDefault, rune(content[i]))
		if trackClear {
			CoordinateToClear = append(CoordinateToClear, &Coordinate{x: startX + i, y: startY})
		}
	}
	Screen.Show()
}

func clearScreen() {
	for _, coordinate := range CoordinateToClear {
		print(coordinate.x, coordinate.y, 1, 1, tcell.StyleDefault, ' ')
	}
}

func printUnfilledRectangle(xOrigin, yOrigin, width, height, borderThickness int, verticalOutline, horizontalOutline, topLeftOutline, topRightOutline, bottomLeftOutline, bottomRightOutline rune) {
	var upperBorder, lowerBorder rune
	verticalBorder := verticalOutline
	for i := 0; i < width; i++ {
		if i == 0 {
			upperBorder = topLeftOutline
			lowerBorder = bottomLeftOutline
		} else if i == width-1 {
			upperBorder = topRightOutline
			lowerBorder = bottomRightOutline
		} else {
			upperBorder = horizontalOutline
			lowerBorder = horizontalOutline
		}
		print(xOrigin+i, yOrigin, borderThickness, borderThickness, tcell.StyleDefault, upperBorder)
		print(xOrigin+i, yOrigin+height-1, borderThickness, borderThickness, tcell.StyleDefault, lowerBorder)
	}

	for i := 1; i < height-1; i++ {
		print(xOrigin, yOrigin+i, borderThickness, borderThickness, tcell.StyleDefault, verticalBorder)
		print(xOrigin+width-1, yOrigin+i, borderThickness, borderThickness, tcell.StyleDefault, verticalBorder)
	}
}

func getBoundaries() (int, int, int, int) {
	originX, originY := getFrameOrigin()
	topY := originY
	bottomY := originY + FRAME_HEIGHT - FRAME_BORDER_THICKNESS
	leftX := originX
	rightX := originX + FRAME_WIDTH - FRAME_BORDER_THICKNESS
	return topY, bottomY, leftX, rightX
}
