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

	gameType := MainMenu()
	if gameType == MenuQuit {
		return

	} else if gameType == MenuSkirmish {
		// skirmish
		field := MakeField(DeductUI(termbox.Size()))
		raster := MakeRaster(field)

		player := MakePlayerCommander(field, raster, RedHill)
		ai := MakePlayerCommander(field, raster, BlueSat)

		var status GameStatus = GameInProgress
		players := [...]*Commander{player, ai}
		i := 0
		for status == GameInProgress {
			status = players[i].Turn()
			i += 1
			i %= len(players)
		}
	}
}

/* --------- Advisor ------------- */
type AdvisorStatus int32

const (
	AdvisorMoreToDo = iota
	AdvisorTurnDone
	AdvisorGameDone
)

type Advisor struct {
	unit   int
	team   *Team
	field  Field
	raster *Raster
}

func (this *Advisor) Move() AdvisorStatus {
	i := this.unit

	for move := 0; move < this.team.units[i].movePoints; move++ {
		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		this.raster.DrawTerrain()
		this.raster.DrawUnits()
		this.raster.DrawUi()
		this.raster.DrawUiMsg(fmt.Sprintf("Health: %d",
			this.team.units[i].health), 0, 0)
		this.raster.DrawUiMsg(fmt.Sprintf("Turns left: %d",
			this.team.units[i].movePoints-move), 0, 1)
		termbox.Flush()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			dx, dy := 0, 0
			switch ev.Key {
			case termbox.KeyEsc:
				return AdvisorGameDone
			case termbox.KeyArrowUp:
				dy = -1
			case termbox.KeyArrowDown:
				dy = 1
			case termbox.KeyArrowLeft:
				dx = -1
			case termbox.KeyArrowRight:
				dx = 1
			}
			this.team.units[i].Move(dx, dy, this.field)
		}
	}

	this.unit++
	if this.unit == len(this.team.units) {
		this.unit = 0
		return AdvisorTurnDone
	}

	return AdvisorMoreToDo
}

func MakeAdvisor(team *Team, field Field, raster *Raster) *Advisor {
	advisor := new(Advisor)
	advisor.field = field
	advisor.team = team
	advisor.raster = raster
	return advisor
}

/* --------- Comander ------------ */
type GameStatus int32

const (
	GameOver = iota
	GameWon
	GameLost
	GameInProgress
)

type Commander struct {
	advisor *Advisor
	team    *Team
	field   Field
}

func MakePlayerCommander(field Field, raster *Raster, aff Affiliation) *Commander {
	this := new(Commander)
	this.team = MakeTeam(aff)
	raster.RegisterTeam(this.team)
	this.field = field
	this.advisor = MakeAdvisor(this.team, this.field, raster)

	this.team.units[0] = Unit{name: '߉', x: 1, y: 2, health: 10, movePoints: 3}
	this.team.units[1] = Unit{name: 'ﾋ', x: 1, y: 3, health: 5, movePoints: 1}
	return this
}

func (this *Commander) Turn() GameStatus {
	var status AdvisorStatus = AdvisorMoreToDo
	for status == AdvisorMoreToDo {
		status = this.advisor.Move()
	}

	if status == AdvisorGameDone {
		return GameOver
	}
	return GameInProgress
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
	n2d := NewNoise2DContext(time.Now().Unix())

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
type Affiliation int32

const (
	BlueSat = iota
	RedHill
)

type Team struct {
	units       []Unit
	affiliation Affiliation
}

func MakeTeam(aff Affiliation) *Team {
	this := new(Team)
	this.units = make([]Unit, 2)
	this.affiliation = aff
	return this
}

/* ------- Raster ----- */
type Raster struct {
	teams        []*Team
	terrain     Field
	chromeColor termbox.Attribute
	textColor   termbox.Attribute
	backColor   termbox.Attribute
}

func MakeRaster(t Field) *Raster {
	this := new(Raster)
	this.terrain = t
	this.chromeColor = termbox.ColorYellow
	this.textColor = termbox.ColorWhite
	this.backColor = termbox.ColorBlack
	this.teams = make([]*Team, 0)
	return this
}

func (this *Raster) RegisterTeam(newteam *Team) {
	this.teams = append(this.teams, newteam)
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

var countryColors = map[Affiliation]termbox.Attribute{
	BlueSat: termbox.ColorBlue,
	RedHill: termbox.ColorRed,
}

func (this Raster) DrawUnits() {
	for i := range this.teams {
		for j := range this.teams[i].units {
			unit := this.teams[i].units[j]
			bg := biomeColors[this.terrain[unit.y][unit.x].terrain]
			fg := countryColors[this.teams[i].affiliation]
			termbox.SetCell(unit.x, unit.y, unit.name, fg, bg)
		}
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
type MenuSelection int32

const (
	MenuQuit = iota
	MenuSkirmish
)

func getWindowMid() int {
	x, _ := termbox.Size()
	return x / 2
}

func MainMenu() int {
	termbox.Clear(termbox.ColorWhite, termbox.ColorWhite)
	logoChannel := make(chan bool)
	go animateLogo(logoChannel)

	DrawMsg(skirmishButton, getWindowMid()-32/2, 17,
		termbox.ColorBlack, termbox.ColorWhite)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				logoChannel <- true
				return MenuQuit
			default:
				logoChannel <- true
				return MenuSkirmish
			}
		}
	}
	return MenuQuit
}

var skirmishButton string = `
+------------------------------+
|    Play a skirmish! ENTER    |
+------------------------------+`

func animateLogo(quit chan bool) {
	frame := 0
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(500 * time.Millisecond)
			DrawMsg(logo[frame], getWindowMid()-50/2, 2,
				termbox.ColorBlack, termbox.ColorWhite)
			termbox.Flush()
			frame++
			frame %= len(logo)
		}
	}
}

var logo = []string{logo_frame1, logo_frame2}
var logo_frame1 string = `                                                  
    _____    ____      _____    ____    +---+      
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
      \/  \/      \____| |_|    |___/   +---+ 
                                                
                                                `

var logo_frame2 string = `   _____    ____        _____    ____       +---+
  /▒▒▒▒▒|  /▒▒▒▒\      /▒▒▒▒▒|  /▒▒▒▒\      |▒▒▒|
 |▒|  __  |▒|  |▒|    |▒|  __  |▒|  |▒|     |▒▒▒|
 |▒| |▒▒| |▒|  |▒|    |▒| |▒▒| |▒|  |▒|    |▒▒▒| 
 |▒|__|▒| |▒|__|▒|    |▒|__|▒| |▒|__|▒|    |▒▒▒| 
  \▒▒▒▒▒|  \▒▒▒▒/      \▒▒▒▒▒|  \▒▒▒▒/    |▒▒▒|  
                                          |▒▒▒|  
                                          |▒▒▒|  
                                          +---+ 
 \ \        / /                                  
  \ \  /\  / /    ____     ____   ___             
   \ \/  \/ /    / _  |   | ___| / __|   +---+  
    \  /\  /    | (_| |   | |    \__ \  |▒▒▒|    
     \/  \/      \____|   |_|    |___/  +---+   `
