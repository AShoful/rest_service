package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"rest/pkg/service"
	mock_service "rest/pkg/service/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_userIdentity_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	mockAuth.EXPECT().ParseToken("valid-token").Return(uint(42), nil)

	services := &service.Service{Authorization: mockAuth}
	h := Handler{services}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", h.userIdentity, func(c *gin.Context) {
		id, _ := c.Get(userCtx)
		c.JSON(http.StatusOK, gin.H{"userId": id})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"userId":42}`, w.Body.String())
}

func TestHandler_userIdentity_MissingHeader(t *testing.T) {
	services := &service.Service{}
	h := Handler{services}

	r := gin.New()
	r.GET("/test", h.userIdentity, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, `{"error":"Missing or invalid Authorization header"}`, w.Body.String())
}

func TestHandler_userIdentity_InvalidPrefix(t *testing.T) {
	services := &service.Service{}
	h := Handler{services}

	r := gin.New()
	r.GET("/test", h.userIdentity, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Token abcdef")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, `{"error":"Missing or invalid Authorization header"}`, w.Body.String())
}

func TestHandler_userIdentity_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	mockAuth.EXPECT().ParseToken("bad-token").Return(uint(0), errors.New("token invalid"))

	services := &service.Service{Authorization: mockAuth}
	h := Handler{services}

	r := gin.New()
	r.GET("/test", h.userIdentity, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer bad-token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, `{"error":"Invalid or expired token"}`, w.Body.String())
}
