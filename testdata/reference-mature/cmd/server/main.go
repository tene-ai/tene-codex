package main

import (
	"github.com/tene-ai/reference-mature/internal/service"
	"net/http"
)

func main() {
	http.HandleFunc("/legacy/orders", func(w http.ResponseWriter, r *http.Request) {
		_ = service.Process(r.Context(), r.FormValue("id"))
		w.WriteHeader(202)
	})
	_ = http.ListenAndServe(":8080", nil)
}
