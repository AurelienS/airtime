package model

import (
	"time"
)

type Point struct {
	Lat          float64
	Lng          float64
	Bearing      float64
	Time         time.Time
	GNSSAltitude int
}

type Flight struct {
	Manufacturer   string
	UniqueID       string
	AdditionalData string
	Date           time.Time
	Site           string
	Pilot          string
	ID             string
	Points         []Point
	Phases         []Phase
}

type PhaseType string

const (
	Cruising PhaseType = "cruising"
	Circling PhaseType = "circling"
)

type Phase struct {
	Type         PhaseType
	NumberOfTurn int
	StartIndex   int
	EndIndex     int
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
}

type CirclingType string

const (
	Mixed CirclingType = "mixed"
	Left  CirclingType = "left"
	Right CirclingType = "right"
)

var a = Phase{Type: "test"}
