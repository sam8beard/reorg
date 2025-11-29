package utils

import (
	"github.com/gorilla/mux"
	api "github.com/sam8beard/reorg/internal/app"
	"net/http"
	"os"
	"path/filepath"
)

/*
Allows SPA serving by giving the path to
the static directory and the path to the index.html
file within that static directory
*/
type spaHandler struct {
	staticPath string
	indexPath  string
}

/*
Serves file located within static directory if found
If no file is found, serves file at the index path
*/
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
Builds router for development

Instead of relying on the SPA fallback, let Vite handle
serving the frontend
*/
func BuildRouterDev() *mux.Router {
	// Initialize router
	r := mux.NewRouter()

	// Add health check
	r.HandleFunc("/health", api.HealthHandler).Methods("GET")

	/*
		Add all handlers here...
	*/

	return r
}

/*
Builds router for production
*/
func BuildRouter() *mux.Router {
	// Initialize router
	r := mux.NewRouter()

	// Add health check
	r.HandleFunc("/health", api.HealthHandler).Methods("GET")

	/*
		Add all handlers here...
	*/

	// Mount files under the base bath and add SPA handler
	spa := spaHandler{staticPath: "dist", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	return r
}
