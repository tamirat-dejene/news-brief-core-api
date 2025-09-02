package middleware

import (
	"net/http"
	"strings"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare(jwtService contract.IJWTService, userUseCase contract.IUserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower((parts[0])) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		tokenString := parts[1]

		claims, err := jwtService.ParseAccessToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		ctx.Set("userID", claims.UserID)
		// fixed failed type assertion, thus passed in the string version of claims.Role. it wasnt working on source_handler.go in the create source
		ctx.Set("userRole", string(claims.Role))

		ctx.Next()
	}
}
