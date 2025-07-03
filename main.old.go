package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Wassssup"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	msg := fmt.Sprintf("a mf snippet with id %d}", id)
	w.Write([]byte(msg))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Wassssssssup. snippet creating form"))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Wassssssssup. creating a snippet"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home)

	mux.HandleFunc("GET /snippets/view/{$}", snippetView)
	mux.HandleFunc("GET /snippets/create", snippetCreate)

	mux.HandleFunc("GET /snippets/create", snippetCreatePost)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error(err.Error(), "err", err)
		os.Exit(1)
	}
}
