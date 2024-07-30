package endpoints

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Auth_WhenAuthorizationIsMissing_ReturnError(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})

	handlerFunc := Auth(nextHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	handlerFunc.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
	assert.Contains(t, res.Body.String(), "request does not contain an authorization header")
}

func Test_Auth_WhenAuthorizationIsInvalid_ReturnError(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})
	ValidateToken = func(token string, ctx context.Context) (string, error) {
		return "", errors.New("invalid token")
	}

	handlerFunc := Auth(nextHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer invalid token")
	res := httptest.NewRecorder()

	handlerFunc.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
	assert.Contains(t, res.Body.String(), "invalid token")
}

func Test_Auth_WhenAuthorizationIsInvalid_CallNextHandler(t *testing.T) {
	emailExpected := "teste@teste.com"
	var email string
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email = r.Context().Value("email").(string)
	})
	ValidateToken = func(token string, ctx context.Context) (string, error) {
		return emailExpected, nil
	}

	handlerFunc := Auth(nextHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer valid token")
	res := httptest.NewRecorder()

	handlerFunc.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, emailExpected, email)
}
