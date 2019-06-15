package main

import (
	"math/rand"
	"time"

	"github.com/donbattery/test/cmd"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	cmd.Execute()
}
