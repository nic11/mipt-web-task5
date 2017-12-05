package main

import "log"
import "net/http"
import "encoding/json"
import "io/ioutil"
import "strconv"
import "github.com/gorilla/mux"

var urls map[string]string

func nextKey() string {
	return strconv.Itoa(len(urls))
}

func NewLink(w http.ResponseWriter, r *http.Request) {
	var body []byte
	defer r.Body.Close()
	body, _ = ioutil.ReadAll(r.Body)
	
	var reqJson map[string]string
	if err := json.Unmarshal(body, &reqJson); err != nil {
		panic(err)
	}
	url := reqJson["url"]
	key := nextKey()
	// println(url)
	// println(key)
	urls[key] = url

	resJson := map[string]string{"key": key}
	res, _ := json.Marshal(resJson)
	w.Write(res)
}

func GetRedirect(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	url, ok := urls[key]
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, 301)
}

func main() {
	urls = make(map[string]string)

	r := mux.NewRouter()
	r.HandleFunc("/", NewLink).Methods("POST")
	r.HandleFunc("/{key}", GetRedirect).Methods("GET")

	log.Fatal(http.ListenAndServe(":8082", r))
}