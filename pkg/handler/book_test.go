package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"rest/models"
	"rest/pkg/service"
	mock_service "rest/pkg/service/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func setUser(ctx *gin.Context, id uint) {
	ctx.Set(userCtx, id)
}

func str(s string) *string {
	return &s
}

func TestHandler_createBook_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)

	input := models.Book{
		Title:  "GoLang Book",
		Author: "AuthorName",
		UserId: 1,
	}

	mockBook.EXPECT().Create(input).Return(uint(10), nil)

	services := &service.Service{Book: mockBook}
	h := Handler{services}

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/books", func(c *gin.Context) {
		setUser(c, 1)
		h.createBook(c)
	})

	body := `{"title":"GoLang Book","author":"AuthorName"}`

	req := httptest.NewRequest("POST", "/books", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"bookId":10}`, w.Body.String())
}

func TestHandler_createBook_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)

	services := &service.Service{Book: mockBook}
	h := Handler{services}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/books", h.createBook)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/books", bytes.NewBufferString(`{`))

	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"invalid input"}`, w.Body.String())
}

func TestHandler_getAllBooks_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)

	mockBook.EXPECT().GetAll().Return([]models.Book{
		{ID: 1, Title: "Book A", Author: "John"},
		{ID: 2, Title: "Book B", Author: "Mike"},
	}, nil)

	services := &service.Service{Book: mockBook}
	h := Handler{services}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/books", h.getAllBook)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/books", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t,
		`{"data":[{"id":1,"title":"Book A","author":"John","userid":0},{"id":2,"title":"Book B","author":"Mike","userid":0}]}`,
		w.Body.String(),
	)
}

func TestHandler_getAllBooks_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)

	mockBook.EXPECT().GetAll().Return(nil, errors.New("db error"))

	services := &service.Service{Book: mockBook}
	h := Handler{services}

	r := gin.New()
	r.GET("/books", h.getAllBook)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/books", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, `{"error":"db error"}`, w.Body.String())
}

func TestHandler_getBookById_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)

	mockBook.EXPECT().GetById(uint(5)).Return(models.Book{
		ID:     5,
		Title:  "Go",
		Author: "Bob",
	}, nil)

	services := &service.Service{Book: mockBook}
	h := Handler{services}

	r := gin.New()
	r.GET("/books/:id", h.getBookById)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/books/5", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"id":5,"title":"Go","author":"Bob","userid":0}`, w.Body.String())
}

func TestHandler_getBookById_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := &service.Service{Book: mock_service.NewMockBook(ctrl)}
	h := Handler{services}

	r := gin.New()
	r.GET("/books/:id", h.getBookById)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/books/xyz", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"invalid book id"}`, w.Body.String())
}

func TestHandler_getBookById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	mockBook.EXPECT().GetById(uint(77)).Return(models.Book{}, errors.New("not found"))

	services := &service.Service{Book: mockBook}
	h := Handler{services}

	r := gin.New()
	r.GET("/books/:id", h.getBookById)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/books/77", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, `{"error":"not found"}`, w.Body.String())
}

func TestHandler_updateBook_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)

	input := models.UpdateBook{
		Title:  str("New Title"),
		Author: str("New Author"),
	}

	mockBook.EXPECT().Update(uint(1), uint(10), input).Return(nil)

	services := &service.Service{Book: mockBook}
	h := Handler{services}

	r := gin.New()
	r.PUT("/books/:id", func(c *gin.Context) {
		setUser(c, 1)
		h.updateBook(c)
	})

	body := `{"title":"New Title","author":"New Author"}`

	req := httptest.NewRequest("PUT", "/books/10", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"book updated"}`, w.Body.String())
}

func TestHandler_updateBook_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := &service.Service{Book: mock_service.NewMockBook(ctrl)}
	h := Handler{services}

	r := gin.New()
	r.PUT("/books/:id", h.updateBook)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/books/abc", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"invalid book id"}`, w.Body.String())
}

func TestHandler_deleteBook_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBook := mock_service.NewMockBook(ctrl)
	mockBook.EXPECT().Delete(uint(1), uint(7)).Return(nil)

	services := &service.Service{Book: mockBook}
	h := Handler{services}

	r := gin.New()
	r.DELETE("/books/:id", func(c *gin.Context) {
		setUser(c, 1)
		h.deleteBook(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/books/7", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"book deleted"}`, w.Body.String())
}

func TestHandler_deleteBook_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := &service.Service{Book: mock_service.NewMockBook(ctrl)}
	h := Handler{services}

	r := gin.New()
	r.DELETE("/books/:id", h.deleteBook)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/books/abc", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"invalid book id"}`, w.Body.String())
}

func TestUserIDFromContext_OK(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(userCtx, uint(42))

	id, err := UserIDFromContext(c)

	assert.Equal(t, nil, err)
	assert.Equal(t, uint(42), id)
}

func TestUserIDFromContext_NotFound(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, err := UserIDFromContext(c)

	assert.Equal(t, "user_id not found in context", err.Error())
}

func TestUserIDFromContext_InvalidType(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(userCtx, "bad type")

	_, err := UserIDFromContext(c)

	assert.Equal(t, "user id is of invalid type", err.Error())
}
