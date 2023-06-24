package repository

import (
	"time"

	"github.com/ol-ilyassov/booking-app/internal/models"
)

// Repository pattern.

// DatabaseRepo
type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(data models.Reservation) (int, error)
	InsertRoomRestriction(data models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRoom(start, end time.Time) ([]models.Room, error)
}
