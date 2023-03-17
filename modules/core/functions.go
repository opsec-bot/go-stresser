package core

import (
	"math/rand"
	"time"
)

func Random() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000
}

func ClearConsole() {
	print("\033[H\033[2J")
}
