package model

import (
	"database/sql"

	// Standard anonymous sql driver import
	_ "github.com/go-sql-driver/mysql"

	"github.com/jadoint/micro/auth"
	"github.com/jadoint/micro/conn"
	"github.com/jadoint/micro/logger"
)

// User contains contains data from user tables
type User struct {
	ID       int64  `db:"id_user,omitempty" json:"id,omitempty"`
	Username string `db:"username,omitempty" json:"username,omitempty"`
	Password string `db:"password,omitempty" json:"-"`
	Email    string `db:"email,omitempty" json:"email,omitempty"`
	Created  string `db:"created,omitempty" json:"created,omitempty"`
	Name     string `db:"name,omitempty" json:"name,omitempty"`
	AboutMe  string `db:"about_me,omitempty" json:"aboutMe,omitempty"`
	Modified string `db:"modified,omitempty" json:"modified,omitempty"`
}

// UserRegistration contains user registration details
type UserRegistration struct {
	Username        string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password        string `json:"password" validate:"required,min=6,max=255,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
}

// UserLogin contains user login details
type UserLogin struct {
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

// GetUser gets single user including password
func GetUser(clients *conn.Clients, idUser int64) (*User, error) {
	db := clients.DB.Read
	u := &User{}
	err := db.QueryRow(`
		SELECT id_user AS id, username, password, created
		FROM user
		WHERE id_user = ?
		LIMIT 1
		# GetUser`, idUser).
		Scan(&u.ID, &u.Username, &u.Password, &u.Created)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return u, err
	}
	return u, nil
}

// GetUserByUsername gets single user including password
func GetUserByUsername(clients *conn.Clients, username string) (*User, error) {
	db := clients.DB.Read
	u := &User{}
	err := db.QueryRow(`
		SELECT id_user AS id, username, password, created
		FROM user
		WHERE username = ?
		LIMIT 1
		# GetUser`, username).
		Scan(&u.ID, &u.Username, &u.Password, &u.Created)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return u, err
	}
	return u, nil
}

// AddUser inserts user into user table
func AddUser(clients *conn.Clients, ur *UserRegistration) (int64, error) {
	passwordHash, err := auth.GenerateHash(ur.Password)
	if err != nil {
		return 0, err
	}
	res, err := clients.DB.Exec(`
		INSERT INTO user(username, password, email)
		VALUES(?, ?, ?)`,
		ur.Username, passwordHash, ur.Email)
	if err != nil {
		return 0, err
	}
	idUser, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return idUser, nil
}
