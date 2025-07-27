package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", helloHandler)

	addr := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello, world!"))
}
