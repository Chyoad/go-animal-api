package repository

import (
	"errors"
	"go-animal-api/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mysqlAnimalRepository struct {
	db *gorm.DB
}

func NewMysqlAnimalRepository(db *gorm.DB) domain.AnimalRepository {
	return &mysqlAnimalRepository{db: db}
}

// Create - POST: "Request to create a duplicate entry should be denied."
func (r *mysqlAnimalRepository) Create(animal *domain.Animal) error {
	// Cek existing animal ID
	var existing domain.Animal
	if err := r.db.First(&existing, animal.ID).Error; err == nil {
		// Jika tidak ada error, berarti record dengan ID tersebut ditemukan
		return domain.ErrAnimalAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	// Jika err adalah gorm.ErrRecordNotFound, berarti ID belum ada, lanjutkan pembuatan.
	return r.db.Create(animal).Error
}

// Upsert - PUT: "update an existing animal or create a new one if the animal doesn't exist yet"
func (r *mysqlAnimalRepository) Upsert(animal *domain.Animal) (*domain.Animal, error) {
	//  clause.OnConflict - INSERT jika primary key (ID) baru, atau UPDATE jika ID sudah ada.
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // Konflik pada kolom 'id'
		DoUpdates: clause.AssignmentColumns([]string{"name", "class", "legs", "updated_at"}), // Field yang diupdate
	}).Create(animal).Error // Create akan melakukan INSERT atau UPDATE tergantung konflik

	if err != nil {
		return nil, err
	}

	// Fetch ulang data terbaru.
	var resultAnimal domain.Animal
	if findErr := r.db.First(&resultAnimal, animal.ID).Error; findErr != nil {
		return nil, findErr
	}
	return &resultAnimal, nil
}

// Delete animal by ID.
func (r *mysqlAnimalRepository) Delete(id int) error {
	// Cek existing animal ID
	var animal domain.Animal
	if err := r.db.First(&animal, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrAnimalNotFound
		}
		return err
	}

	result := r.db.Delete(&domain.Animal{}, id)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return domain.ErrAnimalNotFound
	}
	return nil
}

// Find animal by ID.
func (r *mysqlAnimalRepository) FindByID(id int) (*domain.Animal, error) {
	var animal domain.Animal
	if err := r.db.First(&animal, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAnimalNotFound
		}
		return nil, err
	}
	return &animal, nil
}

// Find all animals.
func (r *mysqlAnimalRepository) FindAll() ([]domain.Animal, error) {
	var animals []domain.Animal
	if err := r.db.Find(&animals).Error; err != nil {
		return nil, err
	}
	return animals, nil
}