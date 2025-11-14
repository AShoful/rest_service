package handler

import (
	"bytes"
	"net/http/httptest"
	"rest/models"
	"rest/pkg/service"
	mock_service "rest/pkg/service/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_create(t *testing.T) {
	type mockBehavior func(r *mock_service.MockBook, book models.Book)

	tests := []struct {
		name                 string
		inputBody            string
		inputBook            models.Book
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"title": "title", "author": "author"}`,
			inputBook: models.Book{
				Title:  "title",
				Author: "author",
				UserId: 1,
			},
			mockBehavior: func(r *mock_service.MockBook, inputBook models.Book) {
				r.EXPECT().Create(inputBook).Return(uint(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"bookId":1}`,
		},
		// {
		// 	name:                 "Wrong Input",
		// 	inputBody:            `{"username": "username"}`,
		// 	inputUser:            models.User{},
		// 	mockBehavior:         func(r *mock_service.MockAuthorization, user models.User) {},
		// 	expectedStatusCode:   400,
		// 	expectedResponseBody: `{"error":"Invalid input"}`,
		// },
		// {
		// 	name:      "Service Error",
		// 	inputBody: `{"username": "username", "name": "Test Name", "password": "qwerty"}`,
		// 	inputUser: models.User{
		// 		Username: "username",
		// 		Name:     "Test Name",
		// 		Password: "qwerty",
		// 	},
		// 	mockBehavior: func(r *mock_service.MockAuthorization, user models.User) {
		// 		r.EXPECT().CreateUser(user).Return(uint(0), errors.New("something went wrong"))
		// 	},
		// 	expectedStatusCode:   500,
		// 	expectedResponseBody: `{"error":"something went wrong"}`,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockBook(c)
			test.mockBehavior(repo, test.inputBook)

			services := &service.Service{Book: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/", handler.createBook)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/books/",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

// func TestHandler_signIn(t *testing.T) {
// 	type mockBehavior func(r *mock_service.MockAuthorization, username, password string)

// 	tests := []struct {
// 		name                 string
// 		inputBody            string
// 		inputUsername        string
// 		inputPassword        string
// 		mockBehavior         mockBehavior
// 		expectedStatusCode   int
// 		expectedResponseBody string
// 	}{
// 		{
// 			name:          "Ok",
// 			inputBody:     `{"username":"user1","password":"pass123"}`,
// 			inputUsername: "user1",
// 			inputPassword: "pass123",
// 			mockBehavior: func(r *mock_service.MockAuthorization, username, password string) {
// 				r.EXPECT().GenerateToken(username, password).Return("mocked_token", nil)
// 			},
// 			expectedStatusCode:   http.StatusOK,
// 			expectedResponseBody: `{"token":"mocked_token"}`,
// 		},
// 		{
// 			name:                 "Invalid input",
// 			inputBody:            `{"user1"}`,
// 			inputUsername:        "user1",
// 			inputPassword:        "",
// 			mockBehavior:         func(r *mock_service.MockAuthorization, username, password string) {},
// 			expectedStatusCode:   http.StatusBadRequest,
// 			expectedResponseBody: `{"error":"Invalid input"}`,
// 		},
// 		{
// 			name:          "Service error",
// 			inputBody:     `{"username":"user1","password":"pass123"}`,
// 			inputUsername: "user1",
// 			inputPassword: "pass123",
// 			mockBehavior: func(r *mock_service.MockAuthorization, username, password string) {
// 				r.EXPECT().GenerateToken("user1", "pass123").Return("", errors.New("auth failed"))
// 			},
// 			expectedStatusCode:   http.StatusInternalServerError,
// 			expectedResponseBody: `{"error":"auth failed"}`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			authMock := mock_service.NewMockAuthorization(ctrl)
// 			tt.mockBehavior(authMock, tt.inputUsername, tt.inputPassword)

// 			services := &service.Service{Authorization: authMock}
// 			handler := Handler{services}

// 			r := gin.New()
// 			r.POST("/sign-in", handler.signIn)

// 			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(tt.inputBody))
// 			req.Header.Set("Content-Type", "application/json")
// 			w := httptest.NewRecorder()

// 			r.ServeHTTP(w, req)

// 			assert.Equal(t, w.Code, tt.expectedStatusCode)
// 			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
// 		})
// 	}
// }
