package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

//func InitDB() *sql.DB {
//	// Abre (o crea si no existe) el archivo gastos.db
//	db, err := sql.Open("sqlite", "./gastos.db")
//	if err != nil {
//		log.Fatal("Error al abrir la base de datos:", err)
//	}

func InitDB(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal("Error al abrir la base de datos:", err)
	}
	if path == ":memory:" {
		db.SetMaxOpenConns(1)
	}

	// Creamos la tabla si todavía no existe
	query := `
	CREATE TABLE IF NOT EXISTS gastos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		monto REAL NOT NULL,
		categoria TEXT NOT NULL,
		fecha TEXT NOT NULL,
		descripcion TEXT
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Error al crear la tabla:", err)
	}

	return db
}
