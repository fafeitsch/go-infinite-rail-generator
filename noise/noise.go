package noise

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"hash/fnv"
	"math"
	"math/rand"
)

func RandomSeed() (string, error) {
	b := make([]byte, 20)
	_, err := cryptoRand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

type Noise struct {
	source   [512]float64
	sampling int
	Seed     string
}

func New(seed string) *Noise {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(seed))
	hashSum := hash.Sum64()
	result := createFromHashSum(hashSum)
	result.Seed = seed
	result.sampling = 30_000
	return result
}

func createFromHashSum(hashSum uint64) *Noise {
	source := rand.NewSource(int64(hashSum))
	random := rand.New(source)
	result := Noise{source: [512]float64{}}
	for i, _ := range result.source {
		result.source[i] = random.Float64()
	}
	return &result
}

func (n *Noise) derive(sampling int) *Noise {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(n.Seed))
	hashSum := hash.Sum64()
	result := createFromHashSum(hashSum + 43)
	result.sampling = sampling
	return result
}

func (n *Noise) interpolate(hectometer int) float64 {
	hectometer = hectometer % n.sampling
	if hectometer < 0 {
		hectometer = n.sampling + hectometer
	}
	x := float64(hectometer) * float64(len(n.source)) / float64(n.sampling)
	xMin := int(x)
	xMax := (xMin + 1) % len(n.source)
	deltaX := x - float64(xMin)
	smoothDeltaX := deltaX * deltaX * (3 - 2*deltaX)

	return n.source[xMin]*(1-smoothDeltaX) + n.source[xMax]*smoothDeltaX
}

func (n *Noise) numberOfTracks(hectometer int) int {
	seed := n.interpolate(hectometer)
	return int(math.Ceil((seed * 100) / 25))
}

func (n *Noise) isStation(hectometer int) bool {
	return n.interpolate(hectometer) <= 0.2
}
