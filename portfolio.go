package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/about", getAbout)

	err := http.ListenAndServe(":9066", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed.")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET `/` request")
	io.WriteString(w, "hello, world!\n")
}

func getAbout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET '/about' request")
	io.WriteString(w, "about me\n")
}
