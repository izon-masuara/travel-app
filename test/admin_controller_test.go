package test

// if you want to test more than 30 s you can change in the settings in vscode go test time out

import (
	"context"
	"encoding/json"
	"io"
	"kautsar/travel-app-api/app"
	"kautsar/travel-app-api/exception"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/middleware"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDb() *mongo.Database {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	helper.PanicIfError(err)
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	helper.PanicIfError(err)

	return client.Database("travel_app_test")
}

func startDb(db *mongo.Database) {
	err := db.Collection("account").Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	pass, _ := helper.HashPassword("admin123")
	_, err = db.Collection("account").InsertOne(context.Background(), bson.M{
		"name":     "admin",
		"username": "admin",
		"password": pass,
		"role":     "admin",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func setupRouter(db *mongo.Database) http.Handler {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	validate := validator.New()
	router := httprouter.New()

	app.RegisterAdminRoutes(router, db, validate)
	app.RegisterUserRoutes(router, db, validate)

	router.PanicHandler = exception.ErrorHandler

	newRouter := middleware.NewAuthMiddleware(router)
	return newRouter
}

func TestAdminLoginSuccess(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	assert.Equal(t, 200, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Contains(t, "OK", responseBody["status"])
}

func TestAdminLoginBadRequestValidationUsernameEmpty(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`
	{
		"username":"",
		"password":"admin123"
	}`)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	assert.Equal(t, 400, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Contains(t, "BAD REQUEST", responseBody["status"])
}

func TestAdminLoginBadRequestValidationPasswordEmpty(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`
	{
		"username":"admin",
		"password":""
	}`)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	assert.Equal(t, 400, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Contains(t, "BAD REQUEST", responseBody["status"])
}

func TestAdminLoginBadRequestValidationUsernameLength(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`
	{
		"username":"ad",
		"password":"admin123"
	}`)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	assert.Equal(t, 400, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Contains(t, "BAD REQUEST", responseBody["status"])
}

func TestAdminLoginBadRequestValidationPasswordLength(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`
	{
		"username":"admin",
		"password":"admi3"
	}`)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	assert.Equal(t, 400, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Contains(t, "BAD REQUEST", responseBody["status"])
}

func TestAdminLoginFailedUnautorized(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`
	{
		"username":"admin3",
		"password":"admin123"
	}`)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	assert.Equal(t, 401, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "Unauthorized", responseBody["status"])
	assert.Equal(t, "Username or password is wrong", responseBody["data"])
}

func TestAdminCreateNewAccountSuccess(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"donggala",
		"password":"donggala123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 200, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 200, int(responseCreate.StatusCode))
	assert.Equal(t, "OK", responseBodyCreate["status"])
	assert.Equal(t, "Success create new account", responseBodyCreate["data"])
}

func TestAdminCreateNewAccountBadRequestEmptyName(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"",
		"username":"donggala",
		"password":"dongga123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 400, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 400, int(responseCreate.StatusCode))
	assert.Equal(t, "BAD REQUEST", responseBodyCreate["status"])
}

func TestAdminCreateNewAccountBadRequestFieldEmptyName(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"username":"donggala",
		"password":"dongga123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 400, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 400, int(responseCreate.StatusCode))
	assert.Equal(t, "BAD REQUEST", responseBodyCreate["status"])
}

func TestAdminCreateNewAccountBadRequestEmptyUsername(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"",
		"password":"dongga123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 400, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 400, int(responseCreate.StatusCode))
	assert.Equal(t, "BAD REQUEST", responseBodyCreate["status"])
}

func TestAdminCreateNewAccountBadRequestEmptyPassword(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"donggala",
		"password":""
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 400, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 400, int(responseCreate.StatusCode))
	assert.Equal(t, "BAD REQUEST", responseBodyCreate["status"])
}

func TestAdminCreateNewAccountBadRequestUsernameLength(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"don",
		"password":"donggala123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 400, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 400, int(responseCreate.StatusCode))
	assert.Equal(t, "BAD REQUEST", responseBodyCreate["status"])
}

func TestAdminCreateNewAccountBadRequestPasswordLength(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"don",
		"password":"don23"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 400, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 400, int(responseCreate.StatusCode))
	assert.Equal(t, "BAD REQUEST", responseBodyCreate["status"])
}

func TestAdminCreateNewAccountUnauthorized(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"donggala",
		"password":"donggala123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", "jskkjsadlsa")
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 401, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		panic(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 401, int(responseCreate.StatusCode))
	assert.Equal(t, "UNAUTHORIZED", responseBodyCreate["status"])
}

func TestAdminGetListAccountSuccess(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"donggala",
		"password":"donggala123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 200, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 200, int(responseCreate.StatusCode))
	assert.Equal(t, "OK", responseBodyCreate["status"])
	assert.Equal(t, "Success create new account", responseBodyCreate["data"])

	//List of account
	requestAccounts := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/account", nil)
	requestAccounts.Header.Add("TOKEN", token)
	recorderAccounts := httptest.NewRecorder()
	router.ServeHTTP(recorderAccounts, requestAccounts)
	responseAccounts := recorderAccounts.Result()
	assert.Equal(t, 200, responseAccounts.StatusCode)
	bodyAccounts, err := io.ReadAll(responseAccounts.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyAccount map[string]interface{}
	err = json.Unmarshal(bodyAccounts, &responseBodyAccount)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 200, int(responseAccounts.StatusCode))
	assert.Equal(t, "OK", responseBodyAccount["status"])
	accountsInterface := responseBodyAccount["data"].([]interface{})

	donggala := accountsInterface[0].(map[string]interface{})["name"]
	assert.Equal(t, "donggala", donggala)
}

func TestAdminGetListAccountUnauthorized(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"donggala",
		"password":"donggala123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 200, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 200, int(responseCreate.StatusCode))
	assert.Equal(t, "OK", responseBodyCreate["status"])
	assert.Equal(t, "Success create new account", responseBodyCreate["data"])

	//List of account
	requestAccounts := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/account", nil)
	requestAccounts.Header.Add("TOKEN", "kasjdkiwqpo")
	recorderAccounts := httptest.NewRecorder()
	router.ServeHTTP(recorderAccounts, requestAccounts)
	responseAccounts := recorderAccounts.Result()
	assert.Equal(t, 401, responseAccounts.StatusCode)
	bodyAccounts, err := io.ReadAll(responseAccounts.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyAccount map[string]interface{}
	json.Unmarshal(bodyAccounts, &responseBodyAccount)
	assert.Equal(t, 401, int(responseAccounts.StatusCode))
	assert.Equal(t, "UNAUTHORIZED", responseBodyAccount["status"])
}

func TestAdminResetPasswordSuccess(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"donggala",
		"password":"donggala123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 200, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 200, int(responseCreate.StatusCode))
	assert.Equal(t, "OK", responseBodyCreate["status"])
	assert.Equal(t, "Success create new account", responseBodyCreate["data"])

	requestAccounts := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/account", nil)
	requestAccounts.Header.Add("TOKEN", token)
	recorderAccounts := httptest.NewRecorder()
	router.ServeHTTP(recorderAccounts, requestAccounts)
	responseAccounts := recorderAccounts.Result()
	assert.Equal(t, 200, responseAccounts.StatusCode)
	bodyAccounts, err := io.ReadAll(responseAccounts.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyAccount map[string]interface{}
	json.Unmarshal(bodyAccounts, &responseBodyAccount)
	assert.Equal(t, 200, int(responseAccounts.StatusCode))
	assert.Equal(t, "OK", responseBodyAccount["status"])
	accountsInterface := responseBodyAccount["data"].([]interface{})

	donggala := accountsInterface[0].(map[string]interface{})["name"]

	assert.Equal(t, "donggala", donggala)

	// Admin reset password
	donggalaId := accountsInterface[0].(map[string]interface{})["_id"]

	requestResetPassword := httptest.NewRequest(http.MethodPatch, "http://localhost:3000/api/v1/account/"+donggalaId.(string), nil)
	requestResetPassword.Header.Add("Content-Type", "application/json")
	requestResetPassword.Header.Add("TOKEN", token)
	recorderResetPassword := httptest.NewRecorder()
	router.ServeHTTP(recorderResetPassword, requestResetPassword)
	responseResetPassword := recorderResetPassword.Result()
	assert.Equal(t, 200, responseResetPassword.StatusCode)
	bodyResetPassword, err := io.ReadAll(responseResetPassword.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyResetPassword map[string]interface{}
	json.Unmarshal(bodyResetPassword, &responseBodyResetPassword)
	assert.Equal(t, 200, int(responseResetPassword.StatusCode))
	assert.Equal(t, "OK", responseBodyResetPassword["status"])
	assert.Equal(t, "Success reset password", responseBodyResetPassword["data"])

	//Operator login with default password
	requestOperatorNewPassword := strings.NewReader(`
	{
		"username":"donggala",
		"password":"defaultpassword"
	}`)
	requestLoginOperator := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestOperatorNewPassword)
	requestLoginOperator.Header.Add("Content-Type", "application/json")
	recorderLoginOperator := httptest.NewRecorder()
	router.ServeHTTP(recorderLoginOperator, requestLoginOperator)
	responseLoginOperator := recorderLoginOperator.Result()
	assert.Equal(t, 200, responseLoginOperator.StatusCode)
	bodyLoginOperator, err := io.ReadAll(responseLoginOperator.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLoginOperator map[string]interface{}
	json.Unmarshal(bodyLoginOperator, &responseBodyLoginOperator)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])

}

func TestAdminResetPasswordNotFoundId(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"donggala",
		"password":"donggala123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 200, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 200, int(responseCreate.StatusCode))
	assert.Equal(t, "OK", responseBodyCreate["status"])
	assert.Equal(t, "Success create new account", responseBodyCreate["data"])

	// Admin reset password
	fakeId := primitive.NewObjectID().Hex()
	requestResetPassword := httptest.NewRequest(http.MethodPatch, "http://localhost:3000/api/v1/account/"+fakeId, nil)
	requestResetPassword.Header.Add("Content-Type", "application/json")
	requestResetPassword.Header.Add("TOKEN", token)
	recorderResetPassword := httptest.NewRecorder()
	router.ServeHTTP(recorderResetPassword, requestResetPassword)
	responseResetPassword := recorderResetPassword.Result()
	assert.Equal(t, 404, responseResetPassword.StatusCode)
	bodyResetPassword, err := io.ReadAll(responseResetPassword.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyResetPassword map[string]interface{}
	json.Unmarshal(bodyResetPassword, &responseBodyResetPassword)
	assert.Equal(t, 404, int(responseResetPassword.StatusCode))
	assert.Equal(t, "NOT FOUND", responseBodyResetPassword["status"])
	assert.Equal(t, "Data not found", responseBodyResetPassword["data"])
}

func TestAdminResetPasswordUnauthorized(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin reset password
	fakeId := primitive.NewObjectID().Hex()
	requestResetPassword := httptest.NewRequest(http.MethodPatch, "http://localhost:3000/api/v1/account/"+fakeId, nil)
	requestResetPassword.Header.Add("Content-Type", "application/json")
	requestResetPassword.Header.Add("TOKEN", "ksdjiiwmmiqmsd")
	recorderResetPassword := httptest.NewRecorder()
	router.ServeHTTP(recorderResetPassword, requestResetPassword)
	responseResetPassword := recorderResetPassword.Result()
	assert.Equal(t, 401, responseResetPassword.StatusCode)
	bodyResetPassword, err := io.ReadAll(responseResetPassword.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyResetPassword map[string]interface{}
	json.Unmarshal(bodyResetPassword, &responseBodyResetPassword)
	assert.Equal(t, 401, int(responseResetPassword.StatusCode))
	assert.Equal(t, "UNAUTHORIZED", responseBodyResetPassword["status"])
}

func TestAdminDeleteSuccess(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin create new operator account
	requestBodyCreateNewAccount := strings.NewReader(`
	{
		"name":"donggala",
		"username":"donggala",
		"password":"donggala123"
	}`)
	requestCreateNewAccount := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/account", requestBodyCreateNewAccount)
	requestCreateNewAccount.Header.Add("Content-Type", "application/json")
	requestCreateNewAccount.Header.Add("TOKEN", token)
	recorderCreate := httptest.NewRecorder()
	router.ServeHTTP(recorderCreate, requestCreateNewAccount)
	responseCreate := recorderCreate.Result()
	assert.Equal(t, 200, responseCreate.StatusCode)
	bodyCreate, err := io.ReadAll(responseCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyCreate map[string]interface{}
	json.Unmarshal(bodyCreate, &responseBodyCreate)
	assert.Equal(t, 200, int(responseCreate.StatusCode))
	assert.Equal(t, "OK", responseBodyCreate["status"])
	assert.Equal(t, "Success create new account", responseBodyCreate["data"])

	//List of account
	requestAccounts := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/account", nil)
	requestAccounts.Header.Add("TOKEN", token)
	recorderAccounts := httptest.NewRecorder()
	router.ServeHTTP(recorderAccounts, requestAccounts)
	responseAccounts := recorderAccounts.Result()
	assert.Equal(t, 200, responseAccounts.StatusCode)
	bodyAccounts, err := io.ReadAll(responseAccounts.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyAccount map[string]interface{}
	json.Unmarshal(bodyAccounts, &responseBodyAccount)
	assert.Equal(t, 200, int(responseAccounts.StatusCode))
	assert.Equal(t, "OK", responseBodyAccount["status"])
	accountsInterface := responseBodyAccount["data"].([]interface{})

	donggala := accountsInterface[0].(map[string]interface{})["name"]

	assert.Equal(t, "donggala", donggala)

	// Admin reset password
	donggalaId := accountsInterface[0].(map[string]interface{})["_id"]

	requestDeleteOperator := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/v1/account/"+donggalaId.(string), nil)
	requestDeleteOperator.Header.Add("Content-Type", "application/json")
	requestDeleteOperator.Header.Add("TOKEN", token)
	recorderDeleteOperator := httptest.NewRecorder()
	router.ServeHTTP(recorderDeleteOperator, requestDeleteOperator)
	responseDeleteOperator := recorderDeleteOperator.Result()
	assert.Equal(t, 200, responseDeleteOperator.StatusCode)
	bodyDeleteOperator, err := io.ReadAll(responseDeleteOperator.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyDeleteOperator map[string]interface{}
	json.Unmarshal(bodyDeleteOperator, &responseBodyDeleteOperator)
	assert.Equal(t, 200, int(responseDeleteOperator.StatusCode))
	assert.Equal(t, "OK", responseBodyDeleteOperator["status"])
	assert.Equal(t, "Success delete account", responseBodyDeleteOperator["data"])

	requestBodyOperator := strings.NewReader(`
	{
		"username":"donggala",
		"password":"donggala123"
	}`)
	requestLoginFailed := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyOperator)
	requestLoginFailed.Header.Add("Content-Type", "application/json")
	recorderLoginFailed := httptest.NewRecorder()
	router.ServeHTTP(recorderLoginFailed, requestLoginFailed)
	responseLoginFailed := recorderLoginFailed.Result()
	assert.Equal(t, 401, responseLoginFailed.StatusCode)
	bodyLoginFailed, err := io.ReadAll(responseLoginFailed.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLoginFailed map[string]interface{}
	json.Unmarshal(bodyLoginFailed, &responseBodyLoginFailed)
	assert.Equal(t, 401, int(responseBodyLoginFailed["code"].(float64)))
	assert.Equal(t, "Unauthorized", responseBodyLoginFailed["status"])
	assert.Equal(t, "Username or password is wrong", responseBodyLoginFailed["data"])
}

func TestStuck(t *testing.T) {
	time.Sleep(80 * time.Millisecond)
}

func TestAdminDeleteNotFoundId(t *testing.T) {
	db := setupTestDb()
	startDb(db)
	router := setupRouter(db)

	// Admin Login to get token
	requestBodyAdmin := strings.NewReader(`
	{
		"username":"admin",
		"password":"admin123"
	}`)
	requestLogin := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBodyAdmin)
	requestLogin.Header.Add("Content-Type", "application/json")
	recorderLogin := httptest.NewRecorder()
	router.ServeHTTP(recorderLogin, requestLogin)
	responseLogin := recorderLogin.Result()
	assert.Equal(t, 200, responseLogin.StatusCode)
	bodyLogin, err := io.ReadAll(responseLogin.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyLogin map[string]interface{}
	json.Unmarshal(bodyLogin, &responseBodyLogin)
	assert.Equal(t, 200, int(responseBodyLogin["code"].(float64)))
	assert.Equal(t, "OK", responseBodyLogin["status"])
	token := responseBodyLogin["data"].(string)

	// Admin delete account with wrong id
	fakeId := primitive.NewObjectID().Hex()
	requestResetPassword := httptest.NewRequest(http.MethodPatch, "http://localhost:3000/api/v1/account/"+fakeId, nil)
	requestResetPassword.Header.Add("Content-Type", "application/json")
	requestResetPassword.Header.Add("TOKEN", token)
	recorderResetPassword := httptest.NewRecorder()
	router.ServeHTTP(recorderResetPassword, requestResetPassword)
	responseResetPassword := recorderResetPassword.Result()
	assert.Equal(t, 404, responseResetPassword.StatusCode)
	bodyResetPassword, err := io.ReadAll(responseResetPassword.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseBodyResetPassword map[string]interface{}
	json.Unmarshal(bodyResetPassword, &responseBodyResetPassword)
	assert.Equal(t, 404, int(responseResetPassword.StatusCode))
	assert.Equal(t, "NOT FOUND", responseBodyResetPassword["status"])
	assert.Equal(t, "Data not found", responseBodyResetPassword["data"])
}
