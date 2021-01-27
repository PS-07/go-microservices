package app

import (
	"fmt"
	"net/http"

	"github.com/PS-07/go-microservices/mvc/controllers"
)

// StartApp func
func StartApp() {
	http.HandleFunc("/users", controllers.GetUser)

	fmt.Println("server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
