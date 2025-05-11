package dto

import (
	"go-animal-api/internal/domain"
	"time"
)

// AnimalPayload untuk request POST dan PUT
type AnimalPayload struct {
	ID    int    `json:"id" binding:"required,gt=0"` // ID wajib dan harus positif
	Name  string `json:"name" binding:"required"`
	Class string `json:"class"`
	Legs  int    `json:"legs" binding:"gte=0"` // Kaki harus >= 0
}

// AnimalResponse untuk format response animal.
type AnimalResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Class     string    `json:"class"`
	Legs      int       `json:"legs"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToAnimalResponse mengubah domain.Animal menjadi AnimalResponse DTO.
func ToAnimalResponse(animal *domain.Animal) AnimalResponse {
	if animal == nil { // Safety check
		return AnimalResponse{}
	}
	return AnimalResponse{
		ID:        animal.ID,
		Name:      animal.Name,
		Class:     animal.Class,
		Legs:      animal.Legs,
		CreatedAt: animal.CreatedAt,
		UpdatedAt: animal.UpdatedAt,
	}
}

// ToAnimalListResponse mengubah slice domain.Animal menjadi slice AnimalResponse DTO.
func ToAnimalListResponse(animals []domain.Animal) []AnimalResponse {
	list := make([]AnimalResponse, len(animals))
	for i, animal := range animals {
		list[i] = ToAnimalResponse(&animal)
	}
	return list
}