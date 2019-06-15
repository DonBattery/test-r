package model

import (
	"time"
)

// Result represents a test result.
type Result int

// Test result constants
const (
	PASS Result = 0
	FAIL Result = 1
)

// Report is a collection of package tests.
type Report struct {
	Apps []App
}

// App contains the test results of a single App.
type App struct {
	Name     string
	Duration time.Duration
	Builds   []*Build
}

// Build contains the results of a single build.
type Build struct {
	Name     string
	Duration time.Duration
	Result   Result
}
