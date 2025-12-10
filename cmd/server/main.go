package main

import (
	"embed"
	"fmt"
	"github.com/sam8beard/reorg/internal/db/pgsql"
	"github.com/sam8beard/reorg/internal/obj-store/minio"
	// "github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sam8beard/reorg/cmd/server/utils"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

/* Development vars */
var devDir string

/* Production vars */
var prodDir string
var prodFS fs.FS

// //go:embed "dist/*"
var dist embed.FS

/*
Runs the development flow
*/
func runDev() {
	// Configure logger for debugging
	log.SetFlags(log.Lshortfile)

	// Load path to frontend files for development
	devDir = os.Getenv("DEV_FRONTEND_DIR")
	port := ":8080"

	// Confirm development environment is correct
	if _, err := os.Stat(devDir); err != nil {
		msg := fmt.Sprintf("development directory is invalid: %v", err)
		log.Fatal(msg)
	}

	// Create router with all handlers
	router := utils.BuildRouterDev()
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

/*
Runs the production flow
*/
func runProd() {
	// Configure logger for debugging
	log.SetFlags(log.Lshortfile)

	// Load path to dist files for production
	prodDir = os.Getenv("DIST_DIR")
	// Get dist subtree if we want to deploy using
	// a binary
	//distSub, _ := fs.Sub(dist, "dist")
	//prodSub := http.FS(distSub)
	port := ":8081"

	// Confirm production environment is correct
	if _, err := os.Stat(prodDir); err != nil {
		msg := fmt.Sprintf("production directory is invalid: %v", err)
		log.Fatal(msg)
	}

	// Create router with all handlers and SPA fallback
	router := utils.BuildRouter()
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
	// Create server
	//	srv := &http.Server{
	//		Addr: port,
	//	}

	// Run
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Panicf("could not load environment variables: %v", err)
	}

	db := pgsql.Init()
	minio := minio.Init()

	runDev()
	//runProd()
}
