package security

import (
	"sync"
	"time"
)

type registro struct {
	cuenta        int
	ventanaInicio time.Time
	bloqueoHasta  time.Time
}

// LimitadorLogin cuenta fallos por clave (IP o correo) y bloquea temporalmente.
// Vive en memoria: sirve para una sola instancia. Con varias, usa Redis.
type LimitadorLogin struct {
	mu          sync.Mutex
	registros   map[string]*registro
	maxIntentos int
	ventana     time.Duration
	bloqueo     time.Duration
}

const maxRegistros = 10_000

func NewLimitadorLogin(maxIntentos int, ventana, bloqueo time.Duration) *LimitadorLogin {
	return &LimitadorLogin{
		registros:   make(map[string]*registro),
		maxIntentos: maxIntentos,
		ventana:     ventana,
		bloqueo:     bloqueo,
	}
}

func (l *LimitadorLogin) Permitir(clave string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	r, ok := l.registros[clave]
	return !ok || !time.Now().Before(r.bloqueoHasta)
}

func (l *LimitadorLogin) Fallo(clave string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	ahora := time.Now()
	r, ok := l.registros[clave]
	if !ok || ahora.Sub(r.ventanaInicio) > l.ventana {
		r = &registro{ventanaInicio: ahora}
		l.registros[clave] = r
	}
	r.cuenta++
	if r.cuenta >= l.maxIntentos {
		r.bloqueoHasta = ahora.Add(l.bloqueo)
		r.cuenta = 0
		r.ventanaInicio = ahora
	}
	if len(l.registros) > maxRegistros {
		l.podar(ahora)
	}
}

func (l *LimitadorLogin) Exito(clave string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.registros, clave)
}

// podar elimina registros que ya no bloquean ni están dentro de la ventana.
func (l *LimitadorLogin) podar(ahora time.Time) {
	for k, r := range l.registros {
		if !ahora.Before(r.bloqueoHasta) && ahora.Sub(r.ventanaInicio) > l.ventana {
			delete(l.registros, k)
		}
	}
}
