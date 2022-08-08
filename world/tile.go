package world

import (
	"fmt"
)

// Tile contains all information that is required to build one tile of the rail network.
type Tile struct {
	Seed        string
	Hectometer  int
	Tracks      Tracks
	Platforms   Platforms
	StationName string
}

// NewTile creates a new tile. The seed is mainly needed for informational purposes.
// The number of tracks determines the number of tracks entering the tile from the left.
// Panics if the number of tracks is less than 0 or greater than 16.
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

// Tracks contains the track information of a certain tile. The tracks enter the
// tile on its left side through the alpha connectors. The beta connectors reside
// in the first third of the tile, the gamma connectors in the second third, and the
// tracks leave the tile through the omega connectors.
type Tracks struct {
	Alpha [16]Connectors
	Beta  [16]Connectors
	Gamma [16]Connectors
}

// Get returns the requested track column. Panics if the column is unknown.
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

// AlphaTracks returns an array of length 16, where result[i] is true if i
// is a track entering on the left side of the tile.
func (t *Tracks) AlphaTracks() [16]bool {
	var result [16]bool
	for i, connectors := range t.Alpha {
		result[i] = len(connectors) > 0
	}
	return result
}

// BuildConnectorMap returns a array of length 16, where result[i] is true if
// the tile contains at least one track from the source column directly to slot i of the target column.
func (t *Tracks) BuildConnectorMap(source Column, target Column) [16]bool {
	var result [16]bool
	for _, connectors := range t.Get(source) {
		for check := 0; check < len(result); check++ {
			result[check] = result[check] || connectors.ConnectsTo(target, check)
		}
	}
	return result
}

// Platforms holds information on where to draw platforms. The organization
// is the same as for the Tracks struct.
type Platforms struct {
	Alpha [17]bool
	Beta  [17]bool
	Gamma [17]bool
}

// Connector specifies a target connection of a track, i.e. go to column gamma, slot 5.
type Connector struct {
	Target Column
	Slot   int
}

type Connectors []*Connector

// ConnectsTo returns true if the current connector connects to the specified target and slot.
func (c Connectors) ConnectsTo(target Column, slot int) bool {
	return c.FindConnector(target, slot) != nil
}

// FindConnector returns the connector that connects to the specified target and slot.
// Returns nil if no such connector exists.
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
