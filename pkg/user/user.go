package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	// Standard anonymous sql driver import
	_ "github.com/go-sql-driver/mysql"

	"github.com/jadoint/micro/pkg/auth"
	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/logger"
)

// User contains contains data from user tables
type User struct {
	ID       int64  `db:"id_user,omitempty" json:"id,omitempty"`
	Username string `db:"username,omitempty" json:"username,omitempty"`
	Password string `db:"password,omitempty" json:"-"`
	Email    string `db:"email,omitempty" json:"email,omitempty"`
	Created  string `db:"created,omitempty" json:"created,omitempty"`
	Title    string `db:"title,omitempty" json:"title,omitempty"`
	About    string `db:"about,omitempty" json:"about,omitempty"`
	Modified string `db:"modified,omitempty" json:"modified,omitempty"`
}

// Registration contains user registration details
type Registration struct {
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=255"`
	Email    string `json:"email" validate:"required,email"`
}

// Login contains user login details
type Login struct {
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

// About contains user_profile details
type About struct {
	Title sql.NullString `db:"title,omitempty" json:"title,omitempty"`
	About sql.NullString `db:"about,omitempty" json:"about,omitempty"`
}

// UserIDs contains a list of user IDs
type UserIDs struct {
	IDs []int64 `json:"ids" validate:"required"`
}

// Username contains username and ID pair
type Username struct {
	ID       int64  `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
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

// GetUsername gets username by id
func GetUsername(clients *conn.Clients, idUser int64) (string, error) {
	db := clients.DB.Read
	var u User
	err := db.QueryRow(`
		SELECT username
		FROM user
		WHERE id_user = ?
		LIMIT 1
		# GetUsername
		`, idUser).
		Scan(&u.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return u.Username, err
	}
	return u.Username, nil
}

// GetUsernames creates a slice of username and ID pairs
func GetUsernames(clients *conn.Clients, uids *UserIDs) ([]*Username, error) {
	var csvUserIds string
	for _, v := range uids.IDs {
		id := strconv.FormatInt(v, 10)
		csvUserIds += fmt.Sprintf(",'%s'", id)
	}
	if len(csvUserIds) == 0 {
		return nil, errors.New("No results found")
	}
	csvUserIds = csvUserIds[1:]

	db := clients.DB.Read
	rows, err := db.Query(`
		SELECT id_user, username
		FROM user
		WHERE id_user IN (` + csvUserIds + `)
		# GetUsernames`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []*Username
	for rows.Next() {
		var u Username
		err := rows.Scan(&u.ID, &u.Username)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Panic(err.Error())
			}
			return nil, err
		}
		names = append(names, &u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return names, nil
}

// AddUser inserts user into user table
func AddUser(clients *conn.Clients, ur *Registration) (int64, error) {
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

// GetAbout gets user_profile details of a user
func GetAbout(clients *conn.Clients, idUser int64) (*User, error) {
	db := clients.DB.Read
	var a About
	err := db.QueryRow(`
		SELECT title, about
		FROM user_profile
		WHERE id_user = ?
		LIMIT 1
		# GetAbout
	`, idUser).
		Scan(&a.Title, &a.About)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return nil, err
	}
	var u User
	if a.Title.Valid {
		u.Title = a.Title.String
	}
	if a.About.Valid {
		u.About = a.About.String
	}
	return &u, nil
}

// UpdateAbout updates user_profile details of a user
func UpdateAbout(clients *conn.Clients, u *User) error {
	_, err := clients.DB.Exec(`
		UPDATE user_profile
		SET title = ?, about = ?
		WHERE id_user = ?
		LIMIT 1`,
		u.Title, u.About, u.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAbout deletes user_profile details of a user
func DeleteAbout(clients *conn.Clients, idUser int64) error {
	_, err := clients.DB.Exec(`
		UPDATE user_profile
		SET title = NULL, about = NULL
		WHERE id_user = ?
		LIMIT 1`,
		idUser)
	if err != nil {
		return err
	}
	return nil
}
