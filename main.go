package main

import "fmt"
import "math/rand" 

func main() {
	field := MakeField()
	team := MakeTeam()
	team[0].name = 'Y'
	team[0].x = 4
	team[0].y = 2
	fmt.Print(team.Draw(field.Draw()))
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
			ret[y][x] = field[y][x].biome
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

