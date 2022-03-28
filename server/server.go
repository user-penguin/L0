package server

import (
	"L0/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const port = ":8080"

var data map[string]model.Order

func Serve(orders *map[string]model.Order) {
	data = *orders
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/order/", orderHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Добро пожаловать в мир Golang\n")
	log.Printf("Method:%s, Request URI:%s", r.Method, r.RequestURI)
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	uid := strings.TrimPrefix(r.URL.Path, "/order/")
	log.Printf("Url path = %s", uid)
	if order, ok := data[uid]; ok {
		w.WriteHeader(http.StatusOK)
		bytesOrder, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshalling in handler 'OrderHandler'. %s", err)
		}
		_, _ = fmt.Fprintf(w, string(bytesOrder))
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "Order with id %s not found", uid)
	}

}
