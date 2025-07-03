package users

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/yousifsabah0/snippets/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Define a new User struct. Notice how the field names and types align
// with the columns in the database "users" table?
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// Define a new UserModel struct which wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

// We'll use the Insert method to add a new record to the "users" table.
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())`
	if _, err = m.DB.Exec(stmt, name, email, string(hashedPassword)); err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			if mysqlError.Number == 1062 && strings.Contains(mysqlError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var (
		id             int
		hashedPassword []byte
	)

	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`
	if err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}

		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}

		return 0, err
	}

	return id, nil
}

// We'll use the Exists method to check if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
