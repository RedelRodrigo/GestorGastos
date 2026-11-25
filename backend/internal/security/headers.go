package security

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func Headers() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")
		h.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' https://challenges.cloudflare.com; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data:; "+
				"connect-src 'self'; "+
				"frame-src https://challenges.cloudflare.com; "+
				"frame-ancestors 'none'; object-src 'none'; base-uri 'self'; form-action 'self'")

		// Los datos financieros nunca deben quedar en caché del navegador o de un proxy.
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			h.Set("Cache-Control", "no-store")
		}
		c.Next()
	}
}
