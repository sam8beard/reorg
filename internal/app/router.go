package app

import (
	"github.com/gorilla/mux"
	"github.com/sam8beard/reorg/internal/auth/middleware"
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

	// Public routes
	r.HandleFunc("/auth/signup", s.SignupHandler).Methods("POST")
	r.HandleFunc("/auth/login", s.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/guest", s.GuestHandler).Methods("POST")

	// Auth protected routes
	authMiddleware := middleware.AuthMiddleware(s.JWTService)
	protected := r.NewRoute().Subrouter()
	protected.Use(authMiddleware)

	// Get user data
	//protected.HandleFunc("/user", s.UserHandler).Methods("POST")

	// Receive file uploads
	protected.HandleFunc("/upload", s.UploadHandler).Methods("POST")

	// Fetch files
	protected.HandleFunc("/files", s.FileHandler).Methods("POST")

	// Preview organized file structure
	protected.HandleFunc("/organize/preview", s.PreviewHandler).Methods("POST")

	// Receive target data
	protected.HandleFunc("/target", s.TargetHandler).Methods("POST")

	// Receive rule data
	protected.HandleFunc("/rule", s.RuleHandler).Methods("POST")

	// Return preview object based on ruleset
	protected.HandleFunc("/preview", s.PreviewHandler).Methods("POST")

	// Return zip file of evaluation result
	protected.HandleFunc("/download/zip", s.DownloadZipHandler).Methods("POST")

	// Write user file organization info to database
	protected.HandleFunc("/organize/save", s.SaveOrg).Methods("POST")

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
	//r.HandleFunc("/user", s.UserHandler).Methods("POST")
	/*
		Add all handlers here...
	*/

	// Mount files under the base bath and add SPA handler
	spa := SpaHandler{staticPath: "dist", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	return r
}
