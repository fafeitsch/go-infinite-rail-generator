package noise

import "math/rand"

type shallowTile struct {
	tracks      int
	bumperLeft  []bool
	bumperRight []bool
}

func (n *Noise) generateShallowTile(hectometer int) shallowTile {
	seed := n.interpolate(hectometer)
	source := rand.NewSource(int64(seed * 10e10))
	random := rand.New(source)

	result := shallowTile{tracks: numberOfTracks(seed)}
	result.bumperRight = make([]bool, result.tracks)
	result.bumperLeft = make([]bool, result.tracks)
	for i := 0; i < result.tracks; i++ {
		result.bumperLeft[i] = random.Float64() < 0.1
		result.bumperRight[i] = random.Float64() < 0.1
	}
	return result
}

func numberOfTracks(seed float64) int {
	if seed < 0.2 {
		return 1
	}
	if seed < 0.6 {
		return 2
	}
	if seed < 0.7 {
		return 3
	}
	return int(seed*10 - 3)
}
