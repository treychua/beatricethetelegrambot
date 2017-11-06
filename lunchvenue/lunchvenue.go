package lunchvenue

import "math/rand"

var lunchVenues = make(map[int]string)

// AddLunchVenue adds a new lunch venue to the map
func AddLunchVenue(venue string) {
	lunchVenues[len(lunchVenues)] = venue
}

// ListLunchVenues will display a list of already defined lunch venues
func ListLunchVenues() []string {
	venues := make([]string, len(lunchVenues))
	for _, v := range lunchVenues {
		venues = append(venues, v)
	}

	return venues
}

// RemoveVenue removes an existing lunch venue
func RemoveVenue(venue string) {

}

// RandVenue returns a random venue from the list of already defined lunch venues
func RandVenue() string {
	if len(lunchVenues) == 0 {
		return "eh.. add some lunch places first leh"
	}
	return lunchVenues[randIntMapKey(lunchVenues)]
}

func randIntMapKey(m map[int]string) int {
	i := rand.Intn(len(m))
	for k := range m {
		if i == 0 {
			return k
		}
		i--
	}
	panic("never")
}
