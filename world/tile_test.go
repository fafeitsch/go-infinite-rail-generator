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
		assert.Equal(t, 15, len(tile.Tracks))
		assert.Equal(t, Connector{
			SourceTrack:  6,
			TargetTrack:  6,
			SourceColumn: Alpha,
			TargetColumn: Beta,
		}, tile.Tracks[0])
		assert.Equal(t, Connector{
			SourceTrack:  6,
			TargetTrack:  6,
			SourceColumn: Beta,
			TargetColumn: Gamma,
		}, tile.Tracks[1])
		assert.Equal(t, Connector{
			SourceTrack:  6,
			TargetTrack:  6,
			SourceColumn: Gamma,
			TargetColumn: Omega,
		}, tile.Tracks[2])
	})
}

func TestTile_Mirror(t *testing.T) {
	tile := NewTile("", 2)
	tile.Platforms = append(tile.Platforms, Platform{Track: 4, Column: Gamma})
	tile.Tracks = append(tile.Tracks, Connector{
		TargetColumn: Omega,
		TargetTrack:  9,
		SourceTrack:  8,
		SourceColumn: Beta,
	})
	tile.Tracks = append(tile.Tracks, Connector{
		TargetColumn: Omega,
		TargetTrack:  8,
		SourceTrack:  7,
		SourceColumn: Gamma,
	})
	tile.Mirror()
	assert.Equal(t, Connector{
		SourceTrack:  9,
		TargetTrack:  8,
		TargetColumn: Gamma,
		SourceColumn: Alpha,
	}, tile.Tracks[6])
	assert.Equal(t, Connector{
		SourceTrack:  8,
		TargetTrack:  7,
		TargetColumn: Beta,
		SourceColumn: Alpha,
	}, tile.Tracks[7])
	assert.Equal(t, []Platform{{Track: 4, Column: Beta}}, tile.Platforms)
}
