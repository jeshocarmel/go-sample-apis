package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

var router *httprouter.Router

//TestMain sets up the testing environment
func TestMain(m *testing.M) {
	router = setupRoutes()
	code := m.Run()
	os.Exit(code)
}

func TestIndex(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)

	expected := "{\"message\":\"Hello World\"}"
	actual := responseRecorder.Body.String()
	assert.Equal(t, expected, actual)

	assert.Equal(t, 200, responseRecorder.Code)
}

func TestSecuredItem1(t *testing.T) {

	req, _ := http.NewRequest("GET", "/secured/item1", nil)
	req.SetBasicAuth("abc", "123")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)

	expected := "{\"message1\":\"Basic Authentication Success\",\"message2\":\"This is item 1\"}"
	actual := responseRecorder.Body.String()
	assert.Equal(t, expected, actual)
	assert.Equal(t, 200, responseRecorder.Code)
}

func TestSecuredItem2(t *testing.T) {

	req, _ := http.NewRequest("GET", "/secured/item2", nil)
	token, _ := generateJWT()
	req.Header.Add("Token", token)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)

	expected := "{\"message1\":\"JWT Authentication Success\",\"message2\":\"This is item 2\"}"
	actual := responseRecorder.Body.String()
	assert.Equal(t, expected, actual)
	assert.Equal(t, 200, responseRecorder.Code)
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Jesho Carmel"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(secretkey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
