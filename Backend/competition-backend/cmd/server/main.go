package main

import (
	"log"
	"net/http"
	"time"

	"competition-backend/internal/attempts/round1"
	"competition-backend/internal/attempts/round2"
	"competition-backend/internal/auth"
	"competition-backend/internal/competition"
	"competition-backend/internal/config"
	"competition-backend/internal/middleware"
	"competition-backend/internal/storage"
)

func main() {
	cfg := config.Load()

	store := storage.New(
		cfg.PocketBaseURL,
		cfg.PocketBaseAdminEmail,
		cfg.PocketBaseAdminPassword,
	)

	competitionHandler := competition.NewHandler(store)
	round1Handler := round1.NewHandler(store)
	round2Handler := round2.NewHandler(store) // Initialize Round 2 handler

	authMiddleware := auth.Middleware(store.PocketBase)

	mux := http.NewServeMux()

	// =========================
	// ROUND 1
	// =========================

	mux.Handle(
		"POST /api/round1/attempt",
		authMiddleware(http.HandlerFunc(round1Handler.StartAttempt)),
	)
	mux.Handle(
		"GET /api/round1/question",
		authMiddleware(http.HandlerFunc(round1Handler.GetQuestion)),
	)
	mux.Handle(
		"POST /api/round1/answer",
		authMiddleware(http.HandlerFunc(round1Handler.SubmitAnswer)),
	)

	// =========================
	// ROUND 2
	// =========================
	mux.Handle(
		"POST /api/round2/attempt",
		authMiddleware(http.HandlerFunc(round2Handler.StartAttempt)),
	)
	mux.Handle(
		"GET /api/round2/state",
		authMiddleware(http.HandlerFunc(round2Handler.GetState)),
	)
	mux.Handle(
		"GET /api/round2/question",
		authMiddleware(http.HandlerFunc(round2Handler.GetQuestion)),
	)
	mux.Handle(
		"POST /api/round2/submit",
		authMiddleware(http.HandlerFunc(round2Handler.SubmitAnswer)),
	)
	mux.Handle(
		"POST /api/round2/run",
		authMiddleware(http.HandlerFunc(round2Handler.RunCode)),
	)

	// =========================
	// COMPETITIONS
	// =========================

	mux.HandleFunc(
		"GET /api/competitions",
		competitionHandler.List,
	)

	mux.HandleFunc(
		"GET /api/competitions/{id}",
		competitionHandler.Get,
	)

	// =========================
	// PUBLIC HEALTH CHECK
	// =========================

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{"status":"ok"}`))
	})

	// =========================
	// POCKETBASE HEALTH CHECK
	// =========================

	mux.HandleFunc("GET /api/pb-health", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := store.PocketBase.HealthCheck(ctx)

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)

			w.Write([]byte(
				`{"status":"error","pocketbase":"unreachable"}`,
			))

			return
		}

		w.WriteHeader(http.StatusOK)

		w.Write([]byte(
			`{"status":"ok","pocketbase":"reachable"}`,
		))
	})

	// =========================
	// PROTECTED USER ENDPOINT
	// =========================

	mux.Handle(
		"GET /api/me",
		authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user, ok := auth.GetUser(r.Context())

			if !ok {
				http.Error(
					w,
					"user not found",
					http.StatusInternalServerError,
				)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			w.Write([]byte(`{
				"id": "` + user.ID + `",
				"name": "` + user.Name + `",
				"email": "` + user.Email + `",
				"verified": ` + boolString(user.Verified) + `
			}`))
		})),
	)

	// =========================
	// HTTP SERVER
	// =========================

	corsMiddleware := middleware.CORS(cfg.CORSOrigins)

	server := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           corsMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server starting on %s", cfg.ServerAddr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// =========================
// BOOLEAN → JSON
// =========================

func boolString(value bool) string {
	if value {
		return "true"
	}

	return "false"
}
