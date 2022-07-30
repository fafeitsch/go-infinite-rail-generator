package noise

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"hash/fnv"
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

type noise struct {
	source   [512]float64
	sampling int
}

type Generator struct {
	*noise
	Seed      string
	TownNames []string
}

func New(seed string) *Generator {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(seed))
	hashSum := hash.Sum64()
	return &Generator{noise: createFromHashSum(hashSum, 30_000), Seed: seed}
}

func createFromHashSum(hashSum uint64, sampling int) *noise {
	source := rand.NewSource(int64(hashSum))
	random := rand.New(source)
	result := noise{source: [512]float64{}, sampling: sampling}
	for i, _ := range result.source {
		result.source[i] = random.Float64()
	}
	return &result
}

func (g *Generator) derive(sampling int) *noise {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(g.Seed))
	hashSum := hash.Sum64()
	result := createFromHashSum(hashSum+43, sampling)
	return result
}

func (n *noise) interpolate(hectometer int) float64 {
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
