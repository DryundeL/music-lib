package router

import (
	"music-lib/internal/http/handlers/artist"
	"music-lib/internal/http/handlers/song"
	"net/http"

	"log/slog"
	mvLog "music-lib/internal/http/middleware/logger"
	"music-lib/internal/storage/pgsql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New создаёт новый Router с подключенными хэндлерами.
// Параметры:
// - storage: экземпляр вашего pgsql хранилища
// - logger: ваш логгер для логирования запросов и ошибок
func New(storage *pgsql.Storage, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(mvLog.New(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte("OK")); err != nil {
			logger.Error("failed to write health check response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	artistHandlers := artist.NewArtistHandlers(storage, logger)
	songHandlers := song.NewSongHandlers(storage, logger)

	r.Route("/artists", func(r chi.Router) {
		r.Get("/", artistHandlers.List)          // GET /artists
		r.Post("/", artistHandlers.Create)       // POST /artists
		r.Get("/{id}", artistHandlers.Get)       // GET /artists/{id}
		r.Put("/{id}", artistHandlers.Update)    // PUT /artists/{id}
		r.Delete("/{id}", artistHandlers.Delete) // DELETE /artists/{id}
	})

	r.Route("/songs", func(r chi.Router) {
		r.Get("/", songHandlers.List)          // GET /songs
		r.Post("/", songHandlers.Create)       // POST /songs
		r.Get("/{id}", songHandlers.Get)       // GET /songs/{id}
		r.Put("/{id}", songHandlers.Update)    // PUT /songs/{id}
		r.Delete("/{id}", songHandlers.Delete) // DELETE /songs/{id}
	})

	return r
}
