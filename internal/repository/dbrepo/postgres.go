package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/ol-ilyassov/booking-app/internal/models"
	"golang.org/x/crypto/bcrypt"
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
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
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

// GetUserByID returns a user by ID.
func (r *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	stmt := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where id = $1`

	row := r.DB.QueryRowContext(ctx, stmt, id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUser updates user data.
func (r *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5`

	_, err := r.DB.ExecContext(ctx, stmt, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())

	if err != nil {
		return err
	}
	return nil
}

// Authenticate authenticates a user.
func (r *postgresDBRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	stmt := `select id, password from users where email = $1`

	row := r.DB.QueryRowContext(ctx, stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// AllReservations returns a slice of all reservations.
func (r *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	stmt := `
	select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at,
	rm.id, rm.room_name
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	order by r.start_date asc
	`

	rows, err := r.DB.QueryContext(ctx, stmt)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// AllNewReservations returns a slice of all reservations.
func (r *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	stmt := `
	select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
	rm.id, rm.room_name
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	where processed = 0
	order by r.start_date asc
	`

	rows, err := r.DB.QueryContext(ctx, stmt)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// GetReservationByID returns one reservation by ID.
func (r *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation

	stmt := `
	select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
	rm.id, rm.room_name
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	where r.id = $1;`

	row := r.DB.QueryRowContext(ctx, stmt, id)
	err := row.Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomID,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
	)

	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

// UpdateReservation updates reservation info.
func (r *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update reservations set first_name = $1, last_name =$2, email = $3, phone = $4,
	updated_at = $5 where id = $6`

	_, err := r.DB.ExecContext(
		ctx,
		stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteReservationByID deletes one reservation by id.
func (r *postgresDBRepo) DeleteReservationByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	delete from reservations  where id = $1`

	_, err := r.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProcessedForReservation updates processed status for a reservation by id.
func (r *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update reservations set processed = $1 where id = $2`

	_, err := r.DB.ExecContext(ctx, stmt, processed, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	stmt := `select id, room_name, created_at, updated_at from rooms order by room_name`

	rows, err := r.DB.QueryContext(ctx, stmt)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var rm models.Room
		err := rows.Scan(
			&rm.ID,
			&rm.RoomName,
			&rm.CreatedAt,
			&rm.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, rm)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

// GetRestrictionsForRoomsByDate returns restrictions for a room by date range.
func (r *postgresDBRepo) GetRestrictionsForRoomsByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	stmt := `
	  select id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date
	  from room_restrictions where $1 < end_date and $2 >= start_date
	  and room_id = $3
	`

	rows, err := r.DB.QueryContext(ctx, stmt, start, end, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rr models.RoomRestriction
		err := rows.Scan(
			&rr.ID,
			&rr.ReservationID,
			&rr.RestrictionID,
			&rr.RoomID,
			&rr.StartDate,
			&rr.EndDate,
		)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, rr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return restrictions, nil
}
