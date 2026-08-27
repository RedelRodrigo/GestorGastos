package db

import (
	"testing"

	_ "modernc.org/sqlite"
)

func TestInitDB(t *testing.T) {
	db := InitDB()
	defer db.Close()

	if db == nil {
		t.Fatal("db should not be nil")
	}
}
