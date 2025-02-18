package raptor

// LegType as it stands its name represents the type of leg in a journey
type LegType int

const (
	LEG_TYPE_TRANSIT LegType = iota
	LEG_TYPE_WALKING
)

var conflictTypeStr = []string{"transit", "walking"}

// String returns string representation of the type
func (iotaIdx LegType) String() string {
	return conflictTypeStr[iotaIdx]
}

// JourneyLeg represents a single leg in a journey
type JourneyLeg struct {
	From          Stop
	To            Stop
	DepartureTime int
	ArrivalTime   int
	Trip          Trip
	Type          LegType // When it is `Walking` then trip should be empty
}

// Journey represents a journey from source to destination
type Journey []JourneyLeg

// ExtractJourney reconstructs the journey from source to target
func (r *RAPTOR) ExtractJourney(source, target Stop) Journey {
	journey := make(Journey, 0)
	currentStop := target
	for currentStop != source {
		label := r.TauStar[currentStop]
		if label.BoardingStop == "" {
			break
		}
		if label.Trip != "" {
			// Transit leg
			leg := JourneyLeg{
				From:          label.BoardingStop,
				To:            currentStop,
				Trip:          label.Trip,
				ArrivalTime:   label.EarliestArrivalTime,
				DepartureTime: r.Network.Trips[r.findRoute(label.Trip)][label.Trip].DepartureTime[label.BoardingStop],
				Type:          LEG_TYPE_TRANSIT,
			}
			journey = append([]JourneyLeg{leg}, journey...)
		} else {
			// Walking leg
			leg := JourneyLeg{
				From:          label.BoardingStop,
				To:            currentStop,
				Trip:          "",
				ArrivalTime:   label.EarliestArrivalTime,
				DepartureTime: label.EarliestArrivalTime - r.Network.FootPaths[label.BoardingStop][currentStop],
				Type:          LEG_TYPE_WALKING,
			}
			journey = append([]JourneyLeg{leg}, journey...)
		}
		currentStop = label.BoardingStop
	}
	return journey
}

// findRoute finds the parent route for a given trip
// returns empty string if trip has not been found in any route
func (r *RAPTOR) findRoute(trip Trip) Route {
	for route, trips := range r.Network.Trips {
		if _, exists := trips[trip]; exists {
			return route
		}
	}
	return ""
}
