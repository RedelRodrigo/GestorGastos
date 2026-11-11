// Package db se encarga de abrir la conexión a SQLite y de preparar el esquema.
package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // registra el driver "sqlite"
)

const esquema = `
CREATE TABLE IF NOT EXISTS gastos (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	monto       REAL NOT NULL,
	categoria   TEXT NOT NULL,
	fecha       TEXT NOT NULL, -- formato ISO: AAAA-MM-DD
	descripcion TEXT
);
CREATE INDEX IF NOT EXISTS idx_gastos_fecha ON gastos(fecha);
`

// Abrir abre (o crea) la base de datos en path y garantiza que el esquema exista.
// Usar ":memory:" como path crea una base temporal, útil para tests.
func Abrir(path string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("abrir base de datos: %w", err)
	}

	// Cada conexión a ":memory:" sería una base distinta; con una sola se comparte.
	if path == ":memory:" {
		database.SetMaxOpenConns(1)
	}

	if _, err := database.Exec(esquema); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("crear esquema: %w", err)
	}
	return database, nil
}
