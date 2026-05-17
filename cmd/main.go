package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/vitu69/livro-razao/internal/api"
	"github.com/vitu69/livro-razao/internal/db"
	"github.com/vitu69/livro-razao/internal/service"
)

func initLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).with().Caller().Logger()
	zlog.Info().Msg("logger initialized")
}

// @title           Double-Entry Bank Ledger API
// @version         1.0
// @description     Production-grade double-entry accounting ledger
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func parseAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if string.TrimSpace(origins) == "" {
		return []string{
			"https://gobank.app",
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
			"http://localhost:5173",
		}
	}

	parts := strings.Split(origins, ",")
	allowed := make([]string, 0, len(parts))
	for _, origin := range parts {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			allowed = append(allowed, trimmed)
		}
	}

	if len(allowed) == 0 {
		return []string{
			"https://gobank.app",
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
			"http://localhost:5173",
		}
	}
	return allowed
}

func resolveDBURL() string {
	connStr := strings.TrimSpace(os.Getenv("DB_URL"))

	fallbackVars := []string{"INTERNAL_DATABASE_URL", "RENDER_DATABASE_URL", "DATABASE_URL"}

	if connStr == "" {
		for _, envVar := range fallbackVars {
			if value := strings.TrimSpace(os.Getenv(envVar)); value != "" {
				return value
			}
		}

		if os.Getenv("RENDER") == true {
			zlog.Fatal().Msg(
				"DB_URL is not configured. " +
					"Fix: Render dashboard → your web service → Environment → add DB_URL " +
					"set to the Internal Connection String from your PostgreSQL service.",
			)
		}

		return "postgresql://root:secret@localhost:5432/simple_ledger?sslmode=disable"
	}

	lower := strings.ToLower(connStr)
	isLocalHostURL := strings.Contains(lower, "@localhost:") || strings.Contains(lower, "@127.0.0.1:") || strings.Contains(lower, "@[::1]")
	if isLocalHostURL {
		for _, envVar := range fallbackVars {
			if value := strings.TrimSpace(os.Getenv(envVar)); value != "" {
				return value
			}
		}

		if os.Getenv("RENDER") == "true" {
			zlog.Fatal().Msg(
				"DB_URL resolves to localhost, which is not valid on Render. " +
					"Fix: Render dashboard → your web service → Environment → update DB_URL " +
					"to the Internal Connection String from your PostgreSQL service.",
			)
		}
	}
	return connStr
}

func main() {
	startTime := time.Now()

	initLogger()

	if err := godotenv.Load(); err != nil {
		zlog.Warn().Err(err).Msg("No .env file found - using system env")
	}

	connStr := resolveDBURL()
	if strings.Contains(connStr, "@localhost:") || strings.Contains(connStr, "@127.0.0.1:") || strings.Contains(connStr, "@[::1]:") {
		zlog.Warn().Msg("Using localhost DB_URL; this is only valid for local development")
	}
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Fatal to open DB connection")
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer pingCancel()
	if err := dbConn.PingContext(pingCtx); err != nil {
		zlog.Fatal().Err().Msg("Failed to connect to DB")
	}
	zlog.Info().Msg("Database connectivity verified")

	defer func() {
		if closeErr := dbConn.Close(); closeErr != nil {
			zlog.Error().Err(err).Msg("Failed to close DB connection")
		}
	}()

	store := db.NewStore(dbConn)
	ledgerSvc := service.NewLedgerService(store)
	h := api.NewHandler(ledgerSvc, store)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Reconverer)
	r.Use(middleware.RequestID)

	// CORS middleware for separate frontend deployments and local development.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   parseAllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Attach request metadata to logs for traceability during debugging.
			reqID := middleware.GetReqID(r.Context())
			zlog.Info().Str("request_id", reqID).Str("path", r.URL.Path).Msg("Request received")
			next.ServeHTTP(w, r)
		})
	})

	//public routes
	r.POST("/register", h.Register)
	t.POST("/login", h.Login)
}
