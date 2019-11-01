package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"

	"github.com/imw-challenge/back/db"
	"github.com/imw-challenge/back/types"
)

var a *API
var testMessages []*types.Message
var err error

func initEmptyDB() *db.MessageDB {
	mdb, err := db.InitMessageDB()
	if err != nil {
		panic(err)
	}
	return mdb
}

func setTestMessages() {
	var newTestMessages []*types.Message
	messageStrings := []string{
		`{"id":"A5D00000-7DE9-7E69-C311-763310C9AA54","name":"Isaac Wilder","email":"name@fake.domain","text":"hi there","time":"2015-11-05T13:15:17-07:00"}`,
		`{"id":"B4F7A417-424E-2B99-87B6-5CA0744B7BBD","name":"Reggie Tester","email":"false@email.address","text":"lorem ipsum dolor sit amet","time":"2016-04-10T15:15:17-07:00"}`,
		`{"id":"2C7BCEC7-CD14-D6E5-3FBF-F9551375429A","name":"Alex Mustermann","email":"fake@site.biz","text":"testing","time":"2017-05-30T15:26:38-07:00"}`,
		`{"id":"B5D11111-7DE9-7E69-C311-763310C9AA54","name":"Isaac Wilder","email":"name@fake.domain","text":"hi again","time":"2018-09-10T00:15:00-07:00"}`,
		`{"id":"C5D22222-7DE9-7E69-C311-763310C9AA54","name":"Isaac Wilder","email":"name@fake.domain","text":"another one","time":"2019-04-10T10:10:10-07:00"}`}
	for _, m := range messageStrings {
		msg := new(types.Message)
		msg.UnmarshalJSON([]byte(m))
		newTestMessages = append(newTestMessages, msg)
	}
	testMessages = newTestMessages
}

func initPopulatedDB() *db.MessageDB {
	mdb := initEmptyDB()
	setTestMessages()
	mdb.InsertMessages(testMessages)
	return mdb
}

func setup() {
	mdb := initPopulatedDB()
	a, err = InitAPI(mdb)
	a.SetRoutes()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

/*
func TestPostMessage(t *testing.T) {
	req, _ := http.NewRequest("POST", "/public/message", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}
func TestPutMessage(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/private/message", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

*/
func TestPostMessage(t *testing.T) {
	//Check that request without body:
	// Return bad request
	req, _ := http.NewRequest("POST", "/public/message", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	//Check that request with malformed body:
	// Returns bad request
	bodyJsonString := `{""}`
	bodyJsonBytes := []byte(bodyJsonString)
	req, _ = http.NewRequest("POST", "/public/message", bytes.NewBuffer(bodyJsonBytes))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	//Check that request with well-formed body:
	// Returns success
	newMessageID := "Z9D00000-XXXX-7E69-C3PO-763310C9AA54"
	newMessageJsonString := `{"id":"` + newMessageID + `","name":"Martin Hipsh","email":"another@fake.email","text":"anyone can post!","time":"2019-11-01T14:09:16+02:00"}`
	bodyJsonBytes = []byte(newMessageJsonString)
	req, _ = http.NewRequest("POST", "/public/message", bytes.NewBuffer(bodyJsonBytes))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// Successfully inserts the message
	bodyJsonString = `{"id":"` + newMessageID + `"}`
	bodyJsonBytes = []byte(bodyJsonString)
	req, _ = http.NewRequest("GET", "/private/message", bytes.NewBuffer(bodyJsonBytes))
	req.SetBasicAuth("admin", "back-challenge")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	var responseMessage types.Message
	err = json.Unmarshal([]byte(body), &responseMessage)
	if err != nil {
		t.Errorf("Expected valid JSON. Got %s. Error: %s", body, err)
	}

	responseJson, err := responseMessage.MarshalJSON()
	if err != nil {
		t.Errorf("Malformed response message: %#v", responseJson)
	}
	if string(responseJson) != newMessageJsonString {
		t.Errorf("Expected message text %s. Got %s", newMessageJsonString, string(responseJson))
	}

}

func TestPutMessage(t *testing.T) {
	//Check that request without auth fails
	req, _ := http.NewRequest("PUT", "/private/message", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	//Check that request with wrong credentials fails
	req, _ = http.NewRequest("PUT", "/private/message", nil)
	req.SetBasicAuth("admin", "badPass")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	//Check that request with correct credentials and malformed body:
	// Returns bad request
	bodyJsonString := `{""}`
	bodyJsonBytes := []byte(bodyJsonString)
	req, _ = http.NewRequest("PUT", "/private/message", bytes.NewBuffer(bodyJsonBytes))
	req.SetBasicAuth("admin", "back-challenge")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	//Check that request with correct credentials and well-formed body:
	// Returns success
	newMessageText := "what once was old"
	bodyJsonString = `{"id":"` + testMessages[0].ID + `","text":"` + newMessageText + `"}`
	bodyJsonBytes = []byte(bodyJsonString)
	req, _ = http.NewRequest("PUT", "/private/message", bytes.NewBuffer(bodyJsonBytes))
	req.SetBasicAuth("admin", "back-challenge")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// Updates the message text
	bodyJsonString = `{"id":"` + testMessages[0].ID + `"}`
	bodyJsonBytes = []byte(bodyJsonString)
	req, _ = http.NewRequest("GET", "/private/message", bytes.NewBuffer(bodyJsonBytes))
	req.SetBasicAuth("admin", "back-challenge")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	var responseMessage types.Message
	err = json.Unmarshal([]byte(body), &responseMessage)
	if err != nil {
		t.Errorf("Expected valid JSON. Got %s. Error: %s", body, err)
	}

	if responseMessage.Text != newMessageText {
		t.Errorf("Expected message text %s. Got %s", newMessageText, responseMessage.Text)
	}

}

func TestGetMessage(t *testing.T) {
	//Check that request without auth fails
	req, _ := http.NewRequest("GET", "/private/message", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	//Check that request with wrong credentials fails
	req, _ = http.NewRequest("GET", "/private/message", nil)
	req.SetBasicAuth("admin", "badPass")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	//Check that request with correct credentials and malformed body:
	// Returns bad request
	bodyJsonString := `{""}`
	bodyJsonBytes := []byte(bodyJsonString)
	req, _ = http.NewRequest("GET", "/private/message", bytes.NewBuffer(bodyJsonBytes))
	req.SetBasicAuth("admin", "back-challenge")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	//Check that request with correct credentials and well-formed body:
	// Returns success
	bodyJsonString = `{"id":"` + testMessages[1].ID + `"}`
	bodyJsonBytes = []byte(bodyJsonString)
	req, _ = http.NewRequest("GET", "/private/message", bytes.NewBuffer(bodyJsonBytes))
	req.SetBasicAuth("admin", "back-challenge")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// Decodes correctly
	body := response.Body.String()
	var responseMessage types.Message
	err = json.Unmarshal([]byte(body), &responseMessage)
	if err != nil {
		t.Errorf("Expected valid JSON. Got %s. Error: %s", body, err)
	}

	// Returns the correct message
	if responseMessage.Name != testMessages[1].Name {
		t.Errorf("Expected %#v correct message to be returned. Got %#v", testMessages[1], responseMessage)
	}

}

func TestDump(t *testing.T) {
	//Reset DB
	setup()
	//Check that request without auth fails
	req, _ := http.NewRequest("GET", "/private/dump", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	//Check that request with wrong credentials fails
	req, _ = http.NewRequest("GET", "/private/dump", nil)
	req.SetBasicAuth("admin", "badPass")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	//Check that request with correct credentials:
	// Returns success
	req, _ = http.NewRequest("GET", "/private/dump", nil)
	req.SetBasicAuth("admin", "back-challenge")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()

	// Decodes correctly
	var responseMessages []types.Message
	err = json.Unmarshal([]byte(body), &responseMessages)
	if err != nil {
		t.Errorf("Expected valid JSON. Got %s. Error: %s", body, err)
	}

	// Has correct number of elements
	if len(responseMessages) != len(testMessages) {
		t.Errorf("Expected %d results. Got %d", len(testMessages), len(responseMessages))
	}

	// Is correctly sorted
	var times []int64
	for _, m := range responseMessages {
		times = append([]int64{m.Time}, times...)
	}

	if !sort.IsSorted(types.Int64Slice(times)) {
		t.Errorf("Expected results in reverse chronological order. Got %v", times)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.GetRouter().ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
