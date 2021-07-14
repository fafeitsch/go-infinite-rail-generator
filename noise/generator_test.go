package noise

import (
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_MandatorySwitches(t *testing.T) {
	tt := []struct {
		left  int
		right int
		div   [][]int
		merge [][]int
	}{
		{left: 4, right: 4, div: [][]int{{}, {}, {}, {}}, merge: [][]int{{}, {}, {}, {}}},
		{left: 1, right: 2, div: [][]int{{-1}}, merge: [][]int{{1}, {}}},
		{left: 1, right: 3, div: [][]int{{-1, 1}}, merge: [][]int{{1}, {}, {-1}}},
		{left: 2, right: 3, div: [][]int{{}, {1}}, merge: [][]int{{}, {}, {-1}}},
		{left: 2, right: 5, div: [][]int{{-1}, {1, 2}}, merge: [][]int{{1}, {}, {}, {-1}, {-2}}},
		{left: 3, right: 4, div: [][]int{{-1}, {}, {}}, merge: [][]int{{1}, {}, {}, {}}},
		{left: 1, right: 5, div: [][]int{{-2, -1, 1, 2}}, merge: [][]int{{2}, {1}, {}, {-1}, {-2}}},
	}
	for _, test := range tt {
		nameDiv := fmt.Sprintf("%d Tracks on %d Tracks", test.left, test.right)
		t.Run(nameDiv, func(t *testing.T) {
			compare(t, test.left, test.right, test.div)
		})
		nameMerge := fmt.Sprintf("%d Tracks on %d Tracks", test.right, test.left)
		t.Run(nameMerge, func(t *testing.T) {
			compare(t, test.right, test.left, test.merge)
		})
	}
}

func compare(t *testing.T,
	left int,
	right int,
	want [][]int) {
	switches := (&tileGenerator{nextTileTracks: right, tracks: left}).mandatorySwitches()
	require.Equal(t, len(want), len(switches), "Number of tracks differs.")
	for track, element := range switches {
		result := assert.Equal(t, len(want[track]), len(element), "Number of switches of track %d is wrong", track)
		if !result {
			continue
		}
		direction := domain.Diverging
		if left > right {
			direction = domain.Merging
		}
		for index, sw := range element {
			assert.Equal(t, domain.Switch{Direction: direction, TrackSpan: want[track][index]}, *sw, "Switch %d of track %d is wrong", index, track)
		}
	}
}

func allSwitchesOfType(switches []*domain.Switch, direction domain.SwitchDirection) bool {
	for _, element := range switches {
		if element.Direction != direction {
			return false
		}
	}
	return true
}
