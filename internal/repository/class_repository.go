package repository

import (
	"glofox/internal/constants"
	"glofox/internal/models"
	"sync"
)

type ClassRepository interface {
	Create(class models.Class) error
	GetByName(name string) (models.Class, bool)
}

// ClassRepo manages the in-memory class data
type ClassRepo struct {
	classes map[string]models.Class
	mu      sync.RWMutex
}

// NewClassRepo creates a new ClassRepo
func NewClassRepo() *ClassRepo {
	return &ClassRepo{
		classes: make(map[string]models.Class),
	}
}

// Create for creating a new class
func (classRepo *ClassRepo) Create(class models.Class) error {
	classRepo.mu.Lock()
	defer classRepo.mu.Unlock()

	if _, exists := classRepo.classes[class.Name]; exists {
		return constants.ErrClassAlreadyExists
	}
	classRepo.classes[class.Name] = class
	return nil
}

// GetByName fetches class by given name
func (classRepo *ClassRepo) GetByName(name string) (models.Class, bool) {
	classRepo.mu.RLock()
	defer classRepo.mu.RUnlock()

	class, exists := classRepo.classes[name]
	return class, exists
}
