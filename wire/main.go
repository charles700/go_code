package main

import "fmt"

type Monster struct {
	Name string
}

func NewMonster() Monster {
	return Monster{
		Name: "kitty",
	}
}

type Player struct {
	Name string
}

func NewPlayer(name string) Player {
	return Player{Name: name}
}

type Mission struct {
	Player  Player
	Monster Monster
}

func NewMission(p Player, m Monster) Mission {
	return Mission{
		Player:  p,
		Monster: m,
	}
}

func (m Mission) Start() {
	fmt.Printf("%s 打败了 %s, 世界和平了\n", m.Player.Name, m.Monster.Name)
}

func main() {
	mission := InitMission("lisi")
	mission.Start()
}
