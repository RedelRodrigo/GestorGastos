package gastos

/*
	type Gasto struct {
		ID          int     `json:"id"`
		Monto       float64 `json:"monto" binding:"required"`
		Categoria   string  `json:"categoria" binding:"required"`
		Fecha       string  `json:"fecha" binding:"required,datetime=02/01/2006"`
		Descripcion string  `json:"descripcion"`
	}
*/
type Gasto struct {
	ID          int     `json:"id"`
	Monto       float64 `json:"monto" binding:"required,gt=0"`
	Categoria   string  `json:"categoria" binding:"required"`
	Fecha       string  `json:"fecha" binding:"required,datetime=02/01/2006"`
	Descripcion string  `json:"descripcion"`
}
