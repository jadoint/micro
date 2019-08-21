package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	// Standard anonymous sql driver import
	_ "github.com/go-sql-driver/mysql"
)

var dsn string

// ClientDB contains database clients
type ClientDB struct {
	Master *sql.DB
	Read   *sql.DB
}

func buildDSN(dbHost string) string {
	hosts := strings.Split(dbHost, ",")
	numHosts := len(hosts)
	randHostIndex := 0
	if numHosts > 1 {
		rand.Seed(time.Now().UnixNano())
		randHostIndex = rand.Intn(numHosts)
	}
	host := hosts[randHostIndex]
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	protocol := os.Getenv("DB_PROTOCOL")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s", user, pass, protocol, host, port, dbName)
}

// InitMasterDB initializes master database connection
func InitMasterDB() (*sql.DB, error) {
	dsn := buildDSN(os.Getenv("DB_MASTER_HOST"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// InitReadDB initializes read (slave) database connection
func InitReadDB() (*sql.DB, error) {
	dsn := buildDSN(os.Getenv("DB_READ_HOST"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetClient initializes database connections and
// returns a ClientDB pointer
func GetClient() (*ClientDB, error) {
	master, err := InitMasterDB()
	if err != nil {
		return nil, err
	}

	read, err := InitReadDB()
	if err != nil {
		return nil, err
	}

	return &ClientDB{
		Master: master,
		Read:   read,
	}, nil
}

// Exec a CREATE/UPDATE/DELETE command on master
func (db *ClientDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	dbMaster := db.Master
	res, err := dbMaster.Exec(query, args...)
	if err != nil {
		return res, err
	}
	return res, nil
}
