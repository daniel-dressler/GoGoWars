package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)

	MainMenu()

	field := MakeField(DeductUI(termbox.Size()))
	team := MakeTeam()
	raster := MakeRaster(team, field)
	team[0] = Unit{name: 'イ', x: 1, y: 2, health: 10, movePoints: 3}
	team[1] = Unit{name: 'ヒ', x: 1, y: 3, health: 5, movePoints: 1}

loop:
	for {
		for i := range team {
			for move := 0; move < team[i].movePoints; move++ {
				termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
				raster.DrawTerrain()
				raster.DrawUnits()
				raster.DrawUi()
				raster.DrawUiMsg(fmt.Sprintf("Health: %d",
					team[i].health), 0, 0)
				raster.DrawUiMsg(fmt.Sprintf("Turns left: %d",
					team[i].movePoints-move), 0, 1)
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
					team[i].Move(dx, dy, field)
				}
			}
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
			v := n2d.Get(float32(i)*0.1, float32(j)*0.1)
			v = v*0.5 + 0.5
			field[i][j].terrain = noiseToBiome[int(v/0.3)]

		}
	}
	return field
}

/* ------- Unit ------- */
type Team []Unit
type Unit struct {
	name       rune
	id         int
	health     int
	movePoints int
	x          int
	y          int
}

func (this *Unit) Move(dx int, dy int, terrain Field) {
	this.x += dx
	this.y += dy

	if this.x < 0 || this.x >= len(terrain[0]) ||
		this.y < 0 || this.y >= len(terrain) ||
		terrain[this.y][this.x].terrain != BiomeGrass {

		this.x -= dx
		this.y -= dy
	}

	return
}

/* ------- Team ------- */
func MakeTeam() Team {
	return make(Team, 2)
}

/* ------- Raster ----- */
type Raster struct {
	units       Team
	terrain     Field
	chromeColor termbox.Attribute
	textColor   termbox.Attribute
	backColor   termbox.Attribute
}

func MakeRaster(u Team, t Field) *Raster {
	this := new(Raster)
	this.units = u
	this.terrain = t
	this.chromeColor = termbox.ColorYellow
	this.textColor = termbox.ColorWhite
	this.backColor = termbox.ColorBlack
	return this
}

var biomeColors = map[Biome]termbox.Attribute{
	BiomeLake:  termbox.ColorBlue,
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

var heightUI int = 4

func DeductUI(width int, height int) (int, int) {
	return width, height - heightUI
}

func (this Raster) DrawUi() {
	windowHeight := len(this.terrain)
	for x := 0; x < len(this.terrain[0]); x++ {
		for y := 0; y < heightUI; y++ {
			var char rune = ' '
			if y == 0 {
				char = '‾'
			}

			termbox.SetCell(x, windowHeight+y, char,
				this.chromeColor, this.backColor)
		}
	}
}

func (this Raster) DrawUiMsg(msg string, x int, y int) {
	DrawMsg(msg, x, len(this.terrain)+y+1, this.textColor, this.backColor)
}

func DrawMsg(msg string, leftCorner int, topCorner int,
	fg termbox.Attribute, bg termbox.Attribute) {
	x := leftCorner
	y := topCorner
	for _, c := range msg {
		if c == '\n' {
			y++
			x = leftCorner
		} else {
			termbox.SetCell(x, y, c, fg, bg)
			x++
		}
	}
}

/* ------- main menu ---------- */

func MainMenu() int {
	for {
		termbox.Clear(termbox.ColorWhite, termbox.ColorWhite)
		DrawMsg(logo, 80/2-45/2, 3,
			termbox.ColorYellow, termbox.ColorWhite)
		termbox.Flush()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return 0
			default:
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
	return 1
}

var logo string = `   _____    ____      _____    ____    +---+
  /▒▒▒▒▒|  /▒▒▒▒\    /▒▒▒▒▒|  /▒▒▒▒\   |▒▒▒|
 |▒|  __  |▒|  |▒|  |▒|  __  |▒|  |▒|  |▒▒▒|
 |▒| |▒▒| |▒|  |▒|  |▒| |▒▒| |▒|  |▒|  |▒▒▒|
 |▒|__|▒| |▒|__|▒|  |▒|__|▒| |▒|__|▒|  |▒▒▒|
  \▒▒▒▒▒|  \▒▒▒▒/    \▒▒▒▒▒|  \▒▒▒▒/   |▒▒▒|
                                       |▒▒▒|
 \ \        / /                        +---+
  \ \  /\  / /    ____   ____   ___
   \ \/  \/ /    / _  | | ___| / __|   +---+
    \  /\  /    | (_| | | |    \__ \   |▒▒▒|
     \/  \/      \____| |_|    |___/   +---+`
