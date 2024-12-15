package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	dbdriver = "postgres"
)

var dbSource = ""

var testQueries *Queries

func TestMain(m *testing.M) {
	if LoadEnv() {
		conn, err := sql.Open(dbdriver, dbSource)
		if err != nil {
			log.Fatal("Cannot connect to database", err)
		}
		testQueries = New(conn)

		os.Exit(m.Run())
	}
}

func LoadEnv() bool {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
		return false
	}
	var dbUser = os.Getenv("DB_USER")
	var dbPass = os.Getenv("DB_PASSWORD")
	var dbPort = os.Getenv("DB_PORT")
	var dbName = os.Getenv("DB_NAME")

	dbSource = fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPass, dbPort, dbName)
	return true
}
