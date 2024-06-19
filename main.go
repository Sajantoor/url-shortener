package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := http.Server{
		Addr:    ":3000",
		Handler: http.HandlerFunc(requestHandler),
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("Error starting server: ", err)
	}

	fmt.Println("Server running on port 3000")
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
