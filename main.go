package main

import (
	"fmt"
	"os"

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

func printUnfilledRectangle(frameOriginX, frameOriginY, FRAME_WIDTH, FRAME_HEIGHT, FRAME_BORDER_THICKNESS int, FRAME_BORDER_VERTICAL, FRAME_BORDER_HORIZONTAL, FRAME_BORDER_TOP_LEFT, FRAME_BORDER_TOP_RIGHT, FRAME_BORDER_BOTTOM_LEFT, FRAME_BORDER_BOTTOM_RIGHT rune) {
	panic("unimplemented")
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
