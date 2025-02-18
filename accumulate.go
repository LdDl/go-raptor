package raptor

// accumulateRoutes accumulates the marked stops into a map of routes to stops
func (r *RAPTOR) accumulateRoutes() map[Route]Stop {
	routeMarkedStops := make(map[Route]Stop)
	for p := range r.MarkedStops {
		for route, stops := range r.Network.Routes {
			for _, stop := range stops {
				if stop == p {
					if existingStop, exists := routeMarkedStops[route]; exists {
						if r.comesBefore(p, existingStop, route) {
							// Substitute (r, p') by (r, p)
							routeMarkedStops[route] = p
						}
					} else {
						routeMarkedStops[route] = p
					}
				}
			}
		}
		delete(r.MarkedStops, p)
	}
	return routeMarkedStops
}

// comesBefore checks if stop p1 comes before stop p2 in the given route
func (r *RAPTOR) comesBefore(p1, p2 Stop, route Route) bool {
	stops := r.Network.Routes[route]
	for _, stop := range stops {
		if stop == p1 {
			return true
		}
		if stop == p2 {
			return false
		}
	}
	return false
}
