package main

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

//secretkey represents the secret key used for JWT token generation
var secretkey = []byte("this_is_a_secret_key_for_007")

//basicAuth is a interceptor before accessing a secured URL
func basicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, pass, ok := r.BasicAuth()

		if ok && checkUsernameAndPassword(user, pass) {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Deny access for the the handle
			http.Error(w, "Not Authorized", 401)
			return
		}
	}
}

func checkUsernameAndPassword(username, password string) bool {
	return username == "abc" && password == "123"
}

//credit - https://tutorialedge.net/golang/authenticating-golang-rest-api-with-jwts/

//jwtAuth is a interceptor before accessing a secured URL
func jwtAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return secretkey, nil
			})

			if err != nil {
				http.Error(w, err.Error(), 500)
			}

			if token.Valid {
				h(w, r, ps)
			}

		} else {
			// Deny access for the the handle
			http.Error(w, "Not Authorized", 401)
			return
		}
	}
}

//index is unsecured
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result := make(map[string]string)
	result["message"] = "Hello World"

	bytes, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

//getItem1 is secured with 'http-basicAuth'
func getItem1(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	result := make(map[string]string)
	result["message1"] = "Basic Authentication Success"
	result["message2"] = "This is item 1"

	bytes, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

//getItem2 is secured with JWT authentication
func getItem2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	result := make(map[string]string)
	result["message1"] = "JWT Authentication Success"
	result["message2"] = "This is item 2"

	bytes, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func setupRoutes() *httprouter.Router {

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/secured/item1", basicAuth(getItem1))
	router.GET("/secured/item2", jwtAuth(getItem2))

	return router

}

func main() {

	router := setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", router))
}
