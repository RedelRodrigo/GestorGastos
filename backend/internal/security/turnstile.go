package security

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const siteverifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

var clienteTurnstile = &http.Client{Timeout: 5 * time.Second}

// VerificarTurnstile valida el token con Cloudflare. Ante cualquier error devuelve false.
func VerificarTurnstile(ctx context.Context, token, ip string) bool {
	secret := os.Getenv("TURNSTILE_SECRET")
	if secret == "" || token == "" {
		return false
	}

	form := url.Values{"secret": {secret}, "response": {token}, "remoteip": {ip}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, siteverifyURL,
		strings.NewReader(form.Encode()))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := clienteTurnstile.Do(req)
	if err != nil {
		slog.WarnContext(ctx, "turnstile no respondió", "error", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.WarnContext(ctx, "turnstile status inesperado", "status", resp.StatusCode)
		return false
	}

	var r struct {
		Success    bool     `json:"success"`
		ErrorCodes []string `json:"error-codes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false
	}
	if !r.Success {
		slog.InfoContext(ctx, "turnstile rechazado", "codigos", r.ErrorCodes)
	}
	return r.Success
}
