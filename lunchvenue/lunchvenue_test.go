package lunchvenue_test

import "testing"
import "github.com/treychua/beatricethetelegrambot/lunchvenue"

func TestInsertLunchVenue(t *testing.T) {
	lvs := lunchvenue.LunchVenues{}
	err := lvs.InsertLunchVenue("location1")

	if 1 != len(lvs) {
		t.Errorf("Expected length of LunchVenues to be %v but got %v", 1, len(lvs))
	}

	if "location1" != lvs[0].Location {
		t.Errorf("Expected %v but got %v", "location1", lvs[0].Location)
	}

	err = lvs.InsertLunchVenue("location1")
	if err == nil {
		t.Errorf("Expected an error but got nil instead")
	}

	err = lvs.InsertLunchVenue("location2")
	if 2 != len(lvs) {
		t.Errorf("Expected length of LunchVenues to be %v but got %v", 2, len(lvs))
	}
}

func TestHasLunchVenue(t *testing.T) {
	lvs := lunchvenue.LunchVenues{}

	if lvs.HasLunchVenue("l1") {
		t.Errorf("Expected no venue found but received true")
	}

	lvs.InsertLunchVenue("l1")

	if false == lvs.HasLunchVenue("l1") {
		t.Error("Expected venue to be fonud but received false")
	}

}

func TestDeleteLunchVenue(t *testing.T) {
	lvs := lunchvenue.LunchVenues{}

	for i := 1; i <= 5; i++ {
		lvs.InsertLunchVenue("lunch venue " + string(i))
	}

	lvs.DeleteLunchVenue("lunch venue 1")
	if lvs.HasLunchVenue("lunch venue 1") {
		t.Error("Expected lunch venue 1 to be removed but it is still inside")
	}

}
