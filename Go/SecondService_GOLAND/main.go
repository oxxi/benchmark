package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//generate UUID
		uuid := uuid.NewString()
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(uuid))
	})

	//listener
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error inicializando el servicio")
		os.Exit(1)
	}
}
