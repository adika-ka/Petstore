package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"petstore/infrastructure"
	"petstore/internal/model"
	"petstore/internal/service"

	"github.com/go-chi/chi"
)

type UserController struct {
	Service   service.UserService
	Responder infrastructure.Responder
}

func RegisterUserRoutes(r chi.Router, uc *UserController) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/", addUser(uc))
		r.Route("/{username}", func(r chi.Router) {
			r.Get("/", getUserByUsername(uc))
			r.Put("/", updateUser(uc))
			r.Delete("/", deleteUser(uc))
		})
		r.Post("/createWithList", addListUsers(uc))
		r.Post("/createWithArray", addListUsers(uc))
		r.Get("/login", uc.Login)
		r.Get("/logout", logout())
	})
}

// AddUser godoc
// @Summary      Create user
// @Description  This can only be done by the logged in user.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body body model.User true "Created user object"
// @Success      200 {object} model.ApiResponse "successful operation"
// @Router       /user [post]
func addUser(uc *UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u model.User

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			uc.Responder.ErrorBadRequest(w, err)
			return
		}

		if err := validateUser(r.Context(), uc.Service, u); err != nil {
			log.Printf("Error validate user: %v", err)
			uc.Responder.ErrorBadRequest(w, err)
			return
		}

		user, err := uc.Service.CreateUser(r.Context(), u)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			uc.Responder.ErrorInternal(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	}
}

// CreateUsersWithList godoc
// @Summary      Creates list of users with given input array
// @Description  Creates list of users with given input array
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body body []model.User true "List of user object"
// @Success      200 {object} model.ApiResponse "successful operation"
// @Router       /user/createWithList [post]
// @Router       /user/createWithArray [post]
func addListUsers(uc *UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u []model.User

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			uc.Responder.ErrorBadRequest(w, err)
			return
		}

		users, err := uc.Service.CreateUserBatch(r.Context(), u)
		if err != nil {
			log.Printf("Error creating users: %v", err)
			uc.Responder.ErrorInternal(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(users)
	}
}

// GetUserByUsername godoc
// @Summary      Get user by user name
// @Description  The name that needs to be fetched. Use user1 for testing.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        username path string true "The name that needs to be fetched"
// @Success      200 {object} model.User "successful operation"
// @Router       /user/{username} [get]
func getUserByUsername(uc *UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		user, err := uc.Service.FindUserByUsername(r.Context(), username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				uc.Responder.ErrorNotFound(w, fmt.Errorf("user not found"))
				return
			}
			log.Printf("Error finding user by username %s: %v", username, err)
			uc.Responder.ErrorInternal(w, err)
			return
		}

		uc.Responder.OutputJSON(w, user)
	}
}

// UpdateUser godoc
// @Summary      Updated user
// @Description  This can only be done by the logged in user.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        username path string true "name that need to be updated"
// @Param        user body model.User true "Updated user object"
// @Success      200 {object} model.User "successful operation"
// @Router       /user/{username} [put]
func updateUser(uc *UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			uc.Responder.ErrorBadRequest(w, err)
			return
		}

		updatedUser, err := uc.Service.UpdateUser(r.Context(), username, u)
		if err != nil {
			log.Printf("Error updating user: %v", err)
			uc.Responder.ErrorInternal(w, err)
			return
		}

		uc.Responder.OutputJSON(w, updatedUser)
	}
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  This can only be done by the logged in user.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        username path string true "The name that needs to be deleted"
// @Success      200 {object} model.ApiResponse "successful operation"
// @Router       /user/{username} [delete]
func deleteUser(uc *UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		if err := uc.Service.DeleteUser(r.Context(), username); err != nil {
			log.Printf("Error deleting user by username %s: %v", username, err)
			uc.Responder.ErrorInternal(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// LoginUser godoc
// @Summary      Logs user into the system
// @Description  Logs user into the system
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        username query string true "The user name for login"
// @Param        password query string true "The password for login in clear text"
// @Success      200 {string} string "token"
// @Router       /user/login [get]
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username == "" || password == "" {
		uc.Responder.ErrorBadRequest(w, fmt.Errorf("missing username or password"))
		return
	}

	token, err := uc.Service.Login(r.Context(), username, password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			uc.Responder.ErrorUnauthorized(w, err)
			return
		}
		uc.Responder.ErrorInternal(w, err)
		return
	}

	uc.Responder.OutputJSON(w, map[string]string{"token": token})
}

// LogoutUser godoc
// @Summary      Logs out current logged in user session
// @Description  Logs out current logged in user session
// @Tags         user
// @Accept       json
// @Produce      json
// @Success 200 {string} string "ok"
// @Router       /user/logout [get]
func logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func validateUser(ctx context.Context, service service.UserService, user model.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}

	_, err := service.FindUserByUsername(ctx, user.Username)
	if err == nil {
		return errors.New("username already exist")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return fmt.Errorf("error checking username: %w", err)
}
