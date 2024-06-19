package main

import (
	"fmt"

	"github.com/Sajantoor/url-shortener/application"
)

func main() {
	app := application.New()

	err := app.Start()

	if err != nil {
		fmt.Println("Failed to start app: ", err)
	}

	fmt.Println("App started successfully")
}
