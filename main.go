package main

import (
"fmt"
"math/rand"
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
	team[0] = Unit{name: 'Y', x:1, y:2}
	
	for {
		fmt.Print(team.Draw(field.Draw()))
		var move int = 0
		fmt.Printf("Movement: ");
		fmt.Scanf("%d", &move)
		fmt.Printf("\n");

		var dx, dy int = 0, 0
		switch move {
		case 8:
			dy = 1
		case 4:
			dx = -1
		case 2:
			dy = -1
		case 6:
			dx = 1
		}

		team[0].x += dx
		team[0].y += dy
	}
	return
}

type Raster [][]byte

func MakeRaster(field Field) Raster {
	ret := make(Raster, len(field))
	for i := range ret {
		ret[i] = make([]byte, len(field[i]))
	}
	return ret
}

func (img Raster) String() string {
	ret := ""
	for _, col := range img {
		for _, char := range col {
			ret += string(char)
		}
		ret += "\n"
	}
	return ret
}
			

type Field [][]FieldCell
type FieldCell struct {
	biome byte
}

func MakeField() Field {
	field := make(Field, 10)
	for i := range field {
		field[i] = make([]FieldCell, 10)
		for j := range field[i] {
			field[i][j].biome = byte(rand.Intn(3)) + 48
		}
	}
	return field
}

func (field Field) Draw() Raster {
	ret := MakeRaster(field)
	for y := range field {
		for x := range field[y] {
			ret[y][x] = ' '
			//ret[y][x] = field[y][x].biome
		}
	}
	return ret
}


type Team []Unit
type Unit struct {
	name byte
	id int
	health int
	x int
	y int
}

func MakeTeam() Team {
	return make(Team, 10)
}

func (team Team) Draw(raster Raster) Raster {
	for _, unit := range team {
		raster[unit.y][unit.x] = unit.name
	}
	return raster
}

