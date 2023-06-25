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
	GetRoomByID(id int) (models.Room, error)

	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, password string) (int, string, error)

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservationByID(id int) error
	UpdateProcessedForReservation(id, processed int) error
}
