// Package db se encarga de abrir la conexión a SQLite y de preparar el esquema.
package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // registra el driver "sqlite"
)

// InitDB inicializa la base de datos SQLite en la ruta especificada y garantiza que el esquema exista.
// Usar ":memory:" como path crea una base temporal, útil para tests.
func InitDB(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal("Error al abrir la base de datos:", err)
	}

	// Cada conexión a ":memory:" sería una base distinta; con una sola conexión se comparte.
	if path == ":memory:" {
		db.SetMaxOpenConns(1)
	}

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

	if _, err := db.Exec(esquema); err != nil {
		log.Fatal("Error al crear el esquema de la base de datos:", err)
	}

	return db
}
