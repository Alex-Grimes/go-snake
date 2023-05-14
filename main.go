package main

import (
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
