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
	raster := MakeRaster(team, field)
	team[0] = Unit{name: 'Y', x: 1, y: 2}

loop:
	for {
		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
			raster.DrawTerrain()
			raster.DrawUnits()
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

/* --------- Terrain / Field ----- */
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



/* ------- Unit ------- */
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

/* ------- Team ------- */
func MakeTeam() Team {
	return make(Team, 10)
}

/* ------- Raster ----- */
type Raster struct {
	units Team
	terrain Field
}

func MakeRaster(u Team, t Field) *Raster {
	this := new(Raster)
	this.units = u
	this.terrain = t
	return this
}

var biomeColors = map[Biome]termbox.Attribute{
	BiomeLake: termbox.ColorBlue,
	BiomeGrass: termbox.ColorGreen,
}
func (this Raster) DrawTerrain() {
	for y := range this.terrain {
		for x := range this.terrain[y] {
			termbox.SetCell(x, y, ' ', termbox.ColorWhite,
							biomeColors[this.terrain[y][x].terrain])
		}
	}
	return
}

func (this Raster) DrawUnits() {
	for _, unit := range this.units {
		bg := biomeColors[this.terrain[unit.y][unit.x].terrain]
		termbox.SetCell(unit.x, unit.y, unit.name, termbox.ColorWhite, bg)
	}
	return
}
