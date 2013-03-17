package main

import (
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	field := MakeField()
	team := MakeTeam()
	team[0] = Unit{name: 'Y', x: 1, y: 2}

loop:
	for {
		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
			field.Draw()
			team.Draw()
		termbox.Flush()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			dx, dy := 0, 0
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyArrowUp:
				dy = -1
			case termbox.KeyArrowDown:
				dy = 1
			case termbox.KeyArrowLeft:
				dx = -1
			case termbox.KeyArrowRight:
				dx = 1
			}
			team[0].x += dx
			team[0].y += dy
		}
	}
	return
}

type Field [][]FieldCell
type FieldCell struct {
	biome int
}

func MakeField() Field {
	field := make(Field, 80)
	n2d := NewNoise2DContext(0)

	for i := range field {
		field[i] = make([]FieldCell, 80)
		for j := range field[i] {
			v := n2d.Get(float32(i) * 0.1, float32(j) * 0.1)
			v = v * 0.5 + 0.5
			field[i][j].biome = int( v / 0.3)
			
		}
	}
	return field
}

var biomeColors = map[int]termbox.Attribute{
	0: termbox.ColorBlue,
	1: termbox.ColorGreen,
	2: termbox.ColorGreen,
	9: termbox.ColorBlack,
	10: termbox.ColorWhite,
}

func (field Field) Draw() {
	for y := range field {
		for x := range field[y] {
			termbox.SetCell(x, y, ' ', termbox.ColorWhite, biomeColors[field[y][x].biome])
		}
	}
	return
}

type Team []Unit
type Unit struct {
	name   rune
	id     int
	health int
	x      int
	y      int
}

func MakeTeam() Team {
	return make(Team, 10)
}

func (team Team) Draw() {
	for _, unit := range team {
		termbox.SetCell(unit.x, unit.y, unit.name, termbox.ColorWhite,
termbox.ColorDefault)
	}
	return
}
