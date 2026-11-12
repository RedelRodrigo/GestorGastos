package gastos

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"GestorGastos/internal/db"

	"github.com/gin-gonic/gin"
)

// setupRouter inicializa la base de datos en memoria y configura las rutas para testing.
func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// Inicia la BD SQLite en memoria
	database := db.InitDB(":memory:")
	t.Cleanup(func() { database.Close() })

	r := gin.New()
	RouterGasto(r, database)
	return r
}

// doRequest ayuda a simular peticiones HTTP hacia el router de Gin.
func doRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

const gastoOK = `{"monto": 1500.5, "categoria": "Comida", "fecha": "2026-11-05", "descripcion": "Almuerzo"}`

func TestCrearGasto(t *testing.T) {
	r := setupRouter(t)
	if w := doRequest(r, "POST", "/gastos", gastoOK); w.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuve %d: %s", w.Code, w.Body)
	}
}

func TestCrearGastoInvalido(t *testing.T) {
	r := setupRouter(t)
	malo := `{"monto": -5, "categoria": "", "fecha": "ayer"}`
	if w := doRequest(r, "POST", "/gastos", malo); w.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, obtuve %d: %s", w.Code, w.Body)
	}
}

func TestListarVacioDevuelveArreglo(t *testing.T) {
	r := setupRouter(t)
	w := doRequest(r, "GET", "/gastos", "")
	if !strings.Contains(w.Body.String(), `"data":[]`) && !strings.Contains(w.Body.String(), `"data": []`) {
		t.Fatalf("esperaba data: [], obtuve %s", w.Body)
	}
}

func TestEditarYBorrarInexistente(t *testing.T) {
	r := setupRouter(t)
	if w := doRequest(r, "PUT", "/gastos/999", gastoOK); w.Code != http.StatusNotFound {
		t.Errorf("PUT: esperaba 404, obtuve %d: %s", w.Code, w.Body)
	}

	if w := doRequest(r, "DELETE", "/gastos/999", ""); w.Code != http.StatusNotFound {
		t.Errorf("DELETE: esperaba 404, obtuve %d: %s", w.Code, w.Body)
	}
}
func TestFlujoCompletoGasto(t *testing.T) {
	r := setupRouter(t)

	// 1. Crear gasto
	w := doRequest(r, "POST", "/gastos", gastoOK)
	if w.Code != http.StatusCreated {
		t.Fatalf("POST falló: esperaba 201, obtuve %d", w.Code)
	}

	// 2. Listar y verificar que existe
	w = doRequest(r, "GET", "/gastos", "")
	if !strings.Contains(w.Body.String(), "Almuerzo") {
		t.Fatalf("GET falló: no encontró el gasto creado")
	}

	// 3. Editar el gasto con ID 1
	gastoEditado := `{"monto": 2000.0, "categoria": "Comida", "fecha": "2026-11-05", "descripcion": "Almuerzo Ejecutivo"}`
	w = doRequest(r, "PUT", "/gastos/1", gastoEditado)
	if w.Code != http.StatusOK {
		t.Fatalf("PUT falló: esperaba 200, obtuve %d", w.Code)
	}

	// 4. Borrar el gasto con ID 1
	w = doRequest(r, "DELETE", "/gastos/1", "")
	if w.Code != http.StatusOK {
		t.Fatalf("DELETE falló: esperaba 200, obtuve %d", w.Code)
	}

	// 5. Verificar que la lista vuelve a estar vacía
	w = doRequest(r, "GET", "/gastos", "")
	if !strings.Contains(w.Body.String(), `"data":[]`) && !strings.Contains(w.Body.String(), `"data": []`) {
		t.Fatalf("GET final falló: el arreglo debería estar vacío")
	}
}
func TestFiltroCategoriaYResumen(t *testing.T) {
	r := setupRouter(t)

	// Crear gastos de distintas categorías
	gasto1 := `{"monto": 100, "categoria": "Comida", "fecha": "2026-01-10", "descripcion": "Cena"}`
	gasto2 := `{"monto": 300, "categoria": "Comida", "fecha": "2026-01-15", "descripcion": "Almuerzo"}`
	gasto3 := `{"monto": 500, "categoria": "Transporte", "fecha": "2026-01-20", "descripcion": "Uber"}`

	doRequest(r, "POST", "/gastos", gasto1)
	doRequest(r, "POST", "/gastos", gasto2)
	doRequest(r, "POST", "/gastos", gasto3)

	// Probar filtro por categoría
	w := doRequest(r, "GET", "/gastos?categoria=Comida", "")
	if !strings.Contains(w.Body.String(), "Cena") || !strings.Contains(w.Body.String(), "Uber") == false {
		// Verificamos que contenga comida y no el uber
		if strings.Contains(w.Body.String(), "Uber") {
			t.Fatalf("El filtro por categoría trajo elementos no deseados: %s", w.Body.String())
		}
	}

	// Probar endpoint de resumen
	wResumen := doRequest(r, "GET", "/gastos/resumen", "")
	if wResumen.Code != http.StatusOK {
		t.Fatalf("Resumen falló: esperaba 200, obtuve %d", wResumen.Code)
	}
	if !strings.Contains(wResumen.Body.String(), `"total":900`) {
		t.Fatalf("El total del resumen es incorrecto, obtuve: %s", wResumen.Body.String())
	}
}
