package domain

type Offset int

type Tile struct {
	Tracks []Track
	Offset int
}

type Track struct {
	Switches    []int
	BumperLeft  bool
	BumperRight bool
}
