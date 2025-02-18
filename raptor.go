package raptor

import "fmt"

type RAPTOR struct {
	// Transit network with stops, routes, trips and footpaths
	Network *TransitNetwork
	// Represents τ_i(p) - earliest known arrival time at `p` with up to `i` trips.
	TauRound map[int]map[Stop]Label
	// Represents τ*(p_i) - the earliest known arrival time at stop `p_i`.
	// Local pruning allows to evade copying labels from the previous round, since
	// it keeps track of the earliest possible time to get to `p_i`
	TauStar map[Stop]Label
	// Tracks stops to process in current round
	MarkedStops map[Stop]bool
}

func NewRAPTOR(network *TransitNetwork) *RAPTOR {
	return &RAPTOR{
		Network:     network,
		TauRound:    make(map[int]map[Stop]Label),
		TauStar:     make(map[Stop]Label),
		MarkedStops: make(map[Stop]bool),
	}
}

func (r *RAPTOR) Initialize(ps Stop, departureTime int) {
	// Initialize round 0 map if doesn't exist
	if r.TauRound[0] == nil {
		r.TauRound[0] = make(map[Stop]Label)
	}

	// Initialize τᵢ(·) ← ∞ and τ*(·) ← ∞ for all stops
	for stop := range r.Network.Stops {
		r.TauRound[0][stop] = Label{Infinity, "", ""}
		r.TauStar[stop] = Label{Infinity, "", ""}
	}
	// Set τ₀(pₛ) ← τ for source stop
	r.TauRound[0][ps] = Label{departureTime, "", ""}
	// Set τ*(pₛ) ← τ for source stop
	r.TauStar[ps] = Label{departureTime, "", ""}
	// Mark source stop pₛ
	r.MarkedStops[ps] = true
}

func (r *RAPTOR) Run(ps, pt Stop, departureTime int, rounds int) {
	// Initialization of the algorithm
	r.Initialize(ps, departureTime)

	for k := 1; k <= rounds; k++ {
		// Copy labels from round k-1 to round k.
		// Ensures τₖ(p) ≤ τₖ₋₁(p) for all stops p.
		r.TauRound[k] = make(map[Stop]Label)
		for stop, label := range r.TauRound[k-1] {
			r.TauRound[k][stop] = label
		}
		// Accumulate routes serving marked stops from previous round.
		Q := r.accumulateRoutes()

		// Traverse each route
		r.traverseRoutes(k, Q, pt)

		// Look at foot-paths
		r.addTransferTime(k)

		// Stopping criterion
		if len(r.MarkedStops) == 0 {
			break
		}
	}
}

// RunAndExtract is simply shorthand for running the algorithm and extracting the journey
func (r *RAPTOR) RunAndExtract(ps, pt Stop, departureTime int, rounds int) Journey {
	r.Run(ps, pt, departureTime, rounds)
	return r.ExtractJourney(ps, pt)
}

// printLabels for debugging purposes
func (r *RAPTOR) printLabels(k int) {
	for stop, label := range r.TauRound[k] {
		fmt.Printf("%s: %d (from %s by %s)\n",
			stop,
			label.EarliestArrivalTime,
			label.BoardingStop,
			label.Trip)
	}
}
