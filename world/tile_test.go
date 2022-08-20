package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTile(t *testing.T) {
	t.Run("test panic", func(t *testing.T) {
		shouldPanic := func() {
			NewTile("", 27)
		}
		assert.PanicsWithValue(t, "the number of tracks must be between 16 and 0 (inclusive), but was 27", shouldPanic)
	})
	t.Run("create new tile", func(t *testing.T) {
		tile := NewTile("test-seed", 5)
		assert.Equal(t, "test-seed", tile.Seed)
		assert.Equal(t, [16]Connectors{
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			{
				{
					TargetColumn: Beta,
					TargetTrack:  6,
				},
			},
			{
				{
					TargetColumn: Beta,
					TargetTrack:  7,
				},
			},
			{
				{
					TargetColumn: Beta,
					TargetTrack:  8,
				},
			},
			{
				{
					TargetColumn: Beta,
					TargetTrack:  9,
				},
			},
			{
				{
					TargetColumn: Beta,
					TargetTrack:  10,
				},
			},
			nil,
			nil,
			nil,
			nil,
			nil,
		}, tile.Tracks.Alpha)
		assert.Equal(t, [16]Connectors{
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			{
				{
					TargetColumn: Gamma,
					TargetTrack:  6,
				},
			},
			{
				{
					TargetColumn: Gamma,
					TargetTrack:  7,
				},
			},
			{
				{
					TargetColumn: Gamma,
					TargetTrack:  8,
				},
			},
			{
				{
					TargetColumn: Gamma,
					TargetTrack:  9,
				},
			},
			{
				{
					TargetColumn: Gamma,
					TargetTrack:  10,
				},
			},
			nil,
			nil,
			nil,
			nil,
			nil,
		}, tile.Tracks.Beta)
		assert.Equal(t, [16]Connectors{
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			{
				{
					TargetColumn: Omega,
					TargetTrack:  6,
				},
			},
			{
				{
					TargetColumn: Omega,
					TargetTrack:  7,
				},
			},
			{
				{
					TargetColumn: Omega,
					TargetTrack:  8,
				},
			},
			{
				{
					TargetColumn: Omega,
					TargetTrack:  9,
				},
			},
			{
				{
					TargetColumn: Omega,
					TargetTrack:  10,
				},
			},
			nil,
			nil,
			nil,
			nil,
			nil,
		}, tile.Tracks.Gamma)
	})
}

func TestTile_Mirror(t *testing.T) {
	tile := NewTile("", 2)
	tile.Tracks.Beta[8] = append(tile.Tracks.Beta[8], &Connector{
		TargetColumn: Omega,
		TargetTrack:  9,
	})
	tile.Tracks.Gamma[7] = append(tile.Tracks.Gamma[7], &Connector{
		TargetColumn: Omega,
		TargetTrack:  8,
	})
	tile.Mirror()
	assert.Equal(t, Connectors{
		{
			TargetTrack:  7,
			TargetColumn: Beta,
		},
	}, tile.Tracks.Alpha[7])
	assert.Equal(t, Connectors{
		{TargetTrack: 7, TargetColumn: Beta},
		{TargetTrack: 8, TargetColumn: Beta},
	}, tile.Tracks.Alpha[8])
	// assert.Equal(t, []*Connector{
	// 	{
	// 		TargetTrack:  8,
	// 		TargetColumn: Gamma,
	// 	},
	// }, tile.Tracks.Alpha[9])
}
