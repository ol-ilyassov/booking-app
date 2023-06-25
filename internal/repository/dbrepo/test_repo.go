package dbrepo

import (
	"errors"
	"time"

	"github.com/ol-ilyassov/booking-app/internal/models"
)

func (r *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the databse.
func (r *testDBRepo) InsertReservation(data models.Reservation) (int, error) {
	if data.RoomID == 2 {
		return 0, errors.New("insert fail")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the databse.
func (r *testDBRepo) InsertRoomRestriction(data models.RoomRestriction) error {
	if data.RoomID == 3 {
		return errors.New("insert fail")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if available, otherwise false.
func (r *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID > 5 {
		return false, errors.New("failed db connection")
	}
	return false, nil
}

// SearchAvailabilityForAllRoom returns a slice of available rooms, if any, for given date range.
func (r *testDBRepo) SearchAvailabilityForAllRoom(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID gets room by id.
func (r *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("no records")
	}

	return room, nil
}
