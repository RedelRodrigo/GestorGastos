package gastos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

package gastos

import (
"context"
"database/sql"
"errors"
"fmt"
"strings"
)

// Repositorio define qué operaciones de persistencia necesita el módulo de gastos.
// El handler depende de esta interfaz y no de SQLite: si mañana se cambia la base
// de datos (o se necesita un falso para tests), solo hay que escribir otra implementación.
type Repositorio interface {
	Crear(ctx context.Context, in GastoInput) (Gasto, error)
	Listar(ctx context.Context, f Filtro) ([]Gasto, error)
	Obtener(ctx context.Context, id int64) (Gasto, error)
	Actualizar(ctx context.Context, id int64, in GastoInput) (Gasto, error)
	Borrar(ctx context.Context, id int64) error
	Resumen(ctx context.Context) ([]TotalCategoria, error)
}

type repoSQL struct {
	db *sql.DB
}

// NuevoRepositorio devuelve un Repositorio respaldado por una base SQL.
func NuevoRepositorio(db *sql.DB) Repositorio {
	return &repoSQL{db: db}
}

// columnas es la lista que se usa en todos los SELECT, en el mismo orden que escanear.
// COALESCE evita errores de Scan si descripcion es NULL en filas antiguas.
const columnas = "id, monto, categoria, fecha, COALESCE(descripcion, '')"

// escaner lo cumplen tanto *sql.Row como *sql.Rows, así escanear sirve para ambos.
type escaner interface {
	Scan(dest ...any) error
}

func escanear(s escaner) (Gasto, error) {
	var g Gasto
	err := s.Scan(&g.ID, &g.Monto, &g.Categoria, &g.Fecha, &g.Descripcion)
	return g, err
}

func (r *repoSQL) Crear(ctx context.Context, in GastoInput) (Gasto, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO gastos (monto, categoria, fecha, descripcion) VALUES (?, ?, ?, ?)",
		in.Monto, in.Categoria, in.Fecha, in.Descripcion)
	if err != nil {
		return Gasto{}, fmt.Errorf("crear gasto: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Gasto{}, fmt.Errorf("crear gasto: obtener id: %w", err)
	}
	return in.aGasto(id), nil
}

func (r *repoSQL) Listar(ctx context.Context, f Filtro) ([]Gasto, error) {
	query := "SELECT " + columnas + " FROM gastos"

	// Los valores SIEMPRE viajan como argumentos (?), nunca concatenados al SQL.
	var condiciones []string
	var args []any
	if f.Categoria != "" {
		condiciones = append(condiciones, "categoria = ?")
		args = append(args, f.Categoria)
	}
	if f.Desde != "" {
		condiciones = append(condiciones, "fecha >= ?")
		args = append(args, f.Desde)
	}
	if f.Hasta != "" {
		condiciones = append(condiciones, "fecha <= ?")
		args = append(args, f.Hasta)
	}
	if len(condiciones) > 0 {
		query += " WHERE " + strings.Join(condiciones, " AND ")
	}
	query += " ORDER BY fecha DESC, id DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listar gastos: %w", err)
	}
	defer rows.Close()

	gastos := []Gasto{} // vacío y no nil, para que el JSON sea [] y no null
	for rows.Next() {
		g, err := escanear(rows)
		if err != nil {
			return nil, fmt.Errorf("listar gastos: %w", err)
		}
		gastos = append(gastos, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listar gastos: %w", err)
	}
	return gastos, nil
}

func (r *repoSQL) Obtener(ctx context.Context, id int64) (Gasto, error) {
	g, err := escanear(r.db.QueryRowContext(ctx,
		"SELECT "+columnas+" FROM gastos WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return Gasto{}, ErrNoEncontrado
	}
	if err != nil {
		return Gasto{}, fmt.Errorf("obtener gasto %d: %w", id, err)
	}
	return g, nil
}

func (r *repoSQL) Actualizar(ctx context.Context, id int64, in GastoInput) (Gasto, error) {
	res, err := r.db.ExecContext(ctx,
		"UPDATE gastos SET monto = ?, categoria = ?, fecha = ?, descripcion = ? WHERE id = ?",
		in.Monto, in.Categoria, in.Fecha, in.Descripcion, id)
	if err != nil {
		return Gasto{}, fmt.Errorf("actualizar gasto %d: %w", id, err)
	}
	if err := verificarAfectadas(res, id); err != nil {
		return Gasto{}, err
	}
	return in.aGasto(id), nil
}

func (r *repoSQL) Borrar(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM gastos WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("borrar gasto %d: %w", id, err)
	}
	return verificarAfectadas(res, id)
}

func (r *repoSQL) Resumen(ctx context.Context) ([]TotalCategoria, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT categoria, SUM(monto) FROM gastos GROUP BY categoria ORDER BY SUM(monto) DESC")
	if err != nil {
		return nil, fmt.Errorf("resumen: %w", err)
	}
	defer rows.Close()

	totales := []TotalCategoria{}
	for rows.Next() {
		var t TotalCategoria
		if err := rows.Scan(&t.Categoria, &t.Total); err != nil {
			return nil, fmt.Errorf("resumen: %w", err)
		}
		totales = append(totales, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("resumen: %w", err)
	}
	return totales, nil
}

// verificarAfectadas devuelve ErrNoEncontrado si la operación no tocó ninguna fila.
func verificarAfectadas(res sql.Result, id int64) error {
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("gasto %d: filas afectadas: %w", id, err)
	}
	if n == 0 {
		return ErrNoEncontrado
	}
	return nil
}