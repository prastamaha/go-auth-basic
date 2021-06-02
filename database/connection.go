package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Creds struct {
	Username string
	Password string
	Database string
	Address  string
	Port     string
}

// DB Connection
var DB *sql.DB

// Connect function will assign connection into DB variable
func Connect(creds *Creds) error {
	dsn := creds.Username + ":" + creds.Password + "@tcp(" + creds.Address + ":" + creds.Port + ")/" + creds.Database + "?parseTime=true"
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Set default configuration
	conn.SetMaxIdleConns(10)
	conn.SetMaxOpenConns(50)
	conn.SetConnMaxIdleTime(1 * time.Minute)
	conn.SetConnMaxLifetime(5 * time.Minute)

	// assign conn into DB variable
	DB = conn

	// return nil if no error occoured
	return nil
}
