package app

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/sam8beard/reorg/internal/auth"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Server struct {
	DB          *pgxpool.Pool
	Minio       *minio.Client
	MinioBucket string
	JWTService  *auth.JWTService
}

/*
Allows SPA serving by giving the path to
the static directory and the path to the index.html
file within that static directory
*/
type SpaHandler struct {
	staticPath string
	indexPath  string
}

func NewServer(jwtService *auth.JWTService, db *pgxpool.Pool, minio *minio.Client) *Server {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("could not load env vars: %v", err)
	}

	return &Server{
		DB:          db,
		Minio:       minio,
		MinioBucket: os.Getenv("MINIO_BUCKET_NAME"),
		JWTService:  jwtService,
	}
}

/*
Serves file located within static directory if found
If no file is found, serves file at the index path
*/
func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Joins static file path with request URL
	path := filepath.Join(h.staticPath, r.URL.Path)

	// Check if file exists or is a directory
	file, err := os.Stat(path)

	// If file does not exist or file is a directory
	if os.IsNotExist(err) || file.IsDir() {
		// File does not exist or is a directory so serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	// Other error besides file does not exist
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serve static file
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

/*
Runs the development flow
*/
func (s *Server) RunDev() {
	// Configure logger for debugging
	log.SetFlags(log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("could not load env vars: %v", err)
	}

	// Load path to frontend files for development
	devDir := os.Getenv("DEV_FRONTEND_DIR")
	port := ":8080"

	// Confirm development environment is correct
	if _, err := os.Stat(devDir); err != nil {
		msg := fmt.Sprintf("development directory is invalid: %v", err)
		log.Fatal(msg)
	}

	// Create router with all handlers
	router := s.BuildRouterDev()
	// Create server
	srv := &http.Server{
		Handler:      router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Run
	if err := srv.ListenAndServe(); err != nil {
		msg := fmt.Sprintf("could not serve frontend from Go server: %v", err)
		log.Fatal(msg)
	}

	log.Printf("listening on port %s", port)
}
