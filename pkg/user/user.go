package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	// Standard anonymous sql driver import
	_ "github.com/go-sql-driver/mysql"

	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/hash"
	"github.com/jadoint/micro/pkg/logger"
)

// User contains contains data from user tables
type User struct {
	ID       int64  `db:"id_user,omitempty" json:"id,omitempty"`
	Username string `db:"username,omitempty" json:"username,omitempty"`
	Password string `db:"password,omitempty" json:"-"`
	Email    string `db:"email,omitempty" json:"email,omitempty"`
	Created  string `db:"created,omitempty" json:"created,omitempty"`
	Modified string `db:"modified,omitempty" json:"modified,omitempty"`
}

// Registration contains user registration details
type Registration struct {
	Username       string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password       string `json:"password" validate:"required,min=6,max=255"`
	Email          string `json:"email" validate:"required,email"`
	RecaptchaToken string `json:"recaptchaToken" validate:"required"`
}

// RecaptchaResponse recaptcha response from verification step
// See: https://developers.google.com/recaptcha/docs/v3
type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ErrorCodes  []string `json:"error-codes"`
}

// Login contains user login details
type Login struct {
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

// About contains user_profile details
type About struct {
	Title string `db:"title,omitempty" json:"title,omitempty"`
	About string `db:"about,omitempty" json:"about,omitempty"`
}

// IDs contains a list of user IDs
type IDs struct {
	IDs []int64 `json:"ids" validate:"required"`
}

// Username contains username and ID pair
type Username struct {
	ID       int64  `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
}

// PasswordChange for a user changing passwords
type PasswordChange struct {
	OldPassword string `json:"oldPassword" validate:"required,min=6,max=255"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=255"`
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
			logger.HandleError(err)
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
			logger.HandleError(err)
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
			logger.HandleError(err)
		}
		return u.Username, err
	}
	return u.Username, nil
}

// GetUsernames creates a slice of username and ID pairs
func GetUsernames(clients *conn.Clients, uids *IDs) ([]*Username, error) {
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
				logger.HandleError(err)
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
func AddUser(clients *conn.Clients, ur *Registration, rr *RecaptchaResponse) (int64, error) {
	passwordHash, err := hash.Generate(ur.Password)
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

	lastError := ""
	if len(rr.ErrorCodes) > 0 {
		lastError = strings.Join(rr.ErrorCodes, ",")
	}

	_, err = clients.DB.Exec(`
		INSERT INTO recaptcha_log(id_user, score, action, last_error)
		VALUES(?, ?, ?, ?)`,
		idUser, rr.Score, rr.Action, lastError)
	if err != nil {
		return 0, err
	}

	return idUser, nil
}

// ChangePassword changes a user's password
func ChangePassword(clients *conn.Clients, idUser int64, newPassword string) error {
	passwordHash, err := hash.Generate(newPassword)
	if err != nil {
		return err
	}

	_, err = clients.DB.Exec(`
		UPDATE user
		SET password = ?
		WHERE id_user = ?
		LIMIT 1`,
		passwordHash, idUser)
	if err != nil {
		return err
	}

	return nil
}

// GetAbout gets user_profile details of a user
func GetAbout(clients *conn.Clients, idUser int64) (*About, error) {
	db := clients.DB.Read
	var title sql.NullString
	var about sql.NullString
	err := db.QueryRow(`
		SELECT title, about
		FROM user_profile
		WHERE id_user = ?
		LIMIT 1
		# GetAbout
	`, idUser).
		Scan(&title, &about)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.HandleError(err)
		}
		return nil, err
	}
	var a About
	if title.Valid {
		a.Title = title.String
	}
	if about.Valid {
		a.About = about.String
	}
	return &a, nil
}

// UpdateAbout updates user_profile details of a user
func UpdateAbout(clients *conn.Clients, idUser int64, a *About) error {
	_, err := clients.DB.Exec(`
		UPDATE user_profile
		SET title = ?, about = ?
		WHERE id_user = ?
		LIMIT 1`,
		a.Title, a.About, idUser)
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
