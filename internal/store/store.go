package store

import (
	"orders/internal/config"
	"orders/models"
	"time"

	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type DBConnector interface {
	writeOrder(order *models.Order) error
	getOrder(orderUID string) (*models.Order, error)
	getAllOrders() ([]*models.Order, error)
	closeConnection() error
}

type Store struct {
	db       DBConnector
	cache    *cache.Cache
	cacheTTL time.Duration
}

func NewStore(cfg *config.Config) *Store {
	dBConnector := NewDBConnector(cfg)

	store := &Store{
		db:       dBConnector,
		cacheTTL: cfg.CacheTTL,
		cache:    cache.New(cfg.CacheTTL, cfg.CacheCleanupInterval),
	}

	if err := store.populateCache(); err != nil {
		logrus.Fatal("Populate cache on start error.")
	}

	return store
}

func (s *Store) Close() {
	s.db.closeConnection()
}

func (s *Store) populateCache() error {
	orders, err := s.db.getAllOrders()
	if err != nil {
		return err
	}

	for _, o := range orders {
		s.cache.Set(o.OrderUID, o, s.cacheTTL)
	}
	return nil
}

func (s *Store) GetAllOrders() []*models.Order {
	orders := []*models.Order{}
	for _, i := range s.cache.Items() {
		order := i.Object.(*models.Order)
		orders = append(orders, order)
	}

	return orders
}

func (s *Store) GetOrder(orderUID string) (*models.Order, error) {
	orderCached, ok := s.cache.Get(orderUID)
	if ok {
		order := orderCached.(*models.Order)
		return order, nil
	}

	order, err := s.db.getOrder(orderUID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, nil
	}

	if err := s.cache.Add(order.OrderUID, order, s.cacheTTL); err != nil {
		logrus.WithField("order_uid", order.OrderUID).Error("unable to set order to cache")
		return order, nil
	}

	return order, nil
}

func (s *Store) SetOrder(order *models.Order) error {
	if err := s.db.writeOrder(order); err != nil {
		return err
	}

	if err := s.cache.Add(order.OrderUID, order, s.cacheTTL); err != nil {
		logrus.WithField("order_uid", order.OrderUID).Error("unable to set order to cache")
	}
	return nil
}
