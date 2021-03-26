package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"time"
)

// local command flags
var port int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve current stored data over http",
	Long: `Serves current database currency data over http, see README for endpoint descriptions`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port for serving API endpoints")
}


func serve() {
	router := mux.NewRouter().StrictSlash(true)
	// main routes
	router.HandleFunc("/", withLogging(healthHandler))
	router.HandleFunc("/new", withLogging(newestHandler))
	router.HandleFunc("/history/{currency}", withLogging(historyHandler))

	// not found route
	router.NotFoundHandler = router.NewRoute().BuildOnly().HandlerFunc(withLogging(http.NotFound)).GetHandler()

	// start server
	log.Printf("Listening on port %v", port)
	http.ListenAndServe(":8080", router)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "curfetch v1.0 @ %v", time.Now())
}

func newestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "newest @ %v", time.Now())
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "history for %s operational @ %v", mux.Vars(r)["currency"], time.Now())
}

// withLogging - middleware - logs incoming request to stdout
func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s:%s  FROM:%s",r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}