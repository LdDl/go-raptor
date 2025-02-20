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
			name: "Can't use advantageous transfer via footpath due departure time",
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
		{
			name: "Need to use footpath due to too long departure time",
			network: &TransitNetwork{
				Stops: map[Stop]struct{}{"A": {}, "B": {}, "C": {}, "D": {}},
				Routes: map[Route][]Stop{
					"R1": {"A", "B", "C"},
					"R2": {"D", "C"},
				},
				Trips: map[Route]map[Trip]Schedule{
					"R1": {
						"T1": {
							ArrivalTime: map[Stop]int{"A": 0, "B": 5, "C": 45},
							// At "B" we have to wait 15 minutes before departing
							DepartureTime: map[Stop]int{"A": 0, "B": 20, "C": 45},
						},
					},
					"R2": {
						"T42": {
							ArrivalTime:   map[Stop]int{"D": 7, "C": 30},
							DepartureTime: map[Stop]int{"D": 7, "C": 30},
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
				{Trip: "T42", Type: LEG_TYPE_TRANSIT, From: "D", To: "C", DepartureTime: 7, ArrivalTime: 30},
			},
		},
		{
			name: "Can reach destination only by footpath due time-related curcumstances at destination",
			network: &TransitNetwork{
				Stops: map[Stop]struct{}{"A": {}, "B": {}, "C": {}, "D": {}},
				Routes: map[Route][]Stop{
					"R1": {"A", "B", "C"},
					"R2": {"D", "C"},
				},
				Trips: map[Route]map[Trip]Schedule{
					"R1": {
						"T1": {
							ArrivalTime: map[Stop]int{"A": 0, "B": 5, "C": 999},
							// At "B" we have to wait 15 minutes before departing
							DepartureTime: map[Stop]int{"A": 0, "B": 20, "C": 999},
						},
					},
					"R2": {
						"T42": {
							ArrivalTime:   map[Stop]int{"D": 7, "C": 999},
							DepartureTime: map[Stop]int{"D": 7, "C": 999},
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
				{Trip: "", Type: LEG_TYPE_WALKING, From: "B", To: "C", DepartureTime: 5, ArrivalTime: 104},
			},
		},
		{
			name: "Can't reach destination at all",
			network: &TransitNetwork{
				Stops: map[Stop]struct{}{"A": {}, "B": {}, "C": {}, "D": {}},
				Routes: map[Route][]Stop{
					"R1": {"A", "B"},
					"R2": {"D", "A"},
				},
				Trips: map[Route]map[Trip]Schedule{
					"R1": {
						"T1": {
							ArrivalTime: map[Stop]int{"A": 0, "B": 5},
							// At "B" we have to wait 15 minutes before departing
							DepartureTime: map[Stop]int{"A": 0, "B": 20},
						},
					},
					"R2": {
						"T42": {
							ArrivalTime:   map[Stop]int{"D": 7, "A": 42},
							DepartureTime: map[Stop]int{"D": 7, "A": 78},
						},
					},
				},
				FootPaths: map[Stop]map[Stop]int{
					"B": {
						"D": 1,
					},
				},
			},
			from:               "A",
			to:                 "C",
			departure:          0,
			rounds:             10,
			correctJourneyLegs: Journey{
				// Should not be able to reach destination, but seems like result is:
				// {Trip: "T1", Type: LEG_TYPE_TRANSIT, From: "A", To: "C", DepartureTime: 0, ArrivalTime: 0},
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
