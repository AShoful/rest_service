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
		// {
		// 	name:      "Ok",
		// 	inputBody: `{"title": "title", "author": "author"}`,
		// 	inputBook: models.Book{
		// 		Title:  "title",
		// 		Author: "author",
		// 		UserId: uint(1),
		// 	},
		// 	mockBehavior: func(r *mock_service.MockBook, inputBook models.Book) {
		// 		r.EXPECT().Create(inputBook).Return(uint(1), nil)
		// 	},
		// 	expectedStatusCode:   200,
		// 	expectedResponseBody: `{"bookId":1}`,
		// },
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
