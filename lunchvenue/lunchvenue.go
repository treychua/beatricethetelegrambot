package lunchvenue

import (
	"errors"
)

type lunchVenue struct {
	Location          string `bson:"location"`
	ChosenFrequency   uint   `bson:"chosenfrequency"`
	SelectedFrequency uint   `bson:"selectedfrequency"`
}

// LunchVenues is a slice of lunch venues
type LunchVenues []lunchVenue

// InsertLunchVenue inserts a new lunch venue into the slice if it doesn't already exist
// O(n)
func (lvs *LunchVenues) InsertLunchVenue(location string) error {
	for _, v := range *lvs {
		if v.Location == location {
			return errors.New("LunchVenues::InsertLunchVenue() - Location already exists")
		}
	}

	*lvs = append(*lvs, lunchVenue{location, 0, 0})
	return nil
}

// DeleteLunchVenue deletes an element from the slice. O(n)
func (lvs *LunchVenues) DeleteLunchVenue(location string) {
	for i, v := range *lvs {
		if location == v.Location {
			*lvs = append((*lvs)[:i], (*lvs)[:i+1]...)
			return
		}
	}
}

// HasLunchVenue is a simple function that checks if a location exists within the lunch venues slice. O(n)
func (lvs *LunchVenues) HasLunchVenue(location string) bool {
	for _, v := range *lvs {
		if v.Location == location {
			return true
		}
	}

	return false
}

// RandVenue returns a random venue from the list of already defined lunch venues
// func RandVenue() string {
// 	if len(lunchVenues) == 0 {
// 		return "eh.. add some lunch places first leh"
// 	}
// 	return lunchVenues[randIntMapKey(lunchVenues)]
// }

// func randIntMapKey(m map[int]string) int {
// 	i := rand.Intn(len(m))
// 	for k := range m {
// 		if i == 0 {
// 			return k
// 		}
// 		i--
// 	}
// 	panic("never")
// }
