package main

import (
	"fmt"

	"github.com/lddl/go-raptor"
)

func main() {
	network := &raptor.TransitNetwork{
		Stops: map[raptor.Stop]struct{}{"A": {}, "B": {}, "C": {}, "D": {}},
		Routes: map[raptor.Route][]raptor.Stop{
			"R1": {"A", "B", "C"},
			"R2": {"D", "C"},
		},
		Trips: map[raptor.Route]map[raptor.Trip]raptor.Schedule{
			"R1": {
				"T1": {
					ArrivalTime: map[raptor.Stop]int{"A": 0, "B": 5, "C": 45},
					// At "B" we have to wait 15 minutes before departing
					DepartureTime: map[raptor.Stop]int{"A": 0, "B": 20, "C": 45},
				},
			},
			"R2": {
				"T42": {
					ArrivalTime:   map[raptor.Stop]int{"D": 7, "C": 30},
					DepartureTime: map[raptor.Stop]int{"D": 7, "C": 30},
				},
			},
		},
		FootPaths: map[raptor.Stop]map[raptor.Stop]int{
			"B": {
				"D": 1,
				"C": 99,
			},
		},
	}

	origin := raptor.Stop("A")
	destination := raptor.Stop("C")
	departureTime := 0
	rounds := 10
	algo := raptor.NewRAPTOR(network)
	journey := algo.RunAndExtractJourney(origin, destination, departureTime, rounds)
	for _, leg := range journey {
		if leg.Type == raptor.LEG_TYPE_TRANSIT {
			fmt.Printf("Take trip %s from %s to %s\n", leg.Trip, leg.From, leg.To)
			continue
		}
		fmt.Printf("Walk from %s to %s\n", leg.From, leg.To)
	}
}
