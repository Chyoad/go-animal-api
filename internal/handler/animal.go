package handler

import (
	"errors"
	"go-animal-api/internal/domain"
	"go-animal-api/internal/dto"
	"go-animal-api/internal/usecase"
	"go-animal-api/pkg/utils/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AnimalHandler struct {
	useCase usecase.AnimalUseCase
}

func NewAnimalHandler(uc usecase.AnimalUseCase) *AnimalHandler {
	return &AnimalHandler{useCase: uc}
}

// POST /v1/animal
func (h *AnimalHandler) CreateAnimal(c *gin.Context) {
	var payload dto.AnimalPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	createdAnimal, err := h.useCase.CreateAnimal(&payload)
	if err != nil {
		if errors.Is(err, domain.ErrAnimalAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil) // Sesuai requirement "duplicate entry should be denied"
		} else if errors.Is(err, domain.ErrInvalidAnimalData) { // Dari validasi DTO binding jika ada
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to create animal", err)
		}
		return
	}
	response.Success(c, http.StatusCreated, "Animal created successfully", dto.ToAnimalResponse(createdAnimal))
}

// PUT /v1/animal (Upsert berdasarkan ID di payload)
func (h *AnimalHandler) UpsertAnimal(c *gin.Context) {
	var payload dto.AnimalPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	upsertedAnimal, err := h.useCase.UpsertAnimal(&payload)
	if err != nil {
		// ErrInvalidAnimalData jika ada validasi tambahan di usecase sebelum ke repo.
		if errors.Is(err, domain.ErrInvalidAnimalData) {
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			// Error lain dari database atau proses upsert
			response.Error(c, http.StatusInternalServerError, "Failed to upsert animal", err)
		}
		return
	}

	statusCode := http.StatusOK // Default ke OK untuk update.
	if upsertedAnimal.CreatedAt.Equal(upsertedAnimal.UpdatedAt) || time.Since(upsertedAnimal.CreatedAt) < 2*time.Second && time.Since(upsertedAnimal.UpdatedAt) < 2*time.Second {
		if upsertedAnimal.CreatedAt.Equal(upsertedAnimal.UpdatedAt) {
			 statusCode = http.StatusCreated
		}
	}
	response.Success(c, statusCode, "Animal upserted successfully", dto.ToAnimalResponse(upsertedAnimal))
}

// DELETE /v1/animal/:id
func (h *AnimalHandler) DeleteAnimal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // Ubah ke Atoi untuk tipe int
	if err != nil || id <= 0 {     // Pastikan ID positif
		response.Error(c, http.StatusBadRequest, "Invalid or non-positive animal ID format in path", err)
		return
	}

	err = h.useCase.DeleteAnimal(id)
	if err != nil {
		if errors.Is(err, domain.ErrAnimalNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
		} else if errors.Is(err, domain.ErrInvalidAnimalData) { // Dari validasi ID di usecase
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to delete animal", err)
		}
		return
	}
	response.Success(c, http.StatusOK, "Animal deleted successfully", nil)
}

// GET /v1/animal/:id
func (h *AnimalHandler) GetAnimalByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, "Invalid or non-positive animal ID format in path", err)
		return
	}

	animal, err := h.useCase.GetAnimalByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrAnimalNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
		} else if errors.Is(err, domain.ErrInvalidAnimalData) {
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to get animal", err)
		}
		return
	}
	response.Success(c, http.StatusOK, "Animal retrieved successfully", dto.ToAnimalResponse(animal))
}

// GET /v1/animal
func (h *AnimalHandler) GetAllAnimals(c *gin.Context) {
	animals, err := h.useCase.GetAllAnimals()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get animals", err)
		return
	}

	if len(animals) == 0 {
		response.Error(c, http.StatusNotFound, "No animals found", nil)
		return
	}
	response.Success(c, http.StatusOK, "Animals retrieved successfully", dto.ToAnimalListResponse(animals))
}