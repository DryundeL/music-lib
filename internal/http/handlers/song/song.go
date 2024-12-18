package song

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"music-lib/internal/lib/api/response"
	"music-lib/internal/models"
	"music-lib/internal/storage/pgsql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SongHandlers struct {
	storage *pgsql.Storage
	logger  *slog.Logger
}

type ResponseList struct {
	response.Response
	Songs []models.Song `json:"songs,omitempty"`
}

type ResponseSingle struct {
	response.Response
	Song models.Song `json:"song,omitempty"`
}

func NewSongHandlers(storage *pgsql.Storage, logger *slog.Logger) *SongHandlers {
	return &SongHandlers{storage: storage, logger: logger}
}

// List возвращает список всех песен
func (h *SongHandlers) List(w http.ResponseWriter, r *http.Request) {
	var songs []models.Song
	if err := h.storage.DB.Find(&songs).Error; err != nil {
		h.logger.Error("failed to list songs", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, ResponseList{
		Response: response.OK(),
		Songs:    songs,
	})
}

type RequestCreate struct {
	Name     string `json:"name" validate:"required,name"`
	ArtistID uint   `json:"artist_id"`
}

// Create создает новую песню
func (h *SongHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req RequestCreate

	if err := validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		h.logger.Error("failed to validate request", err)

		render.JSON(w, r, response.ValidationError(validateErr))

		return
	}

	song := models.Song{Name: req.Name, ArtistID: req.ArtistID}
	if err := h.storage.DB.Create(&song).Error; err != nil {
		h.logger.Error("failed to create song", slog.Any("error", err))
		render.JSON(w, r, response.Error("no able to create song"))
		return
	}

	render.JSON(w, r, ResponseSingle{
		Response: response.OK(),
		Song:     song,
	})
}

// Get возвращает песню по ID
func (h *SongHandlers) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var song models.Song
	if err := h.storage.DB.First(&song, id).Error; err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, ResponseSingle{
		Response: response.OK(),
		Song:     song,
	})
}

type RequestUpdate struct {
	Name string `json:"name" validate:"omitempty,name"`
}

// Update обновляет данные о песне
func (h *SongHandlers) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var req RequestUpdate

	if err := validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		h.logger.Error("failed to validate request", err)

		render.JSON(w, r, response.ValidationError(validateErr))

		return
	}

	var song models.Song
	if err := h.storage.DB.First(&song, id).Error; err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	song.Name = req.Name
	if err := h.storage.DB.Save(&song).Error; err != nil {
		h.logger.Error("failed to update song", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, ResponseSingle{
		Response: response.OK(),
		Song:     song,
	})
}

// Delete удаляет песню по ID
func (h *SongHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.storage.DB.Delete(&models.Song{}, id).Error; err != nil {
		h.logger.Error("failed to delete song", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, ResponseSingle{
		Response: response.OK(),
	})
}
