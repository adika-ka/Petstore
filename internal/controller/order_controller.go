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

	"github.com/go-chi/chi"
)

type OrderController struct {
	Service   service.OrderService
	Responder infrastructure.Responder
}

func RegisterOrderRoutes(r chi.Router, oc *OrderController) {
	r.Route("/store/order", func(r chi.Router) {
		r.Post("/", addOrder(oc))
		r.Route("/{orderId}", func(r chi.Router) {
			r.Get("/", getOrderByID(oc))
			r.Delete("/", deleteOrder(oc))
		})
	})
}

// CreateOrder godoc
// @Summary      Place an order for a pet
// @Description  Places a new order in the system
// @Tags         store
// @Accept       json
// @Produce      json
// @Param        order body model.Order true "order placed for purchasing the pet"
// @Success      200 {object} model.Order
// @Router       /store/order [post]
func addOrder(oc *OrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var o model.Order

		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			oc.Responder.ErrorBadRequest(w, err)
			return
		}

		order, err := oc.Service.CreateOrder(r.Context(), o)
		if err != nil {
			log.Printf("Error creating order %v: %v", order, err)
			oc.Responder.ErrorInternal(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}

// GetOrderById godoc
// @Summary      Find purchase order by ID
// @Description  For valid response try integer IDs with value >= 1 and <= 10. Other values will generate exceptions
// @Tags         store
// @Accept       json
// @Produce      json
// @Param        orderId path int true "ID of pet that needs to be fetched"
// @Success      200 {object} model.Order
// @Router       /store/order/{orderId} [get]
func getOrderByID(oc *OrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderIDStr := chi.URLParam(r, "orderId")

		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			oc.Responder.ErrorBadRequest(w, fmt.Errorf("invalid order ID"))
			return
		}

		order, err := oc.Service.FindOrderByID(r.Context(), orderID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				oc.Responder.ErrorNotFound(w, fmt.Errorf("order not found"))
				return
			}
			log.Printf("Error finding order by ID %d: %v", orderID, err)
			oc.Responder.ErrorInternal(w, err)
			return
		}

		oc.Responder.OutputJSON(w, order)
	}
}

// DeleteOrder godoc
// @Summary      Delete purchase order by ID
// @Description  For valid response try integer IDs with positive integer value. Negative or non-integer values will generate API errors
// @Tags         store
// @Accept       json
// @Produce      json
// @Param        orderId path int true "ID of the order that needs to be deleted"
// @Success      200 {object} model.ApiResponse "successful operation"
// @Router       /store/order/{orderId} [delete]
func deleteOrder(oc *OrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderIDStr := chi.URLParam(r, "orderId")

		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			oc.Responder.ErrorBadRequest(w, fmt.Errorf("invalid order ID"))
			return
		}

		if err := oc.Service.DeleteOrder(r.Context(), orderID); err != nil {
			log.Printf("Error deleting order ID %d: %v", orderID, err)
			oc.Responder.ErrorInternal(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetInventory godoc
// @Summary      Returns pet inventories by status
// @Description  Returns a map of status codes to quantities
// @Tags         store
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]int "successful operation"
// @Security     ApiKeyAuth
// @Router       /store/inventory [get]
func GetInventory(oc *OrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inventory, err := oc.Service.GetInventory(r.Context())
		if err != nil {
			log.Printf("Error finding inventory")
			oc.Responder.ErrorInternal(w, err)
			return
		}

		oc.Responder.OutputJSON(w, inventory)
	}
}
