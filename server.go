package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

var dictionary map[string][]string

func StringToRuneSlice(s string) []rune {
	var r []rune
	for _, runeValue := range s {
		r = append(r, runeValue)
	}
	return r
}

func SortStringByCharacter(s string) string {
	r := StringToRuneSlice(s)
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return string(r)
}

func PrepareString(s string) string {
	return SortStringByCharacter(strings.ToLower(s))
}


func loadWords(res http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	var words []string
	err = json.Unmarshal(b, &words)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	for _, word := range words {
		preparedString := PrepareString(word)
		dictionary[preparedString] = append(dictionary[preparedString], word)
	}
}

func getAnagrams(res http.ResponseWriter, req *http.Request) {
	word := req.URL.Query().Get("word")
	if word == "" {
		http.Error(res, "Error: 'word' parameter was not provided", http.StatusBadRequest)
		return
	}

	preparedString := PrepareString(word)
	output, err := json.Marshal(dictionary[preparedString])
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	res.Write(output)
}

func main() {
	dictionary = make(map[string][]string)

	http.HandleFunc("/load", loadWords)
	http.HandleFunc("/get", getAnagrams)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

