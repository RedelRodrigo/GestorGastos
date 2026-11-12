package db

import (
	"database/sql"
	"log"
	"os"

	// Driver oficial de Turso (libsql) puro en Go, sin CGO
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func InitDB() *sql.DB {
	url := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	var connStr string

	// Si no hay variables de Turso, usa SQLite local (útil para desarrollo en tu PC)
	if url == "" {
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "gastos.db"
		}
		connStr = "file:" + dbPath
	} else {
		// Conexión remota a Turso en la nube
		connStr = url
		if authToken != "" {
			connStr = url + "?authToken=" + authToken
		}
	}

	// Abrimos la conexión usando el driver "libsql"
	db, err := sql.Open("libsql", connStr)
	if err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
	}

	if url == "" {
		db.SetMaxOpenConns(1)
	}

	// Esquema inicial de tablas
	const esquema = `
	CREATE TABLE IF NOT EXISTS usuarios (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		email         TEXT UNIQUE NOT NULL,
		username      TEXT NOT NULL,
		password_hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS gastos (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		usuario_id  TEXT NOT NULL,
		monto       REAL NOT NULL,
		categoria   TEXT NOT NULL,
		fecha       TEXT NOT NULL,
		descripcion TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_gastos_usuario ON gastos(usuario_id);
	CREATE INDEX IF NOT EXISTS idx_gastos_fecha ON gastos(fecha);
	`

	if _, err := db.Exec(esquema); err != nil {
		log.Fatal("Error al crear el esquema de la base de datos:", err)
	}

	log.Println("Base de datos inicializada correctamente.")
	return db
}
