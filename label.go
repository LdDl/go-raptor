package raptor

// Label represents a label in the RAPTOR algorithm
type Label struct {
	// Arrival time τᵢ(p) for stop p in round i
	EarliestArrivalTime int
	// Current trip t used to reach this stop
	Trip Trip
	// Previous stop p' for journey reconstruction
	BoardingStop Stop
}
