package noise

import (
	"fmt"
	"testing"
)

func TestNoise_NumberOfTracks(t *testing.T) {
	noise := New("555")
	for i := -1000; i < 1000; i++ {
		fmt.Printf("%.1f km: %d\n", float64(i*100)/1000, noise.NumberOfTracks(i))
	}
}
