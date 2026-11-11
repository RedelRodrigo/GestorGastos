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
		// 1. Hacemos la consulta a la base de datos
		rows, err := db.Query("SELECT id, monto, categoria, fecha, descripcion FROM gastos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// ¡Muy importante! Asegurarse de cerrar las filas para liberar la conexión a la base de datos
		defer rows.Close()

		// 2. Creamos un slice (arreglo) para almacenar los gastos que vayamos leyendo
		var gastos []Gasto

		// 3. Iteramos fila por fila
		for rows.Next() {
			var g Gasto
			// 4. Escaneamos los valores de la fila actual hacia nuestra estructura
			// El orden en Scan debe coincidir exactamente con el orden del SELECT
			if err := rows.Scan(&g.ID, &g.Monto, &g.Categoria, &g.Fecha, &g.Descripcion); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			gastos = append(gastos, g)
		}

		// 5. Verificamos si ocurrió algún error durante la iteración del for
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Si el slice es nil (porque la tabla está vacía), lo inicializamos como vacío
		// para que el JSON devuelva [] en lugar de null
		if gastos == nil {
			gastos = []Gasto{}
		}

		// 6. Retornamos la lista de gastos
		c.JSON(http.StatusOK, gin.H{"data": gastos})
	}
}

func editarGastosHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var gasto Gasto

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}

		if err := c.ShouldBindJSON(&gasto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := db.Exec(
			"UPDATE gastos SET monto = ?, categoria = ?, fecha = ?, descripcion = ? WHERE id = ?", gasto.Monto, gasto.Categoria, gasto.Fecha, gasto.Descripcion, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "gasto no encontrado"})
			return
		}
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
		id := c.Param("id")

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
