package usecase

import (
	"go-animal-api/internal/domain"
	"go-animal-api/internal/dto"
)

type AnimalUseCase interface {
	CreateAnimal(payload *dto.AnimalPayload) (*domain.Animal, error)
	UpsertAnimal(payload *dto.AnimalPayload) (*domain.Animal, error)
	DeleteAnimal(id int) error
	GetAnimalByID(id int) (*domain.Animal, error)
	GetAllAnimals() ([]domain.Animal, error)
}

type animalUseCase struct {
	repo domain.AnimalRepository
}

func NewAnimalUseCase(repo domain.AnimalRepository) AnimalUseCase {
	return &animalUseCase{repo: repo}
}

func (uc *animalUseCase) CreateAnimal(payload *dto.AnimalPayload) (*domain.Animal, error) {
	// Validasi (ID > 0, Legs >= 0, Name required) di dto.
	animal := &domain.Animal{
		ID:    payload.ID,
		Name:  payload.Name,
		Class: payload.Class,
		Legs:  payload.Legs,
	}

	err := uc.repo.Create(animal)
	if err != nil {
		return nil, err
	}
	
	return animal, nil
}

func (uc *animalUseCase) UpsertAnimal(payload *dto.AnimalPayload) (*domain.Animal, error) {
	// Validasi
	animal := &domain.Animal{
		ID:    payload.ID,
		Name:  payload.Name,
		Class: payload.Class,
		Legs:  payload.Legs,
	}
	return uc.repo.Upsert(animal)
}

func (uc *animalUseCase) DeleteAnimal(id int) error {
	// Validasi ID
	if id <= 0 {
		return domain.ErrInvalidAnimalData
	}
	return uc.repo.Delete(id)
}

func (uc *animalUseCase) GetAnimalByID(id int) (*domain.Animal, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidAnimalData
	}
	return uc.repo.FindByID(id)
}

func (uc *animalUseCase) GetAllAnimals() ([]domain.Animal, error) {
	return uc.repo.FindAll()
}