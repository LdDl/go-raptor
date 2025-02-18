package raptor

// TransitNetwork represents a timetable
type TransitNetwork struct {
	// Set of all stops in the network
	Stops map[Stop]struct{}
	// Set of all routes in the network
	Routes map[Route][]Stop
	// Set of all trips in the network
	Trips map[Route]map[Trip]Schedule
	// Foot-path durations between stops
	FootPaths map[Stop]map[Stop]int
}
