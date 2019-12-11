package main

import (
	"math/rand"
	"time"

	"github.com/donbattery/test-r/cmd"
)

// Some comments here
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	cmd.Execute()
}
