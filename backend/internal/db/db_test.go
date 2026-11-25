package db

import (
	"path/filepath"
	"testing"
)

func TestAbrir_CreaEsquema(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	database, err := Abrir("file:" + dbPath)
	if err != nil {
		t.Fatalf("Abrir devolvió error: %v", err)
	}
	defer database.Close()

	var cantidad int
	if err := database.QueryRow("SELECT COUNT(*) FROM gastos").Scan(&cantidad); err != nil {
		t.Fatalf("la tabla gastos no existe: %v", err)
	}
	if cantidad != 0 {
		t.Errorf("esperaba 0 filas, obtuve %d", cantidad)
	}
}
