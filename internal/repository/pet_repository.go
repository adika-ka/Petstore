package repository

import (
	"context"
	"database/sql"
	"fmt"
	"petstore/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PetRepository interface {
	Create(ctx context.Context, pet model.Pet) (model.Pet, error)
	Update(ctx context.Context, pet model.Pet) (model.Pet, error)
	UpdateFormData(ctx context.Context, petID int, name, status string) (model.Pet, error)
	FindByID(ctx context.Context, petID int) (model.Pet, error)
	FindByStatus(ctx context.Context, statuses []string) ([]model.Pet, error)
	FindByTags(ctx context.Context, tags []string) ([]model.Pet, error)
	Delete(ctx context.Context, petID int) error
	ExistsByID(ctx context.Context, petID int) (bool, error)
}

type petRepo struct {
	db *sqlx.DB
}

func NewPetRepository(db *sqlx.DB) PetRepository {
	return &petRepo{db: db}
}

func (r *petRepo) Create(ctx context.Context, pet model.Pet) (model.Pet, error) {
	petDB, err := model.PetToPetDB(pet)
	if err != nil {
		return pet, fmt.Errorf("failed to convert pet to db: %w", err)
	}

	query := `
		INSERT INTO pets (name, status, category, photo_urls, tags)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	var newID int
	err = r.db.QueryRowContext(ctx, query,
		petDB.Name,
		petDB.Status,
		petDB.Category,
		petDB.PhotoUrls,
		petDB.Tags,
	).Scan(&newID)
	if err != nil {
		return pet, fmt.Errorf("failed to insert pet: %w", err)
	}

	pet.ID = newID
	return pet, nil
}

func (r *petRepo) Update(ctx context.Context, pet model.Pet) (model.Pet, error) {
	petDB, err := model.PetToPetDB(pet)
	if err != nil {
		return pet, fmt.Errorf("failed to convert pet to db: %w", err)
	}

	query := `
		UPDATE pets
		SET name=$1, status=$2, category=$3, photo_urls=$4, tags=$5
		WHERE id=$6
	`
	_, err = r.db.ExecContext(ctx, query,
		petDB.Name,
		petDB.Status,
		petDB.Category,
		petDB.PhotoUrls,
		petDB.Tags,
		petDB.ID,
	)
	if err != nil {
		return pet, fmt.Errorf("failed to update pet: %w", err)
	}

	return pet, nil
}

func (r *petRepo) UpdateFormData(ctx context.Context, petID int, name, status string) (model.Pet, error) {
	query := `
		UPDATE pets
		SET name=$1, status=$2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, name, status, petID)
	if err != nil {
		return model.Pet{}, fmt.Errorf("failed to update pet form data: %w", err)
	}

	pet, err := r.FindByID(ctx, petID)
	if err != nil {
		return model.Pet{}, fmt.Errorf("failed to retrieve updated pet: %w", err)
	}

	return pet, nil
}

func (r *petRepo) FindByID(ctx context.Context, petID int) (model.Pet, error) {
	query := `
		SELECT id, name, status, category, photo_urls, tags
		FROM pets WHERE id = $1
	`
	var petDB model.PetDB
	err := r.db.GetContext(ctx, &petDB, query, petID)
	if err != nil {
		return model.Pet{}, fmt.Errorf("failed to find pet by id: %w", err)
	}

	return model.PetDBToPet(petDB)
}

func (r *petRepo) FindByStatus(ctx context.Context, statuses []string) ([]model.Pet, error) {
	query := `
		SELECT id, name, status, category, photo_urls, tags
		FROM pets
		WHERE status = ANY($1)
	`

	var petDBs []model.PetDB
	err := r.db.SelectContext(ctx, &petDBs, query, pq.Array(statuses))
	if err != nil {
		return nil, fmt.Errorf("failed to select pets by status: %w", err)
	}

	pets := make([]model.Pet, 0, len(petDBs))
	for _, petDB := range petDBs {
		pet, err := model.PetDBToPet(petDB)
		if err != nil {
			return nil, fmt.Errorf("error mapping pet: %w", err)
		}
		pets = append(pets, pet)
	}

	return pets, nil
}

func (r *petRepo) FindByTags(ctx context.Context, tags []string) ([]model.Pet, error) {
	query := `
		SELECT id, name, status, category, photo_urls, tags
		FROM pets
	`
	var petDBs []model.PetDB
	err := r.db.SelectContext(ctx, &petDBs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select pets: %w", err)
	}

	pets := make([]model.Pet, 0, len(petDBs))
	for _, petDB := range petDBs {
		pet, err := model.PetDBToPet(petDB)
		if err != nil {
			return nil, fmt.Errorf("error mapping pet: %w", err)
		}

		// Шаг 3: фильтруем по тегам
		if hasMatchingTag(pet.Tags, tags) {
			pets = append(pets, pet)
		}
	}

	return pets, nil
}

func (r *petRepo) Delete(ctx context.Context, petID int) error {
	query := `DELETE FROM pets WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, petID)
	if err != nil {
		return fmt.Errorf("failed to delete pet: %w", err)
	}
	return nil
}

func (r *petRepo) ExistsByID(ctx context.Context, petID int) (bool, error) {
	query := `SELECT 1 FROM pets WHERE id = $1 LIMIT 1`

	var dummy int
	err := r.db.QueryRowContext(ctx, query, petID).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func hasMatchingTag(petTags []model.Tag, searchTags []string) bool {
	for _, search := range searchTags {
		for _, petTag := range petTags {
			if petTag.Name == search {
				return true
			}
		}
	}
	return false
}
