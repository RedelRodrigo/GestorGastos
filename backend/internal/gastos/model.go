package gastos

import "errors"

// layoutFecha es el formato de fecha que usa toda la API: AAAA-MM-DD (ISO 8601).
// Al ser ISO, ordenar y comparar fechas como texto funciona correctamente.
const layoutFecha = "2006-01-02"

// ErrNoEncontrado se devuelve cuando un gasto con ese id no existe.
var ErrNoEncontrado = errors.New("gasto no encontrado")

// Gasto es la entidad tal como se guarda y se devuelve en las respuestas.
type Gasto struct {
	ID          int64   `json:"id"`
	Monto       float64 `json:"monto"`
	Categoria   string  `json:"categoria"`
	Fecha       string  `json:"fecha"`
	Descripcion string  `json:"descripcion"`
}

// GastoInput son los datos que envía el cliente al crear o editar un gasto.
// Está separado de Gasto para que el cliente nunca pueda fijar el id
// y para concentrar las reglas de validación en un solo lugar.
type GastoInput struct {
	Monto       float64 `json:"monto" binding:"required,gt=0"`
	Categoria   string  `json:"categoria" binding:"required"`
	Fecha       string  `json:"fecha" binding:"required,datetime=2006-01-02"` // debe coincidir con layoutFecha
	Descripcion string  `json:"descripcion"`
}

// aGasto convierte los datos de entrada en una entidad con el id indicado.
func (in GastoInput) aGasto(id int64) Gasto {
	return Gasto{
		ID:          id,
		Monto:       in.Monto,
		Categoria:   in.Categoria,
		Fecha:       in.Fecha,
		Descripcion: in.Descripcion,
	}
}

// Filtro permite acotar el listado de gastos. Los campos vacíos se ignoran.
type Filtro struct {
	Categoria string
	Desde     string // inclusive, AAAA-MM-DD
	Hasta     string // inclusive, AAAA-MM-DD
}

// TotalCategoria es una fila del resumen de gastos agrupado por categoría.
type TotalCategoria struct {
	Categoria string  `json:"categoria"`
	Total     float64 `json:"total"`
}
