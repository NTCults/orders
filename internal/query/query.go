package query

import (
	"bytes"
	"encoding/gob"
	"orders/internal/config"
	"orders/internal/metrics"
	"orders/models"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

const NatsSubject = "orders."

type NatsOrdersQuery struct {
	nc           *nats.Conn
	sub          *nats.Subscription
	ordersWriter OrdersWriter
	debug        bool
}

type OrdersWriter interface {
	SetOrder(order *models.Order) error
}

func NewNatsOrdersQuery(ordersWriter OrdersWriter, cfg *config.Config) (*NatsOrdersQuery, error) {
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	noq := &NatsOrdersQuery{
		nc:           nc,
		ordersWriter: ordersWriter,
		debug:        cfg.Debug,
	}

	sub, err := nc.Subscribe(NatsSubject+"*", noq.ReceiveOrderHandler)
	if err != nil {
		return nil, err
	}

	noq.sub = sub
	return noq, nil
}

func (q *NatsOrdersQuery) Close() {
	if q.sub != nil {
		q.sub.Unsubscribe()
	}
	if q.nc != nil {
		q.nc.Close()
	}
}

func (q *NatsOrdersQuery) ReceiveOrderHandler(m *nats.Msg) {
	var order models.Order
	if q.debug {
		logrus.
			WithField("msg", m).
			WithField("func", "query.ReceiveOrderHandler").
			Debug("Received a message.")
	}

	if err := gob.NewDecoder(bytes.NewBuffer(m.Data)).Decode(&order); err != nil {
		metrics.OrdersReceiveErrors.Inc()
		logrus.WithError(err).Error("unable to decode message.")
		return
	}

	if err := order.Validte(); err != nil {
		metrics.OrdersReceiveErrors.Inc()
		logrus.WithField("order", order).WithError(err).Error("received order is not valid is not valid")
		return
	}

	if err := q.ordersWriter.SetOrder(&order); err != nil {
		metrics.OrdersReceiveErrors.Inc()
		logrus.WithField("order", order).WithError(err).Error("unable to save order")
		return
	}

	metrics.OrdersReceived.Inc()
}
