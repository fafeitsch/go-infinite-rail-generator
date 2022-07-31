package noise

import (
	_ "embed"
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"strings"
)

//go:embed default_town_names.txt
var defaultTownNames string

func (r *rndTile) buildStation() {
	index, startTile, length := r.stationWidth()
	r.replaceTrackLayout(startTile.Tracks.AlphaTracks())
	for i, connected := range startTile.Tracks.AlphaTracks() {
		if !connected {
			continue
		}
		if !r.Platforms.Alpha[i] && !r.Platforms.Alpha[i+1] {
			if startTile.createRandom(98).Float64() < 0.5 && i%2 == 1 {
				r.Platforms.Alpha[i] = true
				r.Platforms.Beta[i] = true
				r.Platforms.Gamma[i] = index != length
			} else {
				r.Platforms.Alpha[i+1] = true
				r.Platforms.Beta[i+1] = true
				r.Platforms.Gamma[i+1] = index != length
			}
		}
	}
	if index == length/2 {
		r.StationName = r.getStationName()
	}
}

func (r *rndTile) stationWidth() (int, *rndTile, int) {
	generator := r.generator.createRndTile
	length := 0
	firstTile := r
	for tile := generator(r.Hectometer - 1); tile.station; tile = generator(r.Hectometer - length - 1) {
		firstTile = tile
		length = length + 1
	}
	index := length
	length = length + 1
	for counter := 1; generator(r.Hectometer + counter).station; counter++ {
		length = length + 1
	}
	return index, firstTile, length
}

func (r *rndTile) replaceTrackLayout(template [16]bool) {
	for i, connected := range template {
		if !connected {
			r.Tracks.Alpha[i] = nil
			r.Tracks.Beta[i] = nil
			r.Tracks.Gamma[i] = nil
			continue
		}
		r.Tracks.Alpha[i] = []*domain.Connector{
			{
				Target: domain.Beta,
				Slot:   i,
			},
		}
		r.Tracks.Beta[i] = []*domain.Connector{
			{
				Target: domain.Gamma,
				Slot:   i,
			},
		}
		r.Tracks.Gamma[i] = []*domain.Connector{
			{
				Target: domain.Omega,
				Slot:   i,
			},
		}
	}
}

func (r *rndTile) getStationName() string {
	if len(r.generator.TownNames) == 0 {
		r.generator.TownNames = strings.Split(defaultTownNames, "\n")
	}
	index := r.createRandom(99).Intn(len(r.generator.TownNames))
	return r.generator.TownNames[index]
}
