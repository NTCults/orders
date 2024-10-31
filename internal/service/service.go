package service

import (
	"net/http"
	"orders/internal/store"
	"orders/internal/utils"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service struct {
	Router *mux.Router
	store  store.Store
}

func NewService(db store.Store) *Service {
	r := mux.NewRouter()

	service := &Service{
		Router: r,
		store:  db,
	}
	r.HandleFunc("/orders", service.getOrdersList).Methods("GET")
	r.HandleFunc("/orders/{orderUID:[0-9]+}", service.getOrder).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	r.Use(utils.LoggingMiddleware)

	return service
}

func (s *Service) getOrdersList(w http.ResponseWriter, r *http.Request) {
	orders := s.store.GetAllOrders()
	utils.ResponseJSON(w, http.StatusOK, orders)
}

func (s *Service) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["orderUID"]

	order, err := s.store.GetOrder(orderUID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if order == nil {
		utils.ResponseError(w, http.StatusNoContent, "no content")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, order)
}
