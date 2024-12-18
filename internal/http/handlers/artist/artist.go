package artist

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

type ArtistHandlers struct {
	storage *pgsql.Storage
	logger  *slog.Logger
}

type ResponseList struct {
	response.Response
	Artists []models.Artist `json:"artists,omitempty"`
}

type ResponseSingle struct {
	response.Response
	Artist models.Artist `json:"artist,omitempty"`
}

func NewArtistHandlers(storage *pgsql.Storage, logger *slog.Logger) *ArtistHandlers {
	return &ArtistHandlers{storage: storage, logger: logger}
}

// List возвращает список всех артистов
func (h *ArtistHandlers) List(w http.ResponseWriter, r *http.Request) {
	var artists []models.Artist
	if err := h.storage.DB.Find(&artists).Error; err != nil {
		h.logger.Error("failed to list artists", slog.Any("error", err))
		render.JSON(w, r, response.Error("fail to load artists"))
		return
	}

	render.JSON(w, r, ResponseList{
		Response: response.OK(),
		Artists:  artists,
	})
}

type RequestCreate struct {
	Name    string `json:"name" validate:"required,name"`
	IsGroup bool   `json:"is_group" validate:"omitempty"`
}

// Create создает нового артиста
func (h *ArtistHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req RequestCreate

	if err := validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		h.logger.Error("failed to validate request", err)

		render.JSON(w, r, response.ValidationError(validateErr))

		return
	}

	artist := models.Artist{Name: req.Name, IsGroup: req.IsGroup}
	if err := h.storage.DB.Create(&artist).Error; err != nil {
		h.logger.Error("failed to create artist", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, ResponseSingle{
		Response: response.OK(),
		Artist:   artist,
	})
}

// Get возвращает артиста по ID
func (h *ArtistHandlers) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var artist models.Artist
	if err := h.storage.DB.First(&artist, id).Error; err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, ResponseSingle{
		Response: response.OK(),
		Artist:   artist,
	})
}

type RequestUpdate struct {
	Name    string `json:"name" validate:"omitempty,name"`
	IsGroup bool   `json:"is_group" validate:"omitempty,is_group"`
}

// Update обновляет данные об артисте по ID
func (h *ArtistHandlers) Update(w http.ResponseWriter, r *http.Request) {
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

	var artist models.Artist
	if err := h.storage.DB.First(&artist, id).Error; err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	artist.Name = req.Name
	artist.IsGroup = req.IsGroup
	if err := h.storage.DB.Save(&artist).Error; err != nil {
		h.logger.Error("failed to update artist", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, ResponseSingle{
		response.OK(),
		artist,
	})
}

// Delete удаляет артиста по ID
func (h *ArtistHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.storage.DB.Delete(&models.Artist{}, id).Error; err != nil {
		h.logger.Error("failed to delete artist", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, ResponseSingle{
		Response: response.OK(),
	})
}
