package main

import (
	"log"
	"net/http"

	"github.com/2twin/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const rootDir = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootDir))))

	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", healthzHandler)
	apiRouter.Get("/chirps", apiCfg.retrieveChirpsHandler)
	apiRouter.Post("/chirps", apiCfg.createChirpsHandler)
	apiRouter.Get("/chirps/{chirpID}", apiCfg.getChirpHandler)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.metricsHandler)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", rootDir, port)
	log.Fatal(srv.ListenAndServe())
}
