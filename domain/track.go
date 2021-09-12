package domain

import (
	"fmt"
)

type Offset int

type Tile struct {
	Seed      string
	Offset    int
	Tracks    Tracks
	Platforms Platforms
	Station   bool
}

func NewTile(seed string, tracks int) *Tile {
	if tracks > 16 || tracks < 0 {
		panic(fmt.Sprintf("the number of tracks must be between 16 and 0 (inclusive), but was %d", tracks))
	}
	first := 8 - tracks/2
	rails := Tracks{}
	for track := first; track < first+tracks; track++ {
		rails.Alpha[track] = []*Connector{{Target: Beta, Slot: track}}
		rails.Beta[track] = []*Connector{{Target: Gamma, Slot: track}}
		rails.Gamma[track] = []*Connector{{Target: Omega, Slot: track}}
	}
	return &Tile{Seed: seed, Tracks: rails}
}

func (t *Tile) countOmegaTracks() int {
	targets := make(map[int]int)
	for _, connectors := range t.Tracks.Gamma {
		for _, connector := range connectors {
			if connector.Target == Omega {
				targets[connector.Slot] = targets[connector.Slot] + 1
			}
		}
	}
	result := 0
	for _, element := range targets {
		if element > 0 {
			result = result + 1
		}
	}
	return result
}

type Tracks struct {
	Alpha [16]Connectors
	Beta  [16]Connectors
	Gamma [16]Connectors
}

type Platforms struct {
	Alpha [17]bool
	Beta  [17]bool
	Gamma [17]bool
}

func (t *Tracks) Get(target Column) [16]Connectors {
	switch target {
	case Alpha:
		return t.Alpha
	case Beta:
		return t.Beta
	case Gamma:
		return t.Gamma
	}
	panic(fmt.Sprintf("no connectors from column %v available", target))
}

func (t *Tracks) AlphaTracks() [16]bool {
	var result [16]bool
	for i, connectors := range t.Alpha {
		result[i] = len(connectors) > 0
	}
	return result
}

func (t *Tracks) BuildConnectorMap(source Column, target Column) [16]bool {
	var result [16]bool
	for _, connectors := range t.Get(source) {
		for check := 0; check < len(result); check++ {
			result[check] = result[check] || connectors.ConnectsTo(target, check)
		}
	}
	return result
}

type Connector struct {
	Target Column
	Slot   int
}

type Connectors []*Connector

func (c Connectors) ConnectsTo(target Column, slot int) bool {
	for _, element := range c {
		if element.Target == target && element.Slot == slot {
			return true
		}
	}
	return false
}

func (c Connectors) FindConnector(target Column, slot int) *Connector {
	for _, element := range c {
		if element.Target == target && element.Slot == slot {
			return element
		}
	}
	return nil
}

type Column int

const (
	Alpha Column = iota
	Beta
	Gamma
	Omega
)
