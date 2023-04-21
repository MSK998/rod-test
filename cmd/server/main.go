package main

import (
	"log"
	"net/http"

	//"rod-test/pkg/database"
	_ "statement-service-poc/pkg/handler"
	"statement-service-poc/pkg/router"
)

func main() {
	/*
		err := database.Open("sqlserver", "<ConnectionString>")
		if err != nil {
			log.Fatal(err.Error())
		}
	*/

	r := router.New()
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err.Error())
	}
}
