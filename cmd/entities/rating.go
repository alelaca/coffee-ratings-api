package entities

import (
	"time"
)

type Rating struct {
	CoffeeType  string
	Stars       Stars
	DateCreated time.Time
}

type Stars struct {
	Given int
	Max   int
}
