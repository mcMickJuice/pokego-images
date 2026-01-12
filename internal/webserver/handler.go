package webserver

import (
	"fmt"
	"net/http"
)

func CreateWebserver() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/pokemon/{pokemon}", func(w http.ResponseWriter, r *http.Request) {

		param := r.PathValue("pokemon")
		fmt.Printf("pokemon param: %s", param)

		response := "hello"
		_, err := w.Write([]byte(response))
		if err != nil {

			fmt.Printf("error writing response: %v", err)
		}
	})

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {

		return err
	}

	return nil
}
