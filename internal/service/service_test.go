package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"orders/models"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type MockStore struct {
	store map[string]*models.Order
}

func (s *MockStore) GetAllOrders() []*models.Order {
	return []*models.Order{}
}
func (s *MockStore) GetOrder(orderUID string) (*models.Order, error) {
	order, ok := s.store[orderUID]
	if ok {
		return order, nil
	}

	return nil, nil
}
func (s *MockStore) SetOrder(order *models.Order) error {
	return nil
}
func (s *MockStore) Close() {}

func TestServiceHandlers(t *testing.T) {

	t.Run("should return 200 ok and correct data", func(t *testing.T) {
		orderUID := "123"
		req, err := http.NewRequest("GET", "/order/"+orderUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		//Hack to fake gorilla/mux vars
		vars := map[string]string{
			"orderUID": orderUID,
		}
		req = mux.SetURLVars(req, vars)

		srv := NewService(&MockStore{
			store: map[string]*models.Order{
				orderUID: {
					OrderUID: orderUID,
				},
			},
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(srv.getOrder)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "must be 200 OK")

		order := &models.Order{}
		err = json.Unmarshal(rr.Body.Bytes(), &order)
		if err != nil {
			panic(err)
		}

		require.Equal(t, orderUID, order.OrderUID, "must be equal")
	})

	t.Run("should return 204 on unexisting orderUUID", func(t *testing.T) {
		unexistingOrderUID := "531"
		req, err := http.NewRequest("GET", "/order/"+unexistingOrderUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		srv := NewService(&MockStore{})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(srv.getOrder)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNoContent, rr.Code, "must be 204")
	})
}
