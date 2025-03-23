package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"petstore/internal/model"
	"petstore/internal/repository"
)

type PetService interface {
	CreatePet(ctx context.Context, pet model.Pet) (model.Pet, error)
	UpdatePet(ctx context.Context, pet model.Pet) (model.Pet, error)
	UpdatePetFormData(ctx context.Context, petID int, name, status string) (model.Pet, error)
	FindPetByID(ctx context.Context, petID int) (model.Pet, error)
	FindPetByStatus(ctx context.Context, statuses []string) ([]model.Pet, error)
	FindPetByTags(ctx context.Context, tags []string) ([]model.Pet, error)
	DeletePet(ctx context.Context, petID int) error
}

type petService struct {
	repo repository.PetRepository
}

func NewPetService(repo repository.PetRepository) PetService {
	return &petService{repo: repo}
}

func (s *petService) CreatePet(ctx context.Context, pet model.Pet) (model.Pet, error) {
	return s.repo.Create(ctx, pet)
}

func (s *petService) UpdatePet(ctx context.Context, pet model.Pet) (model.Pet, error) {
	err := ValidatePet(pet)
	if err != nil {
		return model.Pet{}, fmt.Errorf("incorrect data :%w", err)
	}
	exists, err := s.repo.ExistsByID(ctx, pet.ID)
	if err != nil {
		return model.Pet{}, fmt.Errorf("error checking pet existence: %w", err)
	}
	if !exists {
		log.Printf("Pet with ID %d not found", pet.ID)
		return model.Pet{}, fmt.Errorf("pet with ID %d not found", pet.ID)
	}
	return s.repo.Update(ctx, pet)
}

func (s *petService) UpdatePetFormData(ctx context.Context, petID int, name, status string) (model.Pet, error) {
	err := ValidatePetFormData(name, status)
	if err != nil {
		return model.Pet{}, fmt.Errorf("incorrect data :%w", err)
	}
	exists, err := s.repo.ExistsByID(ctx, petID)
	if err != nil {
		return model.Pet{}, fmt.Errorf("error checking pet existence: %w", err)
	}
	if !exists {
		return model.Pet{}, fmt.Errorf("pet with ID %d not found", petID)
	}

	return s.repo.UpdateFormData(ctx, petID, name, status)

}

func (s *petService) FindPetByID(ctx context.Context, petID int) (model.Pet, error) {
	return s.repo.FindByID(ctx, petID)
}

func (s *petService) FindPetByStatus(ctx context.Context, statuses []string) ([]model.Pet, error) {
	for _, status := range statuses {
		if err := validatePetStatus(status); err != nil {
			return nil, fmt.Errorf("invalid pet status: %w", err)
		}
	}
	pets, err := s.repo.FindByStatus(ctx, statuses)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No pets found for given statuses")
			return nil, fmt.Errorf("no pets found for given statuses")
		}
		return nil, fmt.Errorf("error finding pets: %w", err)
	}
	return pets, nil
}

func (s *petService) FindPetByTags(ctx context.Context, tags []string) ([]model.Pet, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("tag list cannot be empty")
	}
	err := validateTags(tags)
	if err != nil {
		return nil, err
	}

	pets, err := s.repo.FindByTags(ctx, tags)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No pets found for given tags")
			return nil, fmt.Errorf("no pets found for given tags")
		}
		return nil, fmt.Errorf("error finding pets: %w", err)
	}
	return pets, nil
}

func (s *petService) DeletePet(ctx context.Context, petID int) error {
	return s.repo.Delete(ctx, petID)
}

func ValidatePet(pet model.Pet) error {

	if pet.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if pet.ID <= 0 {
		return fmt.Errorf("invalid pet id")
	}
	if validatePetStatus(pet.Status) != nil {
		return fmt.Errorf("invalid pet status")
	}
	return nil
}

func ValidatePetFormData(name, status string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return validatePetStatus(status)
}

func validatePetStatus(status string) error {
	var validStatuses = map[string]struct{}{
		"available": {},
		"pending":   {},
		"sold":      {},
	}
	if _, ok := validStatuses[status]; !ok {
		return fmt.Errorf("invalid pet status")
	}
	return nil
}

func validateTags(tags []string) error {
	for _, tag := range tags {
		if tag == "" {
			return fmt.Errorf("tag cannot be empty")
		}
	}
	return nil
}
