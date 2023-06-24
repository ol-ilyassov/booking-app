package dbrepo

import (
	"context"
	"time"

	"github.com/ol-ilyassov/booking-app/internal/models"
)

func (r *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the databse.
func (r *postgresDBRepo) InsertReservation(data models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := r.DB.QueryRowContext(ctx, stmt,
		data.FirstName,
		data.LastName,
		data.Email,
		data.Phone,
		data.StartDate,
		data.EndDate,
		data.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the databse.
func (r *postgresDBRepo) InsertRoomRestriction(data models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id) values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.DB.ExecContext(ctx, stmt,
		data.StartDate,
		data.EndDate,
		data.RoomID,
		data.ReservationID,
		time.Now(),
		time.Now(),
		data.RestrictionID,
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if available, otherwise false.
func (r *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	stmt := `select count(id) 
	from room_restrictions
	where room_id = $1 and $2 < end_date and $3 > start_date;`

	row := r.DB.QueryRowContext(ctx, stmt, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRoom returns a slice of available rooms, if any, for given date range.
func (r *postgresDBRepo) SearchAvailabilityForAllRoom(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	stmt := `select r.id, r.room_name
	from rooms r
	where r.id not in (select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);`

	rows, err := r.DB.QueryContext(ctx, stmt, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, nil
	}

	return rooms, nil
}

// GetRoomByID gets room by id.
func (r *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	stmt := `select id, room_name, created_at, updated_at from rooms where id = $1`

	row := r.DB.QueryRowContext(ctx, stmt, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil
}
