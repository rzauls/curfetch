package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rzauls/curfetch/db"
	"github.com/spf13/cobra"
)

// local command flags
var port int

// Server - struct for passing around shared resources
type Server struct {
	db db.Storage
}

// NewServeCmd represents the serve command
func NewServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Serve current stored data over http",
		Long:  `Serves current database currency data over http, see README for endpoint descriptions`,
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}
}

// init - initialize command and its flags
func init() {
	serveCmd := NewServeCmd()
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port for serving API endpoints")
}

// serve - main command action
func serve() {
	// set up db connection
	session, err := db.NewSession()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer session.Close()

	// initialize router
	s := Server{db: db.NewStorage(session)}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", withLogging(s.healthHandler))
	router.HandleFunc("/newest", withLogging(s.newestHandler))
	router.HandleFunc("/history/{currency}", withLogging(s.historyHandler))

	// not found route
	router.NotFoundHandler = router.NewRoute().BuildOnly().HandlerFunc(withLogging(http.NotFound)).GetHandler()

	// start server
	log.Printf("Listening on port %v", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), router); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "curfetch v1.0 @ %v", time.Now())
}

func (s *Server) newestHandler(w http.ResponseWriter, r *http.Request) {
	// fetch data from db
	rows, err := s.db.Newest()
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	// serve data
	w.Header().Set("Content-Type", "application/json")
	if len(rows) > 0 {
		err = json.NewEncoder(w).Encode(rows)
		if err != nil {
			log.Fatalf("Failed to serve data: %v", err)
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "no data")
	}
}

func (s *Server) historyHandler(w http.ResponseWriter, r *http.Request) {
	currencyCode := strings.ToUpper(mux.Vars(r)["currency"])
	// fetch data from db
	rows, err := s.db.History(currencyCode)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	// serve data
	w.Header().Set("Content-Type", "application/json")
	if len(rows) > 0 {
		err = json.NewEncoder(w).Encode(rows)
		if err != nil {
			log.Fatalf("Failed to serve data: %v", err)
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "no data for: %v", currencyCode)
	}
}

// withLogging - middleware - logs incoming request to stdout
func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s  FROM: %s (%s)", r.Method, r.RequestURI, r.RemoteAddr, r.UserAgent())
		next.ServeHTTP(w, r)
	}
}
