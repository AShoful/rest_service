package handler

import (
	"net/http"
	"net/http/httptest"
	"rest/pkg/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestHandler_InitRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewHandler(&service.Service{})
	router := h.InitRoutes()

	tests := []struct {
		name       string
		method     string
		path       string
		statusCode int
	}{
		{"SignUp route exists", "POST", "/auth/sign-up", http.StatusBadRequest},
		{"SignIn route exists", "POST", "/auth/sign-in", http.StatusBadRequest},

		{"Create book exists", "POST", "/books/", http.StatusUnauthorized},
		{"Get all books exists", "GET", "/books/", http.StatusUnauthorized},
		{"Get book by id exists", "GET", "/books/1", http.StatusUnauthorized},
		{"Update book exists", "PUT", "/books/1", http.StatusUnauthorized},
		{"Delete book exists", "DELETE", "/books/1", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
