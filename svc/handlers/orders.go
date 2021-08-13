package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-flexible/example-order-system/svc/domain"
	"github.com/gorilla/mux"
)

func createOrder(deps Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order domain.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctrl := domain.OrderController{
			Dependencies: deps,
			Nats:         deps.Nats(), // uses a different intereface within domain.
		}

		if err := ctrl.New(r.Context(), order); err != nil {
			log.Println(err)
			httpError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func getOrderByID(deps Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			http.Error(w, "missing order id", http.StatusBadRequest)
			return
		}

		ctrl := domain.OrderController{Dependencies: deps}
		order, err := ctrl.GetOrderByID(r.Context(), id)
		if err != nil {
			log.Println(err)
			httpError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func httpError(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case domain.DatabaseQueryError:
		log.Println(err.Stmt)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	case domain.PublishingError:
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
