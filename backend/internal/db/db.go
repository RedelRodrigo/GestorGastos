package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

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

func Abrir(connStr string) (*sql.DB, error) {
	driver := "libsql"
	dsn := connStr
	if strings.HasPrefix(connStr, "file:") {
		driver = "sqlite"
		dsn = strings.TrimPrefix(connStr, "file:")
	}

	db, err := sql.Open(driver, dsn)
	// db, err := sql.Open("libsql", connStr)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir la base de datos: %w", err)
	}

	if _, err := db.Exec(esquema); err != nil {
		db.Close()
		return nil, fmt.Errorf("no se pudo crear el esquema de la base de datos: %w", err)
	}

	return db, nil
}

func InitDB() *sql.DB {
	url := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	var connStr string
	if url == "" {
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "gastos.db"
		}
		connStr = "file:" + dbPath
	} else {
		connStr = url
		if authToken != "" {
			connStr = url + "?authToken=" + authToken
		}
	}

	// Nunca registrar connStr: puede contener el token de Turso.
	database, err := Abrir(connStr)
	if err != nil {
		slog.Error("error inicializando base de datos", "error", err)
		os.Exit(1)
	}

	if url == "" {
		database.SetMaxOpenConns(1)
	}

	slog.Info("base de datos inicializada correctamente")
	return database
}
