package gastos

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockDB struct {
	queries map[string][]map[string]interface{}
}

func (m *mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if _, exists := m.queries[query]; !exists {
		return nil, nil
	}
	// Return a mock result
	result := &mockResult{}
	return result, nil
}

func (m *mockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if rows, exists := m.queries[query]; exists {
		return &mockRows{rows: rows}, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockDB) QueryRow(query string, args ...interface{}) sql.Row {
	return &mockRow{}
}

type mockResult struct {
	rowsAffected int64
}

func (m *mockResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (m *mockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

type mockRows struct {
	rows []map[string]interface{}
	idx  int
}

func (m *mockRows) Next() bool {
	if m.idx < len(m.rows) {
		m.idx++
		return true
	}
	return false
}

func (m *mockRows) Close() error {
	return nil
}

func (m *mockRows) Err() error {
	return nil
}

func (m *mockRows) Scan(dest ...interface{}) error {
	if m.idx-1 >= len(m.rows) {
		return sql.ErrNoRows
	}
	row := m.rows[m.idx-1]
	for i, d := range dest {
		if val, ok := row[i]; ok {
			dest[i] = val
		}
	}
	return nil
}

type mockRow struct {
	scanVal interface{}
}

func (m *mockRow) Scan(dest ...interface{}) error {
	if len(dest) > 0 {
		dest[0] = m.scanVal
	}
	return nil
}

func setupMockDB() *mockDB {
	mdb := &mockDB{
		queries: make(map[string][]map[string]interface{}),
	}
	// Setup mock data for obtenerGastosHandler
	mdb.queries = map[string][]map[string]interface{}{
		"SELECT id, monto, categoria, fecha, descripcion FROM gastos": {
			{"1", "100.50", "Alimentación", "2024-01-15", "Comida"},
		},
	}
	return mdb
}

func TestCrearGastoHandler(t *testing.T) {
	db := setupMockDB()
	handler := crearGastosHandler(db)

	c := gin.New()
	c.Request, _ = http.NewRequest("POST", "/gastos", `{"monto":50.0,"categoria":"Transporte","fecha":"2024-01-20","descripcion":"Combustible"}`)
	w := gin.NewResponseWriter()

	handler(w, c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status Created, got %d", w.Code)
	}
}

func TestObtenerGastosHandler(t *testing.T) {
	db := setupMockDB()
	handler := obtenerGastosHandler(db)

	c := gin.New()
	w := gin.NewResponseWriter()

	handler(w, c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %d", w.Code)
	}
}

func TestEditarGastosHandler(t *testing.T) {
	db := setupMockDB()
	handler := editarGastosHandler(db)

	c := gin.New()
	c.Request, _ = http.NewRequest("PUT", "/gastos/1", `{"monto":75.0,"categoria":"Alimentación","fecha":"2024-01-20","descripcion":"Comida actualizada"}`)
	w := gin.NewResponseWriter()

	handler(w, c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %d", w.Code)
	}
}

func TestBorrarGastoHandler(t *testing.T) {
	db := setupMockDB()
	handler := borrarGastoHandler(db)

	c := gin.New()
	w := gin.NewResponseWriter()

	handler(w, c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %d", w.Code)
	}
}
