package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"br.com.mlabs/api"
	"br.com.mlabs/models"
	"br.com.mlabs/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var tx *gorm.DB

func TestMain(t *testing.M) {
	storage.ConnectTest()
	tx = storage.StartTest()

	t.Run()

	tx.Exec("TRUNCATE parkings CASCADE;")
	tx.Exec("ALTER SEQUENCE parkings_id_seq RESTART WITH 1")
	tx.Exec("ALTER SEQUENCE payments_id_seq RESTART WITH 1")
}

func TestReservationHappyPath(t *testing.T) {
	request := models.ParkingRequest{
		Plate: "ABC-1234",
	}
	jsonBytes, _ := json.Marshal(request)

	req, _ := http.NewRequest(http.MethodPost, "/parking", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusOK)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, strings.Contains(string(bts), "{\"id\":"), true)
}

func TestReservationValidationError(t *testing.T) {
	request := models.ParkingRequest{
		Plate: "ABC-123",
	}
	jsonBytes, _ := json.Marshal(request)

	req, _ := http.NewRequest(http.MethodPost, "/parking", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusBadRequest)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"Plate must be valid, format: AAA-1234\"}")
}

func TestHistoryHappyPath(t *testing.T) {
	plate := "ABC-1234"
	url := fmt.Sprintf("/parking/%s", plate)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusOK)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, strings.Contains(string(bts), "{\"id\":"), true)
}

func TestHistoryNotFound(t *testing.T) {
	plate := "ZZZ-1234"
	url := fmt.Sprintf("/parking/%s", plate)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusNotFound)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"Not found\"}")
}

func TestHistoryNotValid(t *testing.T) {
	plate := "ZZZ-123"
	url := fmt.Sprintf("/parking/%s", plate)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusBadRequest)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"Plate must be valid, format: AAA-1234\"}")
}

// Tests checkin

func TestPayHappyPath(t *testing.T) {
	id := 3
	url := fmt.Sprintf("/parking/%d/pay", id)

	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusOK)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"Paid\"}")
}

func TestPayAlreadyPaid(t *testing.T) {
	id := 1
	url := fmt.Sprintf("/parking/%d/pay", id)

	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusOK)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"You have already paid\"}")
}

func TestPayNotFound(t *testing.T) {
	id := 9999
	url := fmt.Sprintf("/parking/%d/pay", id)

	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusNotFound)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"Not found\"}")
}

func TestPayBadRequest(t *testing.T) {
	id := "as"
	url := fmt.Sprintf("/parking/%s/pay", id)

	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusBadRequest)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"ID must be valid\"}")
}

// Tests Checkout

func TestCheckoutHappyPath(t *testing.T) {
	id := 1
	url := fmt.Sprintf("/parking/%d/out", id)

	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusOK)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"Checked out\"}")
}

func TestCheckoutAlreadyDone(t *testing.T) {
	id := 2
	url := fmt.Sprintf("/parking/%d/out", id)

	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusOK)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"You have already checked out\"}")
}

func TestCheckoutPayPending(t *testing.T) {
	id := 5
	url := fmt.Sprintf("/parking/%d/out", id)

	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusPaymentRequired)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"You have to pay first\"}")
}

func TestCheckoutNotFound(t *testing.T) {
	id := 9999
	url := fmt.Sprintf("/parking/%d/out", id)

	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusNotFound)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"Not found\"}")
}

func TestCheckoutBadRequest(t *testing.T) {
	id := "as"
	url := fmt.Sprintf("/parking/%s/out", id)

	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, api.NewParkingRouter)

	assert.Equal(t, response.Code, http.StatusBadRequest)
	bts, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(bts), "{\"response\":\"ID must be valid\"}")
}

func executeRequest(req *http.Request, subRouter func(router *mux.Router)) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()

	router := mux.NewRouter()
	subRouter(router)

	router.ServeHTTP(responseRecorder, req)

	return responseRecorder
}
