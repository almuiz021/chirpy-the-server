package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const fileRootPath = "."
	const port = "8080"

	apiConfig := &apiConfig{}

	mux := http.NewServeMux()

	handlerFS := http.StripPrefix("/app", http.FileServer(http.Dir(fileRootPath)))

	mux.Handle("/app/", apiConfig.middlewareMetricsInc(handlerFS))
	mux.HandleFunc("/healthz", handlerHealth)
	mux.HandleFunc("/metrics", apiConfig.countHits)
	mux.HandleFunc("/reset", apiConfig.resetHits)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) countHits(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write((fmt.Appendf([]byte{}, "Hits: %d", hits)))

}
