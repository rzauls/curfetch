package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rzauls/curfetch/db"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// local command flags
var port int

// db handler
var DB db.CurrencyModel

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

func init() {
	serveCmd := NewServeCmd()
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port for serving API endpoints")
}

func serve() {
	// set up db connection
	cluster := db.InitDB(db.CassandraConfig{
		Hosts:    []string{os.Getenv("CASS_HOST")}, // potentially you can pass multiple cassandra nodes here
		Keyspace: "curfetch",
	})
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer session.Close()
	DB = db.CurrencyModel{Session: session}

	router := mux.NewRouter().StrictSlash(true)
	// main routes
	router.HandleFunc("/", withLogging(healthHandler))
	router.HandleFunc("/newest", withLogging(newestHandler))
	router.HandleFunc("/history/{currency}", withLogging(historyHandler))

	// not found route
	router.NotFoundHandler = router.NewRoute().BuildOnly().HandlerFunc(withLogging(http.NotFound)).GetHandler()

	// start server
	log.Printf("Listening on port %v", port)
	if err := http.ListenAndServe(":" + strconv.Itoa(port), router); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "curfetch v1.0 @ %v", time.Now())
}

func newestHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Newest()
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	if len(rows) > 0 {
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rows)
		if err != nil {
			log.Fatalf("Failed to serve data: %v", err)
		}
	} else {
		fmt.Fprintf(w, "no data")
	}
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	currencyCode := strings.ToUpper(mux.Vars(r)["currency"])
	rows, err := DB.History(currencyCode)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	if len(rows) > 0 {
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rows)
		if err != nil {
			log.Fatalf("Failed to serve data: %v", err)
		}
	} else {
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
