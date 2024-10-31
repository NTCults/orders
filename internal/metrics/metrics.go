package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrdersReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_received",
		Help: "The total number of orders received",
	})
)

var (
	OrdersReceiveErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_receive_errors",
		Help: "Errors on order receive",
	})
)
