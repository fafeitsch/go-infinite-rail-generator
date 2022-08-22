package world

import (
	"fmt"
)

// Tile contains all information that is required to build one tile of the rail network.
type Tile struct {
	Seed        string
	Hectometer  int
	Tracks      Tracks
	Platforms   []Platform
	StationName string
}

// NewTile creates a new tile. The seed is mainly needed for informational purposes.
// The number of tracks determines the number of tracks entering the tile from the left.
// Panics if the number of tracks is less than 0 or greater than 16.
func NewTile(seed string, tracks int) Tile {
	if tracks > 16 || tracks < 0 {
		panic(fmt.Sprintf("the number of tracks must be between 16 and 0 (inclusive), but was %d", tracks))
	}
	first := 8 - tracks/2
	rails := make([]Connector, 0, tracks*3)
	for track := first; track < first+tracks; track++ {
		rails = append(rails, Connector{
			SourceTrack:  track,
			SourceColumn: Alpha,
			TargetTrack:  track,
			TargetColumn: Beta,
		}, Connector{
			SourceTrack:  track,
			SourceColumn: Beta,
			TargetTrack:  track,
			TargetColumn: Gamma,
		},
			Connector{
				SourceTrack:  track,
				SourceColumn: Gamma,
				TargetTrack:  track,
				TargetColumn: Omega,
			},
		)
	}
	return Tile{Seed: seed, Tracks: rails}
}

func (t *Tile) Mirror() {
	for index := range t.Tracks {
		connector := &t.Tracks[index]
		connector.SourceTrack, connector.TargetTrack = connector.TargetTrack, connector.SourceTrack
		connector.SourceColumn, connector.TargetColumn = connector.TargetColumn.Mirror(), connector.SourceColumn.Mirror()
	}
	for index := range t.Platforms {
		t.Platforms[index].Column = t.Platforms[index].Column.Mirror()
	}
}

type Tracks []Connector

// AlphaTracks returns an array of length 16, where result[i] is true if i
// is a track entering on the left side of the tile.
func (t Tracks) AlphaTracks() [16]bool {
	var result [16]bool
	for _, connector := range t {
		if connector.SourceColumn == Alpha {
			result[connector.SourceTrack] = true
		}
	}
	return result
}

type Platform struct {
	Column Column
	Track  int
}

// Connector specifies a target connection of a track, i.e. go to column gamma, slot 5.
type Connector struct {
	SourceColumn Column
	SourceTrack  int
	TargetColumn Column
	TargetTrack  int
	Junction     bool
	Label        string
}

type Column int

func (c Column) Mirror() Column {
	switch c {
	case Alpha:
		return Omega
	case Beta:
		return Gamma
	case Gamma:
		return Beta
	case Omega:
		return Alpha
	}
	panic(fmt.Sprintf("unknown column %d", c))
}

const (
	Alpha Column = iota
	Beta
	Gamma
	Omega
)
