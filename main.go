/*
 * Power control block
 *
 * This API was created to monitor states of Power Control Block and send some commands to it.
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package main

import (
	"log"
	"net/http"
	"time"

	"./control"

	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//

	sw "./go"
)

func main() {
	log.Printf("Hello programm")

	time.Sleep(time.Second)
	control.Init()

	log.Printf("Server started")

	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
