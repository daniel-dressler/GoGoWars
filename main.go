package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
	"math"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	gameType := MenuSkirmish//MainMenu()
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
	cur_unit   int
	team   *Team
	field  Field
	raster *Raster
}

func (this *Advisor) Move() AdvisorStatus {
	if this.cur_unit == len(this.team.units) {
		this.cur_unit = 0
		return AdvisorTurnDone
	}

	unit := this.team.units[this.cur_unit]
	for move := 0; move < unit.movePoints; move++ {
		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		this.raster.DrawTerrain()
		this.raster.DrawUnits()
		this.raster.DrawUi()
		this.raster.DrawActiveUnit(this.team, unit)
		this.raster.DrawUiMsg(fmt.Sprintf("Health: %d",
			unit.GetDisplayHealth()), 4, 0)
		this.raster.DrawUiMsg(fmt.Sprintf("Turns left: %d",
			unit.movePoints - move), 4, 1)
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
			unit.Move(dx, dy, this.field)
		}
	}

	this.cur_unit++
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
	this.team.Recruit(0, field)
	this.team.Recruit(1, field)

	this.advisor = MakeAdvisor(this.team, this.field, raster)

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
type UnitPrototype struct {
	name rune
	movePoints int
}

var UnitTypeTable = []UnitPrototype{{'߉', 3},{'ﾋ', 2}}
	
type Unit struct {
	name       rune
	id         int
	health     float64
	movePoints int
	x          int
	y          int
}

func BuildUnit(id int) *Unit {
	this := new(Unit)
	this.name = UnitTypeTable[id].name
	this.movePoints = UnitTypeTable[id].movePoints
	this.health = 10.0
	this.id = id
	return this
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

func (this *Unit) Fire(target int) float64 {
	return DamageBaseMulti[this.id][target] * this.health
}

func (this *Unit) TakeDamage(damage float64) bool {
	this.health -= damage
	return this.IsAlive()
}

func (this *Unit) IsAlive() bool {
	return this.health - 0.1 < 0
}

func (this *Unit) GetDisplayHealth() int {
	return int(math.Ceil(this.health))
}

var DamageBaseMulti = [][]float64{
{1, 1},
{1, 1}}
	

/* ------- Team ------- */
type Affiliation int32

const (
	BlueSat = iota
	RedHill
)

type Point struct {
	x int
	y int
}

var FactoryLocation = map[Affiliation]Point {
	BlueSat: {1, 1},
	RedHill: {3, 3},
}
	

type Team struct {
	units       []*Unit
	affiliation Affiliation
}

func MakeTeam(aff Affiliation) *Team {
	this := new(Team)
	this.units = make([]*Unit, 0)
	this.affiliation = aff
	return this
}

func (this *Team) Recruit(id int, field Field) {
	unit := BuildUnit(id)
	this.units = append(this.units, unit)
	delta := FactoryLocation[this.affiliation]
	unit.Move(delta.x, delta.y, field)
}

/* ------- Raster ----- */
type Raster struct {
	teams       []*Team
	terrain     Field
	uiStart     int
	uiWidth     int
	chromeColor termbox.Attribute
	textColor   termbox.Attribute
	backColor   termbox.Attribute
}

func MakeRaster(t Field) *Raster {
	this := new(Raster)
	this.terrain = t
	this.uiStart = len(t)
	this.uiWidth = len(t[0])
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

func (this *Raster) DrawTerrain() {
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

func (this *Raster) DrawUnit(x int, y int, country Affiliation, unit *Unit) {
	bg := biomeColors[this.terrain[unit.y][unit.x].terrain]
	fg := countryColors[country]
	termbox.SetCell(x, y, unit.name, fg, bg)
}

func (this *Raster) DrawUnits() {
	for i := range this.teams {
		for j := range this.teams[i].units {
			unit := this.teams[i].units[j]
			this.DrawUnit(unit.x, unit.y, this.teams[i].affiliation, unit)
		}
	}
	return
}
var heightUI int = 3

func DeductUI(width int, height int) (int, int) {
	return width, height - heightUI
}

func (this *Raster) DrawChromeCell(x int, y int, char rune) {
	termbox.SetCell(x, this.uiStart+y, char,
		this.chromeColor, this.backColor)
}
	

func (this *Raster) DrawUi() {
	for x := 0; x < this.uiWidth; x++ {
		for y := 0; y < heightUI; y++ {
			var char rune = ' '
			if y == 0 {
				char = '-'
			}
			this.DrawChromeCell(x, y, char)
		}
	}
}

func (this *Raster) DrawActiveUnit(team *Team, active *Unit) {
	// draw corners
	for x, y := 0, 0; y < 2; {
		this.DrawChromeCell(x*2, y*2, '+')
		x++
		y += x / 2
		x %= 2
	}
	// draw -'s
	for y := 0; y < 2; y++ {
		this.DrawChromeCell(1, y*2, '-')
	}
	// draw |'s
	for x := 0; x < 2; x++ {
		this.DrawChromeCell(x*2, 1, '|')
	}
	
	// draw unit
	color := countryColors[team.affiliation]
	termbox.SetCell(1, this.uiStart + 1, active.name, color, this.backColor)
	termbox.SetCell(active.x, active.y, active.name, termbox.ColorWhite, color)
}

	

func (this *Raster) DrawUiMsg(msg string, x int, y int) {
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
