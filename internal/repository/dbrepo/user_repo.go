package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// Get user by ID
func (m *postgresDBRepo) GetUser(userID int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		*
	from
		users
	where
		user_id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, userID)

	var u models.User
	err := row.Scan(
		&u.UserID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.ImgID,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// Authenticates a user
func (m *postgresDBRepo) Authenticate(email, password string) (int, string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	var accessLevel int

	row := m.DB.QueryRowContext(ctx, "select user_id, password, access_level from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword, &accessLevel)
	if err != nil {
		return id, "", 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, "", 0, err
	}

	return id, hashedPassword, accessLevel, nil
}