package model

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

// Task is a task that can be runned
type Task struct{}

// NewTask returns a pointer to a new Task
func NewTask() *Task {
	return &Task{}
}

// Run runs the Task
func (t *Task) Run() {
	log.Info("This task has been triggered!")
	time.Sleep(time.Duration(randInt(500, 15000)) * time.Millisecond)
	log.Info("This task has been finished.")
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
