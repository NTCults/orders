package store

import (
	"orders/models"
	"testing"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"
)

type mockDBConnector struct {
	writeOrderCallsCount int
	getOrderCallsCount   int
}

func (m *mockDBConnector) writeOrder(order *models.Order) error {
	m.writeOrderCallsCount++
	return nil
}
func (m *mockDBConnector) getOrder(orderUID string) (*models.Order, error) {
	m.getOrderCallsCount++
	return nil, nil
}
func (m *mockDBConnector) getAllOrders() ([]*models.Order, error) {
	return nil, nil
}
func (m *mockDBConnector) closeConnection() error {
	return nil
}

func TestStore(t *testing.T) {
	mockdbConn := &mockDBConnector{}
	store := Store{
		db:    mockdbConn,
		cache: cache.New(-1, -1),
	}

	order := &models.Order{OrderUID: "1"}
	store.SetOrder(order)

	require.Equal(t, 1, mockdbConn.writeOrderCallsCount, "must be 1")

	order, err := store.GetOrder("1")
	require.NotEmpty(t, order)
	require.Equal(t, "1", order.OrderUID)
	require.NoError(t, err)
	require.Empty(t, mockdbConn.getOrderCallsCount)
}
