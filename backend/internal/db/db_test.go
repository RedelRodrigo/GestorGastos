package db

import "testing"

func TestAbrir_CreaEsquema(t *testing.T) {
	database, err := Abrir(":memory:")
	if err != nil {
		t.Fatalf("Abrir devolvió error: %v", err)
	}
	defer database.Close()

	// Si la tabla existe, esta consulta funciona aunque no haya filas.
	var cantidad int
	if err := database.QueryRow("SELECT COUNT(*) FROM gastos").Scan(&cantidad); err != nil {
		t.Fatalf("la tabla gastos no existe: %v", err)
	}
	if cantidad != 0 {
		t.Errorf("esperaba 0 filas, obtuve %d", cantidad)
	}
}
