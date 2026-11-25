package gastos

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"GestorGastos/internal/db"
	"GestorGastos/internal/security"

	"github.com/gin-gonic/gin"
)

func setupRouter(t *testing.T) (*gin.Engine, *http.Cookie) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "secreto-de-prueba-con-mas-de-32-caracteres-123456")
	if err := security.InicializarSecreto(); err != nil {
		t.Fatalf("no se pudo inicializar secreto: %v", err)
	}

	dbPath := filepath.Join(t.TempDir(), "gastos_test.db")
	database, err := db.Abrir("file:" + dbPath)
	if err != nil {
		t.Fatalf("no se pudo abrir BD de test: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	rec := httptest.NewRecorder()
	if err := security.EmitirSesion(rec, "test@example.com", "tester"); err != nil {
		t.Fatalf("no se pudo emitir sesión de test: %v", err)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("no se generó cookie de sesión")
	}

	r := gin.New()
	api := r.Group("/api")
	RouterGasto(api, database)
	return r, cookies[0]
}

func doRequest(r http.Handler, cookie *http.Cookie, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
