package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	AllUsers() bool
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, int, error)
}

type authRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewAuthRepo(conn *sql.DB, a *config.AppConfig) AuthRepo {
	return &authRepo{
		App: a,
		DB:  conn,
	}
}

func (m *authRepo) AllUsers() bool {
	return true
}

// Returns user by ID
func (m *authRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		user_id, first_name, last_name, email, password, access_level
	from
		users
	where
		user_id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.UserID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// Updates a user in a database
func (m *authRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	update
		users
	set
		first_name = $1,
		last_name = $2,
		email = $3,
		access_level = $4 
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
	)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user
func (m *authRepo) Authenticate(email, testPassword string) (int, string, int, error) {
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

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, "", 0, err
	}

	return id, hashedPassword, accessLevel, nil
}