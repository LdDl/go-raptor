package raptor

// Schedule is a struct that holds the arrival and departure times for a given stop
type Schedule struct {
	ArrivalTime   map[Stop]int
	DepartureTime map[Stop]int
}
