package main

import (
	"fmt"
	"mcmickjuice/pokego/internal/webserver"
)

func main() {
	fmt.Println("starting webserver")
	err := webserver.CreateWebserver()
	if err != nil {
		fmt.Printf("error starting webserver: %v", err)
	}
	fmt.Println("exiting webserver")
}
