package raptor

// addTransferTime scans possible footpaths from marked stops
func (r *RAPTOR) addTransferTime(k int) {
	for p := range r.MarkedStops {
		if _, exists := r.Network.FootPaths[p]; !exists {
			continue
		}
		for pPrime, walkingTime := range r.Network.FootPaths[p] {
			newArrival := r.TauRound[k][p].EarliestArrivalTime + walkingTime
			if _, ok := r.TauRound[k][pPrime]; !ok {
				r.TauRound[k][pPrime] = Label{Infinity, "", ""}
			}
			if newArrival > r.TauRound[k][pPrime].EarliestArrivalTime {
				// We should consider footpaths only if they improve the label
				// rather taking min(x,y) as stated in algorithm section in the paper
				continue
			}
			label := r.TauRound[k][pPrime]
			label.EarliestArrivalTime = newArrival
			// Reassign modified label
			r.TauRound[k][pPrime] = label
			r.TauStar[pPrime] = Label{newArrival, "", p} // Do we need this?
			r.MarkedStops[pPrime] = true
		}
	}
}
