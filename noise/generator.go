package noise

import (
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"math"
	"math/rand"
)

func (n *Noise) Generate(hectometer int) domain.Tile {
	tile := n.generateShallowTile(hectometer)
	left := n.generateShallowTile(hectometer - 1)
	right := n.generateShallowTile(hectometer + 1)

	tracks := mandatorySwitches(tile, left, right)

	return domain.Tile{Tracks: tracks}
}

type tileGenerator struct {
	seed           float64
	random         rand.Rand
	leftGenerator  *tileGenerator
	rightGenerator *tileGenerator
}

func mandatorySwitches(tile shallowTile, left shallowTile, right shallowTile) []domain.Track {
	result := make([]domain.Track, tile.tracks)
	for track := 0; track < tile.tracks; track++ {
		result[track] = domain.Track{Switches: make([]int, 0, 0)}
	}
	if tile.tracks > right.tracks {
		mergingSwitches(tile, right.tracks, result)
	} else {
		divergingSwitches(tile.tracks, right, result)
	}
	for track := range result {
		span := getSwitchSpan(track, tile.tracks, left.tracks)
		if span != 0 && tile.bumperLeft[track] {
			result[track].BumperLeft = true
		}
	}
	return result
}

func mergingSwitches(left shallowTile, right int, tracks []domain.Track) {
	for index := range tracks {
		span := getSwitchSpan(index, left.tracks, right)
		if left.bumperRight[index] {
			tracks[index].BumperRight = span != 0
			continue
		}
		if span != 0 {
			tracks[index].Switches = append(tracks[index].Switches, span)
		}
	}
}

func divergingSwitches(left int, right shallowTile, tracks []domain.Track) {
	for track := 0; track < right.tracks; track++ {
		span := -getSwitchSpan(track, right.tracks, left)
		if right.bumperLeft[track] {
			continue
		}
		if span < 0 {
			tracks[0].Switches = append(tracks[0].Switches, span)
		} else if span > 0 {
			tracks[len(tracks)-1].Switches = append(tracks[len(tracks)-1].Switches, span)
		}
	}
}

func getSwitchSpan(track int, tracks int, otherTracks int) int {
	offset := -tracks / 2
	lowerLimit := -otherTracks / 2
	upperLimit := lowerLimit + otherTracks - 1
	if (offset + track) < lowerLimit {
		return int(math.Abs(math.Abs(float64(lowerLimit)) - math.Abs(float64(offset+track))))
	} else if (offset + track) > upperLimit {
		return -(offset + track - upperLimit)
	}
	return 0
}
