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

	field := MakeField(termbox.Size())
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
			team[0].Move(dx, dy, field)
		}
	}
	return
}

type Biome int32
const (
	BiomeGrass = iota
	BiomeLake
)

type Field [][]FieldCell
type FieldCell struct {
	terrain Biome
}

var noiseToBiome = map[int]Biome{
	0: BiomeLake,
	1: BiomeGrass,
	2: BiomeGrass,
}

func MakeField(width int, height int) Field {
	field := make(Field, height)
	n2d := NewNoise2DContext(0)

	for i := range field {
		field[i] = make([]FieldCell, width)
		for j := range field[i] {
			v := n2d.Get(float32(i) * 0.1, float32(j) * 0.1)
			v = v * 0.5 + 0.5
			field[i][j].terrain = noiseToBiome[int( v / 0.3)]
			
		}
	}
	return field
}

var biomeColors = map[Biome]termbox.Attribute{
	BiomeLake: termbox.ColorBlue,
	BiomeGrass: termbox.ColorGreen,
}

func (field Field) Draw() {
	for y := range field {
		for x := range field[y] {
			termbox.SetCell(x, y, ' ', termbox.ColorWhite,
							biomeColors[field[y][x].terrain])
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

func (this *Unit) Move(dx int, dy int, terrain Field) {
	this.x += dx
	this.y += dy
	
	if  this.x < 0 || this.x >= len(terrain[0]) ||
		this.y < 0 || this.y >= len(terrain) ||
		terrain[this.y][this.x].terrain != BiomeGrass {

		this.x -= dx
		this.y -= dy
	}

	return
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
