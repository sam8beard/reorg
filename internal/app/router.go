package app

import (
	"github.com/gorilla/mux"
)

/*
Builds router for development

Instead of relying on the SPA fallback, let Vite handle
serving the frontend
*/
func (s *Server) BuildRouterDev() *mux.Router {
	// Initialize router
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", s.HealthHandler).Methods("GET")
	// Get user data
	r.HandleFunc("/user", s.UserHandler).Methods("POST")
	// Receive file uploads
	r.HandleFunc("/upload", s.UploadHandler).Methods("POST")
	// Fetch files
	r.HandleFunc("/files", s.FileHandler).Methods("POST")
	// Preview organized file structure
	r.HandleFunc("/organize/preview", s.Preview).Methods("POST")
	/*
		Add all handlers here...
	*/

	return r
}

/*
Builds router for production
*/
func (s *Server) BuildRouter() *mux.Router {
	// Initialize router
	r := mux.NewRouter()

	// Add health check
	r.HandleFunc("/health", s.HealthHandler).Methods("GET")
	// Add user account endpoint
	r.HandleFunc("/user", s.UserHandler).Methods("POST")
	/*
		Add all handlers here...
	*/

	// Mount files under the base bath and add SPA handler
	spa := SpaHandler{staticPath: "dist", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	return r
}
