package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-flexible/example-order-system/svc/domain"
)

func createOrder(deps Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order domain.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctrl := domain.OrderController{Dependencies: deps}

		if err := ctrl.New(r.Context(), order); err != nil {
			log.Println(err)
			switch err := err.(type) {
			case domain.DatabaseQueryError:
				log.Println(err.Stmt)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			case domain.PublishingError:
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
