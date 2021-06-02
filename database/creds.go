package database

import "os"

var (
	Username = os.Getenv("MYSQL_USER")
	Password = os.Getenv("MYSQL_PASSWORD")
	Database = os.Getenv("MYSQL_DATABASE")
	Address  = os.Getenv("MYSQL_ADDRESS")
	Port     = os.Getenv("MYSQL_PORT")
)
