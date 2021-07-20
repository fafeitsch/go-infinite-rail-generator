package domain

type Offset int

type Tile struct {
	Seed   string
	Tracks []Track
	Offset int
}

type Track struct {
	Switches    []int
	BumperLeft  bool
	BumperRight bool
}
