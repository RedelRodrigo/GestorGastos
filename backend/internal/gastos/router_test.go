package gastos

import (
	"testing"
)

func TestGastoStruct(t *testing.T) {
	g := Gasto{
		ID:        1,
		Monto:     50.0,
		Fecha:     "2024-01-01",
		Categoria: "Test",
	}

	if g.ID != 1 {
		t.Errorf("expected ID 1, got %d", g.ID)
	}
	if g.Monto != 50.0 {
		t.Errorf("expected Monto 50, got %f", g.Monto)
	}
	if g.Categoria != "Test" {
		t.Errorf("expected Categoria Test, got %s", g.Categoria)
	}
}
