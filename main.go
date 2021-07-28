package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8000", "http service address")
var ROOMS = make(map[string]*Hub)

// getHub return a hub by it's id
// and bool value indicating if the hub is newly created or not
func getHub(id string) (h *Hub, isNew bool) {
	isNew = true

	h, ok := ROOMS[id]

	if ok {
		isNew = false
	} else {
		h = newHub()
		ROOMS[id] = h
	}

	return
}

// ignore this please
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var message map[string]interface{}
	json.Unmarshal([]byte(`{"hello": "world"}`), &message)
	fmt.Println(message)
	json.NewEncoder(w).Encode(message)
}

func RoomHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h, isNew := getHub(id)
	if isNew {
		go h.run()
	}

	serveWs(h, w, r, id)
}

func CreateHandle(w http.ResponseWriter, r *http.Request) {

}

func main() {
	flag.Parse()
	r := mux.NewRouter()

	// public hub for any connection in /ws
	hub := newHub()
	go hub.run()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r, "0")
	})

	r.HandleFunc("/", HelloWorld)
	r.HandleFunc("/room/{id}", RoomHandle)

	log.Fatal(http.ListenAndServe(*addr, r))
}
