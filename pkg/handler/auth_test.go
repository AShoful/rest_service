package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"rest/models"
	"rest/pkg/service"
	mock_service "rest/pkg/service/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuthorization, user models.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "name": "Test Name", "password": "qwerty"}`,
			inputUser: models.User{
				Username: "username",
				Name:     "Test Name",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(user).Return(uint(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"username": "username"}`,
			inputUser:            models.User{},
			mockBehavior:         func(r *mock_service.MockAuthorization, user models.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Invalid input"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "username", "name": "Test Name", "password": "qwerty"}`,
			inputUser: models.User{
				Username: "username",
				Name:     "Test Name",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(user).Return(uint(0), errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockAuthorization(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{Authorization: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_service.NewMockAuthorization(ctrl)
	services := &service.Service{Authorization: mockAuth}
	h := Handler{services}

	tests := []struct {
		name               string
		inputBody          string
		mockBehavior       func()
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"user1","password":"password"}`,
			mockBehavior: func() {
				mockAuth.EXPECT().GenerateToken("user1", "password").Return("mocked_token", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"token":"mocked_token"}`,
		},
		{
			name:               "Invalid input",
			inputBody:          `{"user1"}`,
			mockBehavior:       func() {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid input"}`,
		},
		{
			name:      "Service error",
			inputBody: `{"username":"user1","password":"password"}`,
			mockBehavior: func() {
				mockAuth.EXPECT().GenerateToken("user1", "password").Return("", errors.New("auth failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"auth failed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/sign-in", h.signIn)

			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
