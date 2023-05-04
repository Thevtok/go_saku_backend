package main

import (
	"net/http"

	"github.com/ReygaFitra/inc-final-project.git/delivery"
)

func main() {
	http.HandleFunc("/", Handler)
	delivery.RunServer()
}
func Handler(w http.ResponseWriter, r *http.Request) {
	delivery.RunServer()
}
