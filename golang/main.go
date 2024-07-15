package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/kawabatas/auth0-react-go-sample/middleware"
)

func public(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Hello from a public endpoint"}`))
}

func private(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Hello from a private endpoint"}`))
}

func rbac(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Hello from a private rbac endpoint"}`))
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	audience := os.Getenv("AUTH0_AUDIENCE")
	domain := os.Getenv("AUTH0_DOMAIN")
	frontend := os.Getenv("FRONTEND_URL")
	port := "8080"

	router := http.NewServeMux()

	publicHandler := http.HandlerFunc(public)
	privateHandler := http.HandlerFunc(private)
	rbacHandler := http.HandlerFunc(rbac)
	router.Handle("/api/public",
		middleware.UseAccessLog(
			middleware.CORS(frontend, publicHandler),
		),
	)
	router.Handle("/api/private",
		middleware.UseAccessLog(
			middleware.CORS(frontend,
				middleware.ValidateJWT(audience, domain, privateHandler),
			),
		),
	)
	router.Handle("/api/privaterbac",
		middleware.UseAccessLog(
			middleware.CORS(frontend,
				middleware.ValidateJWT(audience, domain,
					middleware.ValidatePermissions([]string{"read:messages"}, rbacHandler),
				),
			),
		),
	)

	slog.Info("Server listening on http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		slog.Error("There was an error with the http server: " + err.Error())
	}
}
