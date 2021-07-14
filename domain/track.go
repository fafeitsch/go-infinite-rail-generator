package domain

type Offset int

type Tile struct {
	Tracks []Track
	Offset int
}

type Track struct {
	Switches []Switch
}

type SwitchDirection int

const (
	Diverging SwitchDirection = iota
	Merging
)

type Switch struct {
	Direction SwitchDirection
	TrackSpan int
}
