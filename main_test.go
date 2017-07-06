package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

  "github.com/contd/links"
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize("apiuser", "wj5np47dn", "links_test")
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/links", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentLink(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/link/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Link not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Link not found'. Got '%s'", m["error"])
	}
}

func TestCreateLink(t *testing.T) {
	clearTable()

	payload := []byte(`{"url":"new url","category":"general","created_on":"2017-02-02T12:00:00TZ","done":0}`)

	req, _ := http.NewRequest("POST", "/link", bytes.NewBuffer(payload))
	response := executeRequest(req)

	//checkResponseCode(t, http.StatusCreated, response.Code)
	if response.Code != http.StatusCreated {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusCreated, response.Code)
	}

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["url"] != "new url" {
		t.Errorf("Expected link url to be 'new url'. Got '%v'", m["url"])
	}

	if m["category"] != "general" {
		t.Errorf("Expected link category to be 'general'. Got '%v'", m["category"])
	}

	if m["created_on"] != "2017-02-02T12:00:00TZ" {
		t.Errorf("Expected link created_on to be '2017-02-02T12:00:00TZ'. Got '%v'", m["created_on"])
	}
}

func TestGetLink(t *testing.T) {
	clearTable()
	addLink()

	req, _ := http.NewRequest("GET", "/link/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateLink(t *testing.T) {
	clearTable()
	addLink()

	req, _ := http.NewRequest("GET", "/link/1", nil)
	response := executeRequest(req)

	var originalLink map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalLink)

	payload := []byte(`{"url":"orig url 1 - updated name","category":"specific"}`)

	req, _ = http.NewRequest("PUT", "/link/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalLink["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalLink["id"], m["id"])
	}

	if m["url"] == originalLink["url"] {
		t.Errorf("Expected the url to change from '%v' to '%v'. Got '%v'", originalLink["url"], m["url"], m["url"])
	}

	if m["category"] == originalLink["price"] {
		t.Errorf("Expected the category to change from '%v' to '%v'. Got '%v'", originalLink["category"], m["category"], m["category"])
	}
}

func TestDeleteLink(t *testing.T) {
	clearTable()
	addLink()

	req, _ := http.NewRequest("GET", "/link/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/link/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/link/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM links")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addLink() {
	_, err := a.DB.Exec("INSERT INTO links (id,url,category,created_on,done) VALUES (1,?,?,?,?)", "new url ", "general", "2017-02-02T12:00:00TZ", 0)
	if err != nil {
		log.Fatal(err)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS links
(
	id INT NOT NULL AUTO_INCREMENT,
	url VARCHAR(255) NOT NULL,
	category VARCHAR(50) NOT NULL,
	created_on CHAR(25) NOT NULL,
	done INT,
	PRIMARY KEY (id)
)`
