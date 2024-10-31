package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"orders/internal/query"
	"orders/models"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var natsConn *nats.Conn

func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	natsURL := os.Getenv("NATS_URL")

	var err error
	natsConn, err = nats.Connect(natsURL)
	if err != nil {
		panic(err)
	}

	logrus.WithField("port", port).Info("Starting server.")
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		natsConn.Close()
		logrus.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var order models.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := natsConn.Publish(query.NatsSubject+order.OrderUID, buf.Bytes()); err != nil {
		logrus.Error(err)
	}
}

//JSON EXAMPLE
/*
{
    "order_uid": "551",
    "track_number": "1",
    "entry": "",
    "delivery": {
        "order_uid": "1",
        "phone": "12341",
        "zip": "1232",
        "city": "Bubuga",
        "address": "123",
        "region": "KAfdf",
        "email": "dfDF@FDSFDF.COM"
    },
    "payment": {
        "transaction": "1",
        "request_id": "",
        "currency": "",
        "provider": "",
        "amount": 0,
        "payment_dt": 0,
        "bank": "",
        "delivery_cost": 0,
        "goods_total": 0,
        "custom_fee": 0
    },
    "items": [
        {
            "chart_id": 1,
            "track_number": "",
            "price": 0,
            "rid": "",
            "name": "",
            "sale": 0,
            "size": "",
            "total_price": 0,
            "nm_id": 0,
            "brand": "",
            "status": ""
        },
        {
            "chart_id": 2,
            "track_number": "asfaf",
            "price": 123142413,
            "rid": "dasd",
            "name": "dsdad",
            "sale": 213123,
            "size": "asda",
            "total_price": 0,
            "nm_id": 0,
            "brand": "",
            "status": ""
        }
    ],
    "locale": "",
    "internal_signature": "",
    "customer_id": "",
    "delivery_service": "",
    "sharedkey": "",
    "sm_id": 0,
    "date_created": "0001-01-01T00:00:00Z",
    "oof_shard": ""
}
*/
