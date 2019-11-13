package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	httpPort := 13000

	log.Println("Server available at following address:")
	log.Printf("    http://localhost:%d/", httpPort)

	handler := NewRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), logRequest(handler)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
