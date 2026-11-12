package gastos

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	Crear(usuarioID string, g Gasto) (int, error)
	Listar(usuarioID string, f Filtro) ([]Gasto, error)
	Actualizar(usuarioID string, id int, g Gasto) error
	Borrar(usuarioID string, id int) error
	ObtenerResumen(usuarioID string) ([]TotalCategoria, float64, error)
}

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) Repository {
	return &sqliteRepository{db: db}
}

func (r *sqliteRepository) Crear(usuarioID string, g Gasto) (int, error) {
	query := "INSERT INTO gastos (usuario_id, monto, categoria, fecha, descripcion) VALUES (?, ?, ?, ?, ?)"
	res, err := r.db.Exec(query, usuarioID, g.Monto, g.Categoria, g.Fecha, g.Descripcion)
	if err != nil {
		return 0, fmt.Errorf("error al insertar gasto: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener ID: %w", err)
	}
	return int(id), nil
}

func (r *sqliteRepository) Listar(usuarioID string, f Filtro) ([]Gasto, error) {
	query := "SELECT id, monto, categoria, fecha, COALESCE(descripcion, '') FROM gastos WHERE usuario_id = ?"
	args := []any{usuarioID}

	if f.Categoria != "" {
		query += " AND categoria = ?"
		args = append(args, f.Categoria)
	}
	if f.Desde != "" {
		query += " AND fecha >= ?"
		args = append(args, f.Desde)
	}
	if f.Hasta != "" {
		query += " AND fecha <= ?"
		args = append(args, f.Hasta)
	}

	query += " ORDER BY fecha DESC, id DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al listar gastos: %w", err)
	}
	defer rows.Close()

	gastos := []Gasto{}
	for rows.Next() {
		var g Gasto
		if err := rows.Scan(&g.ID, &g.Monto, &g.Categoria, &g.Fecha, &g.Descripcion); err != nil {
			return nil, fmt.Errorf("error al escanear gasto: %w", err)
		}
		gastos = append(gastos, g)
	}
	return gastos, nil
}

func (r *sqliteRepository) Actualizar(usuarioID string, id int, g Gasto) error {
	query := "UPDATE gastos SET monto = ?, categoria = ?, fecha = ?, descripcion = ? WHERE id = ? AND usuario_id = ?"
	res, err := r.db.Exec(query, g.Monto, g.Categoria, g.Fecha, g.Descripcion, id, usuarioID)
	if err != nil {
		return fmt.Errorf("error al actualizar: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNoEncontrado
	}
	return nil
}

func (r *sqliteRepository) Borrar(usuarioID string, id int) error {
	query := "DELETE FROM gastos WHERE id = ? AND usuario_id = ?"
	res, err := r.db.Exec(query, id, usuarioID)
	if err != nil {
		return fmt.Errorf("error al borrar: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNoEncontrado
	}
	return nil
}

func (r *sqliteRepository) ObtenerResumen(usuarioID string) ([]TotalCategoria, float64, error) {
	query := "SELECT categoria, SUM(monto) FROM gastos WHERE usuario_id = ? GROUP BY categoria ORDER BY SUM(monto) DESC"
	rows, err := r.db.Query(query, usuarioID)
	if err != nil {
		return nil, 0, fmt.Errorf("error al obtener resumen: %w", err)
	}
	defer rows.Close()

	items := []TotalCategoria{}
	var total float64
	for rows.Next() {
		var it TotalCategoria
		if err := rows.Scan(&it.Categoria, &it.Total); err != nil {
			return nil, 0, fmt.Errorf("error al escanear resumen: %w", err)
		}
		total += it.Total
		items = append(items, it)
	}
	return items, total, nil
}
