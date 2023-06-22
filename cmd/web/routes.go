package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rlr524/bookings/pkg/config"
	"github.com/rlr524/bookings/pkg/handlers"
	"net/http"
	"path/filepath"
)

func Routes(_ *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/madison", handlers.Repo.Madison)

	fs := http.FileServer(neuteredFileSystem{http.Dir("./static/")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/*", http.StripPrefix("/static", fs))

	return mux
}

// Custom file system to disable directory listing
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
