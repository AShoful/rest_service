package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const authorizationHeader = "Authorization"

func (h *Handler) userIdentity(c *gin.Context) {
	authHeader := c.GetHeader(authorizationHeader)

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	userId, err := h.services.Authorization.ParseToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.Set(userCtx, userId)
}
