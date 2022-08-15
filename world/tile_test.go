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
					Target: Beta,
					Track:  6,
				},
			},
			{
				{
					Target: Beta,
					Track:  7,
				},
			},
			{
				{
					Target: Beta,
					Track:  8,
				},
			},
			{
				{
					Target: Beta,
					Track:  9,
				},
			},
			{
				{
					Target: Beta,
					Track:  10,
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
					Target: Gamma,
					Track:  6,
				},
			},
			{
				{
					Target: Gamma,
					Track:  7,
				},
			},
			{
				{
					Target: Gamma,
					Track:  8,
				},
			},
			{
				{
					Target: Gamma,
					Track:  9,
				},
			},
			{
				{
					Target: Gamma,
					Track:  10,
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
					Target: Omega,
					Track:  6,
				},
			},
			{
				{
					Target: Omega,
					Track:  7,
				},
			},
			{
				{
					Target: Omega,
					Track:  8,
				},
			},
			{
				{
					Target: Omega,
					Track:  9,
				},
			},
			{
				{
					Target: Omega,
					Track:  10,
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

func TestTracks(t *testing.T) {
	t.Run("get column", func(t *testing.T) {
		tracks := Tracks{
			Alpha: [16]Connectors{
				[]*Connector{},
			},
			Beta:  [16]Connectors{nil, []*Connector{}},
			Gamma: [16]Connectors{nil, nil, []*Connector{}},
		}
		assert.Equal(t, [16]Connectors{[]*Connector{}}, tracks.Get(Alpha))
		assert.Equal(t, [16]Connectors{nil, []*Connector{}}, tracks.Get(Beta))
		assert.Equal(t, [16]Connectors{
			nil,
			nil,
			[]*Connector{},
		}, tracks.Get(Gamma))
	})
	t.Run("get should panic", func(t *testing.T) {
		shouldPanic := func() { (&Tracks{}).Get(4) }
		assert.PanicsWithValue(t, "no connectors from column 4 available", shouldPanic)
	})
	t.Run("alpha tracks", func(t *testing.T) {
		tracks := Tracks{
			Alpha: [16]Connectors{
				[]*Connector{},
				[]*Connector{{}},
			},
		}
		assert.Equal(t, [16]bool{
			false,
			true,
		}, tracks.AlphaTracks())
	})
	t.Run("build connector map", func(t *testing.T) {
		tracks := Tracks{
			Alpha: [16]Connectors{
				nil,
				[]*Connector{{Target: Gamma, Track: 3}},
				[]*Connector{
					{Target: Gamma, Track: 2},
					{Target: Beta, Track: 1},
				},
				[]*Connector{{Target: Beta, Track: 1}},
				[]*Connector{{Target: Gamma, Track: 2}},
			},
		}
		result := tracks.BuildConnectorMap(Alpha, Gamma)
		assert.Equal(t, [16]bool{false, false, true, true, false}, result)
	})
}

func TestConnectors(t *testing.T) {
	var connectors Connectors = []*Connector{
		{
			Target: Beta,
			Track:  2,
		}, {Target: Alpha, Track: 1},
	}
	t.Run("connectsTo", func(t *testing.T) {
		assert.True(t, connectors.ConnectsTo(Alpha, 1))
		assert.False(t, connectors.ConnectsTo(Alpha, 2))
	})
	t.Run("find Connector", func(t *testing.T) {
		assert.Equal(t, &Connector{
			Target: Alpha,
			Track:  1,
		}, connectors.FindConnector(Alpha, 1))
		assert.Nil(t, connectors.FindConnector(Beta, 4))
	})
}
