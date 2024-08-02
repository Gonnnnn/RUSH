package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"rush/auth"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockIdentifierFetcher struct {
	valueToReturn auth.UserIdentifier
	errToReturn   error
}

func (m *mockIdentifierFetcher) GetUserIdentifier(token string) (auth.UserIdentifier, error) {
	return m.valueToReturn, m.errToReturn
}

func TestGinAuthMiddleware(t *testing.T) {
	t.Run("Fails if the identifierFetcher fails to get the user identifier", func(t *testing.T) {
		mockIdentifierFetcher := &mockIdentifierFetcher{valueToReturn: auth.UserIdentifier{}, errToReturn: errors.New("mock error")}

		middleware := ginAuthMiddleware(mockIdentifierFetcher)
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)
		ginContext.Request = httptest.NewRequest("GET", "/", nil)
		ginContext.Request.Header.Set("Authorization", "token")

		middleware(ginContext)

		assert.Equal(t, http.StatusUnauthorized, ginContext.Writer.Status())
	})

	t.Run("Succeeds if the identifierFetcher gets the user identifier", func(t *testing.T) {
		mockIdentifierFetcher := &mockIdentifierFetcher{valueToReturn: auth.UserIdentifier{}, errToReturn: nil}

		middleware := ginAuthMiddleware(mockIdentifierFetcher)
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)
		ginContext.Request = httptest.NewRequest("GET", "/", nil)
		ginContext.Request.Header.Set("Authorization", "token")

		middleware(ginContext)

		assert.Equal(t, http.StatusOK, ginContext.Writer.Status())
	})
}
