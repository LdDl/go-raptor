package raptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleRoute(t *testing.T) {
	tests := []struct {
		name               string
		network            *TransitNetwork
		from               Stop
		to                 Stop
		departure          int
		rounds             int
		correctJourneyLegs Journey
	}{
		{
			name: "Advantageous transfer via footpath",
			network: &TransitNetwork{
				Stops: map[Stop]struct{}{"A": {}, "B": {}, "C": {}, "D": {}},
				Routes: map[Route][]Stop{
					"R1": {"A", "B", "C"},
					"R2": {"D", "C"},
				},
				Trips: map[Route]map[Trip]Schedule{
					"R1": {
						"T1": {
							ArrivalTime:   map[Stop]int{"A": 0, "B": 5, "C": 16},
							DepartureTime: map[Stop]int{"A": 0, "B": 5, "C": 16},
						},
					},
					"R2": {
						"T42": {
							ArrivalTime:   map[Stop]int{"D": 7, "C": 11},
							DepartureTime: map[Stop]int{"D": 7, "C": 11},
						},
					},
				},
				FootPaths: map[Stop]map[Stop]int{
					"B": {
						"D": 1,
						"C": 99,
					},
				},
			},
			from:      "A",
			to:        "C",
			departure: 0,
			rounds:    10,
			correctJourneyLegs: Journey{
				{Trip: "T1", Type: LEG_TYPE_TRANSIT, From: "A", To: "B", DepartureTime: 0, ArrivalTime: 5},
				{Trip: "", Type: LEG_TYPE_WALKING, From: "B", To: "D", DepartureTime: 5, ArrivalTime: 6},
				{Trip: "T42", Type: LEG_TYPE_TRANSIT, From: "D", To: "C", DepartureTime: 7, ArrivalTime: 11},
			},
		},
		{
			name: "Can use advantageous transfer via footpath due departure time",
			network: &TransitNetwork{
				Stops: map[Stop]struct{}{"A": {}, "B": {}, "C": {}, "D": {}},
				Routes: map[Route][]Stop{
					"R1": {"A", "B", "C"},
					"R2": {"D", "C"},
				},
				Trips: map[Route]map[Trip]Schedule{
					"R1": {
						"T1": {
							ArrivalTime:   map[Stop]int{"A": 0, "B": 5, "C": 16},
							DepartureTime: map[Stop]int{"A": 0, "B": 5, "C": 16},
						},
					},
					"R2": {
						"T42": {
							ArrivalTime:   map[Stop]int{"D": 4, "C": 11},
							DepartureTime: map[Stop]int{"D": 4, "C": 11},
						},
					},
				},
				FootPaths: map[Stop]map[Stop]int{
					"B": {
						"D": 1,
						"C": 99,
					},
				},
			},
			from:      "A",
			to:        "C",
			departure: 0,
			rounds:    10,
			correctJourneyLegs: Journey{
				{Trip: "T1", Type: LEG_TYPE_TRANSIT, From: "A", To: "C", DepartureTime: 0, ArrivalTime: 16},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raptor := NewRAPTOR(tt.network)
			journey := raptor.RunAndExtractJourney(tt.from, tt.to, tt.departure, tt.rounds)
			assert.Equal(t, len(tt.correctJourneyLegs), len(journey))
			for i, leg := range journey {
				assert.Equalf(t, tt.correctJourneyLegs[i], leg, "Incorrect journey leg at index %d", i)
			}
		})
	}
}
