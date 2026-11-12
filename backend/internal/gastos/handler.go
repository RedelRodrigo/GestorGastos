package gastos

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func crearGastosHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var gasto Gasto

		if err := c.ShouldBindJSON(&gasto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := db.Exec(
			"INSERT INTO gastos (monto, categoria, fecha, descripcion) VALUES (?, ?, ?, ?)",
			gasto.Monto, gasto.Categoria, gasto.Fecha, gasto.Descripcion)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		id, _ := result.LastInsertId()
		gasto.ID = int(id)

		c.JSON(http.StatusCreated, gin.H{"data": gasto})
	}
}

func obtenerGastosHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := "SELECT id, monto, categoria, fecha, COALESCE(descripcion, '') FROM gastos WHERE 1=1"
		args := []any{}

		if v := c.Query("categoria"); v != "" {
			query += " AND categoria = ?"
			args = append(args, v)
		}
		if v := c.Query("desde"); v != "" {
			query += " AND fecha >= ?"
			args = append(args, v)
		}
		if v := c.Query("hasta"); v != "" {
			query += " AND fecha <= ?"
			args = append(args, v)
		}

		query += " ORDER BY fecha DESC, id DESC"

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		gastos := []Gasto{}
		for rows.Next() {
			var g Gasto
			if err := rows.Scan(&g.ID, &g.Monto, &g.Categoria, &g.Fecha, &g.Descripcion); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			gastos = append(gastos, g)
		}

		c.JSON(http.StatusOK, gin.H{"data": gastos})
	}
}

func editarGastosHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}

		var gasto Gasto
		if err := c.ShouldBindJSON(&gasto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := "UPDATE gastos SET monto = ?, categoria = ?, fecha = ?, descripcion = ? WHERE id = ?"
		res, err := db.Exec(query, gasto.Monto, gasto.Categoria, gasto.Fecha, gasto.Descripcion, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := res.RowsAffected()
		if err != nil || rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "gasto no encontrado"})
			return
		}

		gasto.ID = id
		c.JSON(http.StatusOK, gin.H{"data": gasto})
	}
}

func borrarGastoHandler(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result, err := database.Exec("DELETE FROM gastos WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		filas, _ := result.RowsAffected()
		if filas == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "gasto no encontrado"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "gasto eliminado"})
	}
}

func obtenerGastoHandlerId(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}

		result, err := database.Query("SELECT id, monto, categoria, fecha, descripcion FROM gastos WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer result.Close()

		if result.Next() {
			var gasto Gasto
			if err := result.Scan(&gasto.ID, &gasto.Monto, &gasto.Categoria, &gasto.Fecha, &gasto.Descripcion); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": gasto})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "gasto no encontrado"})

	}

}

func resumenHandler(db *sql.DB) gin.HandlerFunc {
	type item struct {
		Categoria string  `json:"categoria"`
		Total     float64 `json:"total"`
	}
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT categoria, SUM(monto) FROM gastos GROUP BY categoria ORDER BY SUM(monto) DESC")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		items := []item{}
		var total float64
		for rows.Next() {
			var it item
			if err := rows.Scan(&it.Categoria, &it.Total); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			total += it.Total
			items = append(items, it)
		}
		c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
	}
}
