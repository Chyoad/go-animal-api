package domain

import (
	"errors"
	"time"
)

// Error constants
var (
	ErrAnimalNotFound      = errors.New("animal not found")
	ErrAnimalAlreadyExists = errors.New("animal with this ID already exists") // Untuk POST dengan ID duplikat
	ErrInvalidAnimalData   = errors.New("invalid animal data")
)

type Animal struct {
	ID        int       `json:"id" gorm:"primaryKey;unique"`
	Name      string    `json:"name" gorm:"not null"`
	Class     string    `json:"class"`
	Legs      int       `json:"legs"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AnimalRepository interface {
	Create(animal *Animal) error            // Gagal jika ID sudah ada
	Upsert(animal *Animal) (*Animal, error) // Membuat atau memperbarui berdasarkan animal.ID
	Delete(id int) error 
	FindByID(id int) (*Animal, error)
	FindAll() ([]Animal, error)
}