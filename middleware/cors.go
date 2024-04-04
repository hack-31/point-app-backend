package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware は、CORSの設定を行う
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost*",
			"https://*.dkjrwfcbom7qp.amplifyapp.com",
			"https://hack-31.github.io",
			"https://192.168.*",
			"http://192.168.*",
			"http://kaki-local.*",
		},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowWildcard:    true,
		AllowCredentials: true,
	})
}
