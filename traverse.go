package raptor

// traverseRoutes processes each route in the timetable exactly once
func (r *RAPTOR) traverseRoutes(k int, Q map[Route]Stop, pt Stop) {
	for route, p := range Q {
		stops := r.Network.Routes[route]
		startIndex := r.indexOfStop(p, stops)
		if startIndex == -1 {
			continue
		}
		// t = ⊥ (undefined trip initially)
		var t Trip = ""
		for i := startIndex; i < len(stops); i++ {
			pi := stops[i]
			depTime := Infinity
			if t != "" {
				depTime = r.Network.Trips[route][t].DepartureTime[pi]
				arrTime := r.Network.Trips[route][t].ArrivalTime[pi]
				minTime := MinInt(r.TauStar[pi].EarliestArrivalTime, r.TauStar[pt].EarliestArrivalTime)
				if arrTime < minTime {
					// Can the label be improved in this round?
					// Includes local and target pruning
					r.TauRound[k][pi] = Label{arrTime, t, p}
					r.TauStar[pi] = Label{arrTime, t, p}
					r.MarkedStops[pi] = true
				}
			}
			// Can we catch an earlier trip at pi?
			if r.TauRound[k-1][pi].EarliestArrivalTime <= depTime {
				t = r.earliestTrip(route, pi, r.TauRound[k-1][pi].EarliestArrivalTime)
			}
		}
	}
}

// earliestTrip returns the earliest trip in route `r` that one can catch at stop `p`
// i.e. the earliest trip such that τ_dep (t, p) ≤ τ_k₋₁(p).
// This trip may not exist, in which case result is ⊥ (undefined).
func (r *RAPTOR) earliestTrip(route Route, p Stop, time int) Trip {
	earliest := Trip("")
	minTime := Infinity
	for trip, schedule := range r.Network.Trips[route] {
		depTime := schedule.DepartureTime[p]
		if depTime >= time && depTime < minTime {
			minTime = depTime
			earliest = trip
		}
	}
	return earliest
}

// indexOfStop returns the index of a stop in a slice of stops
func (r *RAPTOR) indexOfStop(p Stop, stops []Stop) int {
	for i, stop := range stops {
		if stop == p {
			return i
		}
	}
	return -1
}
