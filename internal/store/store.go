package store

import (
	"database/sql"
	"orders/internal/config"
	"orders/models"
	"time"

	"github.com/avast/retry-go"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Store interface {
	GetAllOrders() []*models.Order
	GetOrder(orderUID string) (*models.Order, error)
	SetOrder(order *models.Order) error
	Close()
}

type PGStore struct {
	db       *sql.DB
	cache    cache.Cache
	cacheTTL time.Duration
}

func NewPGStore(cfg *config.Config) Store {
	var db *sql.DB
	err := retry.Do(
		func() error {
			var err error
			db, err = sql.Open("postgres", cfg.DBConnString)
			if err != nil {
				logrus.Info("Trying to connect to db.")
				return err
			}
			return nil
		},
		retry.Attempts(5),
		retry.Delay(time.Second),
	)

	if err != nil {
		logrus.WithField("DB_CONN_STR", cfg.DBConnString).Fatal("Unable to connect to db.")
	}

	store := &PGStore{
		db:       db,
		cacheTTL: cfg.CacheTTL,
		cache:    *cache.New(cfg.CacheTTL, cfg.CacheCleanupInterval),
	}

	if err := store.populateCache(); err != nil {
		logrus.Fatal("Populate cache on start error.")
	}

	return store
}

func (s *PGStore) Close() {
	s.db.Close()
}

func (s *PGStore) populateCache() error {
	orders, err := s.getAllOrdersFromDB()
	if err != nil {
		return err
	}

	for _, o := range orders {
		s.cache.Set(o.OrderUID, o, s.cacheTTL)
	}
	return nil
}

func (s *PGStore) GetAllOrders() []*models.Order {
	orders := []*models.Order{}
	for _, i := range s.cache.Items() {
		order := i.Object.(*models.Order)
		orders = append(orders, order)
	}

	return orders
}

func (s *PGStore) GetOrder(orderUID string) (*models.Order, error) {
	orderCached, ok := s.cache.Get(orderUID)
	if ok {
		order := orderCached.(*models.Order)
		return order, nil
	}

	order, err := s.getOrderFromDB(orderUID)
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

func (s *PGStore) SetOrder(order *models.Order) error {
	if err := s.writeOrderTX(order); err != nil {
		return err
	}

	if err := s.cache.Add(order.OrderUID, order, s.cacheTTL); err != nil {
		logrus.WithField("order_uid", order.OrderUID).Error("unable to set order to cache")
	}
	return nil
}
