package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"petstore/infrastructure"
	"petstore/internal/model"
	"petstore/internal/service"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type PetController struct {
	Service   service.PetService
	Responder infrastructure.Responder
}

func RegisterPetRoutes(r chi.Router, pc *PetController) {
	r.Route("/pet", func(r chi.Router) {
		r.Post("/", addPet(pc))
		r.Put("/", updatePet(pc))
		r.Get("/findByStatus", getPetsByStatus(pc))
		r.Get("/findByTags", getPetsByTags(pc))

		r.Route("/{petId}", func(r chi.Router) {
			r.Get("/", getPetByID(pc))
			r.Post("/", updatePetForm(pc))
			r.Delete("/", deletePet(pc))

			r.Route("/uploadImage", func(r chi.Router) {
				r.Post("/", uploadPetImage(pc))
			})
		})
	})
}

// @Summary Add a new pet to the store
// @Tags pet
// @Accept json
// @Produce json
// @Param pet body model.Pet true "Pet object that needs to be added to the store"
// @Success 200 {object} model.Pet
// @Security ApiKeyAuth
// @Router /pet [post]
func addPet(pc *PetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p model.Pet

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			pc.Responder.ErrorBadRequest(w, err)
			return
		}

		if err := service.ValidatePet(p); err != nil {
			pc.Responder.ErrorBadRequest(w, err)
			return
		}

		pet, err := pc.Service.CreatePet(r.Context(), p)
		if err != nil {
			log.Printf("Error creating pet %v: %v", pet, err)
			pc.Responder.ErrorInternal(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pet)
	}
}

// UpdatePet godoc
// @Summary      Update an existing pet
// @Description  Update an existing pet in the store
// @Tags         pet
// @Accept       json
// @Produce      json
// @Param        pet  body      model.Pet  true  "Pet object that needs to be added to the store"
// @Success      200  {object}  model.Pet
// @Security ApiKeyAuth
// @Router       /pet [put]
func updatePet(pc *PetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p model.Pet

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			pc.Responder.ErrorBadRequest(w, err)
			return
		}

		if err := service.ValidatePet(p); err != nil {
			pc.Responder.ErrorBadRequest(w, err)
			return
		}

		pet, err := pc.Service.UpdatePet(r.Context(), p)
		if err != nil {
			log.Printf("Error updating pet: %v", err)
			pc.Responder.ErrorInternal(w, err)
			return
		}

		pc.Responder.OutputJSON(w, pet)
	}
}

// FindPetByStatus godoc
// @Summary      Finds Pets by status
// @Description  Multiple status values can be provided with comma separated strings
// @Tags         pet
// @Accept       json
// @Produce      json
// @Param        status query []string true "Status values that need to be considered for filter" Enums(available, pending, sold)
// @Success      200 {array} model.Pet "successful operation"
// @Security ApiKeyAuth
// @Router       /pet/findByStatus [get]
func getPetsByStatus(pc *PetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		if status == "" {
			pc.Responder.ErrorBadRequest(w, fmt.Errorf("status query parameter is required"))
			return
		}

		statuses := strings.Split(status, ",")

		pets, err := pc.Service.FindPetByStatus(r.Context(), statuses)
		if err != nil {
			log.Printf("Error finding pets by status %s: %v", status, err)
			pc.Responder.ErrorInternal(w, err)
			return
		}

		if len(pets) == 0 {
			pc.Responder.ErrorNotFound(w, fmt.Errorf("no pets found for given statuses"))
			return
		}

		pc.Responder.OutputJSON(w, pets)
	}
}

// FindPetByTags godoc
// @Summary      Finds Pets by tags
// @Description  Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.
// @Tags         pet
// @Accept       json
// @Produce      json
// @Param        tags query []string true "Tags to filter by"
// @Success      200 {array} model.Pet "successful operation"
// @Router       /pet/findByTags [get]
// @Security ApiKeyAuth
// @Deprecated
func getPetsByTags(pc *PetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagsParam := r.URL.Query().Get("tags")
		if tagsParam == "" {
			pc.Responder.ErrorBadRequest(w, fmt.Errorf("tags query parameter is required"))
			return
		}

		tags := strings.Split(tagsParam, ",")

		pets, err := pc.Service.FindPetByTags(r.Context(), tags)
		if err != nil {
			log.Printf("Error finding pets by tags %v: %v", tags, err)
			pc.Responder.ErrorInternal(w, err)
			return
		}

		if len(pets) == 0 {
			pc.Responder.ErrorNotFound(w, fmt.Errorf("no pets found for given tags"))
			return
		}

		pc.Responder.OutputJSON(w, pets)
	}
}

// GetPetByID godoc
// @Summary      Find pet by ID
// @Description  Returns a single pet
// @Tags         pet
// @Accept       json
// @Produce      json
// @Param        petId path int true "ID of pet to return"
// @Success      200 {object} model.Pet "successful operation"
// @Security ApiKeyAuth
// @Router       /pet/{petId} [get]
func getPetByID(pc *PetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		petIDStr := chi.URLParam(r, "petId")

		petID, err := strconv.Atoi(petIDStr)
		if err != nil {
			pc.Responder.ErrorBadRequest(w, fmt.Errorf("invalid pet ID"))
			return
		}

		pet, err := pc.Service.FindPetByID(r.Context(), petID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				pc.Responder.ErrorNotFound(w, fmt.Errorf("pet not found"))
				return
			}
			log.Printf("Error finding pet by ID %d: %v", petID, err)
			pc.Responder.ErrorInternal(w, err)
			return
		}

		pc.Responder.OutputJSON(w, pet)
	}
}

// UpdatePetWithForm godoc
// @Summary      Updates a pet in the store with form data
// @Description  Updates name and status of pet
// @Tags         pet
// @Accept       multipart/form-data
// @Produce      json
// @Param        petId path int true "ID of pet that needs to be updated"
// @Param        name formData string false "Updated name of the pet"
// @Param        status formData string false "Updated status of the pet"
// @Success      200 {object} model.Pet "successful operation"
// @Security ApiKeyAuth
// @Router       /pet/{petId} [post]
func updatePetForm(pc *PetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		petIDStr := chi.URLParam(r, "petId")

		petID, err := strconv.Atoi(petIDStr)
		if err != nil {
			pc.Responder.ErrorBadRequest(w, fmt.Errorf("invalid pet ID"))
			return
		}

		name := r.FormValue("name")
		status := r.FormValue("status")

		if err := service.ValidatePetFormData(name, status); err != nil {
			pc.Responder.ErrorBadRequest(w, err)
			return
		}

		pet, err := pc.Service.UpdatePetFormData(r.Context(), petID, name, status)
		if err != nil {
			log.Printf("Error updating pet form data: %v", err)
			pc.Responder.ErrorInternal(w, err)
			return
		}

		pc.Responder.OutputJSON(w, pet)
	}
}

// DeletePet godoc
// @Summary      Deletes a pet
// @Description  Deletes a pet by ID
// @Tags         pet
// @Accept       json
// @Produce      json
// @Param        petId path int true "Pet id to delete"
// @Success      200 {object} model.ApiResponse "successful operation"
// @Security ApiKeyAuth
// @Router       /pet/{petId} [delete]
func deletePet(pc *PetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		petIDStr := chi.URLParam(r, "petId")

		id, err := strconv.Atoi(petIDStr)
		if err != nil {
			pc.Responder.ErrorBadRequest(w, fmt.Errorf("invalid pet ID"))
			return
		}

		if err := pc.Service.DeletePet(r.Context(), id); err != nil {
			log.Printf("Error deleting pet ID %d: %v", id, err)
			pc.Responder.ErrorInternal(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary uploads an image
// @Tags pet
// @Accept multipart/form-data
// @Produce json
// @Param petId path int true "ID of pet to update"
// @Param additionalMetadata formData string false "Additional data to pass to server"
// @Param file formData file true "File to upload"
// @Success 200 {object} model.ApiResponse
// @Security ApiKeyAuth
// @Router /pet/{petId}/uploadImage [post]
func uploadPetImage(pc *PetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		petIDStr := chi.URLParam(r, "petId")

		id, err := strconv.Atoi(petIDStr)
		if err != nil {
			pc.Responder.ErrorBadRequest(w, fmt.Errorf("invalid pet ID"))
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			pc.Responder.ErrorBadRequest(w, fmt.Errorf("file is required"))
			return
		}
		defer file.Close()

		pc.Responder.OutputJSON(w, map[string]string{
			"message": fmt.Sprintf("Image uploaded for pet ID %d (stub)", id),
		})
	}
}
