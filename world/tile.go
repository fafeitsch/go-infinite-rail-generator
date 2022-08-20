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
	// alpha := [16]Connectors{}
	// for track := 0; track < 16; track++ {
	// 	for _, connector := range t.Tracks.Gamma[track] {
	// 		if connector.TargetColumn != Omega {
	// 			continue
	// 		}
	// 		if alpha[connector.TargetTrack] == nil {
	// 			alpha[connector.TargetTrack] = make([]*Connector, 0, 0)
	// 		}
	// 		alpha[connector.TargetTrack] = append(alpha[connector.TargetTrack], &Connector{
	// 			TargetColumn: Beta,
	// 			TargetTrack:  track,
	// 		})
	// 	}
	// }
	// t.Tracks.Alpha = alpha
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

// FindConnector returns the connector that connects to the specified target and slot.
// Returns nil if no such connector exists.
func (t *Tile) FindConnector(target Column, slot int) *Connector {
	for index, element := range t.Tracks {
		if element.TargetColumn == target && element.TargetTrack == slot {
			return &t.Tracks[index]
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
