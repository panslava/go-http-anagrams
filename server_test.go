package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmptyGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	GetAnagrams(w, req)
	resp := w.Result()

	responseStatusCode := resp.StatusCode
	expectedStatusCode := http.StatusBadRequest
	if responseStatusCode != expectedStatusCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseStatusCode, expectedStatusCode)
	}
}

func TestNotAdded(t *testing.T) {
	req, err := http.NewRequest("GET", "/get?word=haha", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	GetAnagrams(w, req)

	responseString := w.Body.String()
	expectedString := "null"
	if responseString != expectedString {
		t.Errorf("handler returned wrong response: got %v want %v",
			responseString, expectedString)
	}
}

func TestBasic(t *testing.T) {
	dictionary = make(map[string][]string)

	inputWords := []string{"foobar", "aabb", "baba", "boofar", "test"}
	body, _ := json.Marshal(inputWords)
	req1, err1 := http.NewRequest("POST", "/load", bytes.NewReader(body))
	if err1 != nil {
		t.Fatal(err1)
	}
	w1 := httptest.NewRecorder()
	LoadWords(w1, req1)

	req, err := http.NewRequest("GET", "/get?word=raboof", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	GetAnagrams(w, req)

	b, err := ioutil.ReadAll(w.Body)
	var responseWords []string
	err = json.Unmarshal(b, &responseWords)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := []string{"foobar", "boofar"}
	equals := true
	if len(responseWords) != len(expectedResult) {
		equals = false
	}
	for i := range expectedResult {
		if expectedResult[i] != responseWords[i] {
			equals = false
		}
	}

	if !equals {
		t.Errorf("handler returned array %v, but expected %v", strings.Join(responseWords, ","), strings.Join(expectedResult, ","))
	}
}

func TestCaseInsensitive(t *testing.T) {
	dictionary = make(map[string][]string)

	inputWords := []string{"FooBar", "aabb", "bAba", "bOOfar", "test"}
	body, _ := json.Marshal(inputWords)
	req1, err1 := http.NewRequest("POST", "/load", bytes.NewReader(body))
	if err1 != nil {
		t.Fatal(err1)
	}
	w1 := httptest.NewRecorder()
	LoadWords(w1, req1)

	req, err := http.NewRequest("GET", "/get?word=BOOFRA", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	GetAnagrams(w, req)

	b, err := ioutil.ReadAll(w.Body)
	var responseWords []string
	err = json.Unmarshal(b, &responseWords)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := []string{"FooBar", "bOOfar"}

	equals := true
	if len(responseWords) != len(expectedResult) {
		equals = false
	}
	for i := range expectedResult {
		if expectedResult[i] != responseWords[i] {
			equals = false
		}
	}

	if !equals {
		t.Errorf("handler returned array %v, but expected %v", strings.Join(responseWords, ","), strings.Join(expectedResult, ","))
	}
}
