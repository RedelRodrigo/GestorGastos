package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Registrar(in RegistroInput) (Usuario, error)
	Login(in LoginInput) (Usuario, error)
}

type sqliteAuthRepository struct {
	db *sql.DB
}

func NewSQLiteAuthRepository(db *sql.DB) Repository {
	return &sqliteAuthRepository{db: db}
}

func (r *sqliteAuthRepository) Registrar(in RegistroInput) (Usuario, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	username := strings.TrimSpace(in.Username)

	// Generar hash seguro de la contraseña
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return Usuario{}, fmt.Errorf("error al cifrar contraseña: %w", err)
	}

	query := "INSERT INTO usuarios (email, username, password_hash) VALUES (?, ?, ?)"
	res, err := r.db.Exec(query, email, username, string(hashBytes))
	if err != nil {
		// Verificamos si el error es por duplicidad de email (UNIQUE constraint)
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "unique") {
			return Usuario{}, ErrEmailDuplicado
		}
		return Usuario{}, fmt.Errorf("error al registrar usuario: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Usuario{}, fmt.Errorf("error al obtener ID: %w", err)
	}

	return Usuario{
		ID:       int(id),
		Email:    email,
		Username: username,
	}, nil
}

func (r *sqliteAuthRepository) Login(in LoginInput) (Usuario, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	var id int
	var username, passwordHash string

	query := "SELECT id, username, password_hash FROM usuarios WHERE email = ?"
	err := r.db.QueryRow(query, email).Scan(&id, &username, &passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Usuario{}, ErrCredencialesInv
		}
		return Usuario{}, fmt.Errorf("error al consultar usuario: %w", err)
	}

	// Validar la contraseña contra el hash almacenado
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(in.Password))
	if err != nil {
		return Usuario{}, ErrCredencialesInv
	}

	return Usuario{
		ID:       id,
		Email:    email,
		Username: username,
	}, nil
}
