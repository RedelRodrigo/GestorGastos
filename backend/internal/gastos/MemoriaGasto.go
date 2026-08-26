package gastos

type Gasto struct {
	ID          int     `json:"id"`
	Monto       float64 `json:"monto"`
	Categoria   string  `json:"categoria"`
	Fecha       string  `json:"fecha"`
	Descripcion string  `json:"descripcion"`
}
