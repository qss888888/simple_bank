package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	DB_DRIVER = "postgres"
	DB_SOURCE = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	// 建立数据库连接
	testDB, err = sql.Open(DB_DRIVER, DB_SOURCE)
	if err != nil {
		log.Fatal("Failed to open database connection", err)
	}

	// 暴露数据库连接接口
	testQueries = New(testDB)

	os.Exit(m.Run())
}
