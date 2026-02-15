package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"mon-go/internal/handler"
	"mon-go/internal/store"
)

// Route binds a path to a handler and which HTTP methods to register.
// Methods is http.MethodGet (register both GET and POST) or http.MethodPost (register only POST).
type Route struct {
	Path    string
	Methods string // http.MethodGet or http.MethodPost
	Handler http.HandlerFunc
}

// Server holds the HTTP server and dependencies.
type Server struct {
	httpServer *http.Server
	db         *store.DB
}

// itemRoutes is the hardcoded table of item routes and their handlers.
func itemRoutes(h *handler.ItemHandler) []Route {
	return []Route{
		{"/items.create", http.MethodGet, h.CreateItem},
		{"/items.get", http.MethodGet, h.GetItem},
		{"/items.list", http.MethodGet, h.ListItems},
		{"/items.delete", http.MethodGet, h.DeleteItem},
	}
}

// objectMemberRoutes is the hardcoded table of object-members routes and their handlers.
func objectMemberRoutes(h *handler.ObjectMemberHandler) []Route {
	return []Route{
		{"/object-members.create", http.MethodGet, h.Create},
		{"/object-members.delete", http.MethodGet, h.Delete},
	}
}

// registerRoutes registers all routes in the given slice (GET and/or POST per route).
func registerRoutes(r chi.Router, routes []Route) {
	for _, route := range routes {
		switch route.Methods {
		case http.MethodGet:
			r.Get(route.Path, route.Handler)
			r.Post(route.Path, route.Handler)
		case http.MethodPost:
			r.Post(route.Path, route.Handler)
		}
	}
}

// New builds the router and HTTP server. Call Run() to start listening.
func New(port int, db *store.DB, itemHandler *handler.ItemHandler, objectMemberHandler *handler.ObjectMemberHandler) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("Hello World"))
	})

	registerRoutes(r, itemRoutes(itemHandler))
	registerRoutes(r, objectMemberRoutes(objectMemberHandler))

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
		db: db,
	}
}

// Run starts the HTTP server. Blocks until the server is shut down.
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server and closes the DB (if connected).
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	if s.db != nil {
		return s.db.Close(ctx)
	}
	return nil
}
