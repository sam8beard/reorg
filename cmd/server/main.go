package main

import (
	"embed"
	"github.com/joho/godotenv"
	"io/fs"
	"log"
	"net/http"
	"os"
)

/* Development vars */
var devDir string
var devFS http.Dir
var devPort string

/* Production vars */
var prodDir string
var prodFS fs.FS

//go:embed "dist/*"
var dist embed.FS
var prodPort string

/*
Runs the development flow using the specified port and file system
*/
func runDev(port string, fSys http.Dir) {

	if err := http.ListenAndServe(
		port, http.FileServer(fSys),
	); err != nil {
		log.Panicf("could not serve frontend from Go server: %v", err)
	}
	log.Printf("listening on port %s", port)
}

/*
Runs the production flow using the specified port and file system
*/
func runProd(port string, fSys fs.FS) {
	// Run production
	if err := http.ListenAndServe(
		port, http.FileServerFS(fSys),
	); err != nil {
		log.Panicf("could not serve frontend from Go server: %v", err)
	}
}
func main() {
	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Panicf("could not load environment: %v", err)
	}

	// Load path to frontend files for development
	devDir = os.Getenv("DEV_FRONTEND_DIR")
	devFS = http.Dir(devDir)
	devPort = ":8080"

	// Load path to dist files for production
	prodDir = os.Getenv("DIST_DIR")
	// Get dist subtree
	distSub, _ := fs.Sub(dist, "dist")
	prodFS = fs.FS(distSub)
	prodPort = ":8081"

	// Configure logger for debugging
	log.SetFlags(log.Lshortfile)

	// Confirm development environment is correct
	if _, err := os.Stat(devDir); err != nil {
		log.Fatalf("development directory is invalid: %v", err)
	}

	// Confirm production environment is correct
	if _, err := os.Stat(prodDir); err != nil {
		log.Fatalf("development directory is invalid: %v", err)
	}

	runDev(devPort, devFS)
	//runProd(prodPort, prodFS)

	/*

		// Run development
		if err := http.ListenAndServe(
			devPort, http.FileServer(devFS),
		); err != nil {
			log.Panicf("could not serve frontend from Go server: %v", err)
		}
		log.Printf("listening on port %s", port)

	*/

	// Run production
	//	if err := http.ListenAndServe(
	//		prodPort, http.FileServerFS(prodFS),
	//	); err != nil {
	//		log.Panicf("could not serve frontend from Go server: %v", err)
	//	}

}
