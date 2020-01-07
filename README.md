# Sample REST APIs in golang with authentication

The application contains 3 resources that can be accessed using HTTP REST (GET) protocol.
```go
router := httprouter.New()
router.GET("/", index)
router.GET("/secured/item1", basicAuth(getItem1))
router.GET("/secured/item2", jwtAuth(getItem2))
```
 1. **/** is unsecured.
 2. **/secured/item1** is secured using **HTTP basic access authentication** i.e. with a username and password
 3. **/secured/item2** is secured using **JWT authentication**
 
To start server execute ```go run main.go```
To run tests execute ```go test -v```
