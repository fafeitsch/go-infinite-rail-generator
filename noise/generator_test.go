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
		left := shallowTile{}
		center := shallowTile{tracks: test.left, bumperLeft: make([]bool, test.left), bumperRight: make([]bool, test.left)}
		right := shallowTile{tracks: test.right, bumperLeft: make([]bool, test.right), bumperRight: make([]bool, test.right)}
		t.Run(nameDiv, func(t *testing.T) {
			_ = compareSwitches(t, center, left, right, test.div)
		})
		nameMerge := fmt.Sprintf("%d Tracks on %d Tracks", test.right, test.left)
		t.Run(nameMerge, func(t *testing.T) {
			_ = compareSwitches(t, right, left, center, test.merge)
		})
	}
}

func compareSwitches(t *testing.T, center shallowTile, left shallowTile, right shallowTile, want [][]int) []domain.Track {
	tracks := mandatorySwitches(center, left, right)
	require.Equal(t, len(want), len(tracks), "Number of tracks differs.")
	for track, element := range tracks {
		result := assert.Equal(t, len(want[track]), len(element.Switches), "Number of tracks of track %d is wrong", track)
		if !result {
			continue
		}
		for index, sw := range element.Switches {
			assert.Equal(t, want[track][index], sw, "Switch %d of track %d is wrong", index, track)
		}
	}
	return tracks
}

func Test_Bumpers(t *testing.T) {
	tests := []struct {
		left             int
		center           int
		right            int
		leftBumpers      []bool
		centerBumpers    []bool
		rightBumpers     []bool
		wantLeftBumpers  []bool
		wantRightBumpers []bool
		wantSwitches     [][]int
	}{
		{
			left:             3,
			center:           1,
			right:            3,
			leftBumpers:      []bool{false, false, false},
			centerBumpers:    []bool{false},
			rightBumpers:     []bool{true, false, true},
			wantLeftBumpers:  []bool{false},
			wantRightBumpers: []bool{false},
			wantSwitches:     [][]int{{}},
		},
		{
			left:             1,
			center:           3,
			right:            3,
			leftBumpers:      []bool{false},
			centerBumpers:    []bool{true, true, true},
			rightBumpers:     []bool{false, false, false},
			wantLeftBumpers:  []bool{true, false, true},
			wantRightBumpers: []bool{false, false, false},
			wantSwitches:     [][]int{{}, {}, {}},
		},
		{
			left:             1,
			center:           3,
			right:            3,
			leftBumpers:      []bool{false},
			centerBumpers:    []bool{false, false, false},
			rightBumpers:     []bool{false, false, false},
			wantLeftBumpers:  []bool{false, false, false},
			wantRightBumpers: []bool{false, false, false},
			wantSwitches:     [][]int{{}, {}, {}},
		},
		{
			left:             3,
			center:           3,
			right:            1,
			leftBumpers:      []bool{false, false, false},
			centerBumpers:    []bool{false, false, true},
			rightBumpers:     []bool{false, false, false},
			wantLeftBumpers:  []bool{false, false, false},
			wantRightBumpers: []bool{false, false, true},
			wantSwitches:     [][]int{{1}, {}, {}},
		},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%d Tracks on %d Tracks with bumpers", test.left, test.right)
		left := shallowTile{tracks: test.left, bumperRight: test.leftBumpers}
		right := shallowTile{tracks: test.right, bumperLeft: test.rightBumpers}
		center := shallowTile{tracks: test.center, bumperLeft: test.centerBumpers, bumperRight: test.centerBumpers}
		t.Run(name, func(t *testing.T) {
			tracks := compareSwitches(t, center, left, right, test.wantSwitches)
			for index, track := range tracks {
				assert.Equal(t, test.wantLeftBumpers[index], track.BumperLeft, "left bumper on track %d wrong", index)
				assert.Equal(t, test.wantRightBumpers[index], track.BumperRight, "right bumper on track %d wrong", index)
			}
		})
	}
}
