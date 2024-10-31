package store

import (
	"database/sql"
	"orders/models"

	"github.com/sirupsen/logrus"
)

func (s *PGStore) writeOrderTX(order *models.Order) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if err := writeOrderToDB(tx, order); err != nil {
		return err
	}

	if err := writeDeliveryToDB(tx, order); err != nil {
		return err
	}

	if err := writePaymentToDB(tx, order); err != nil {
		return err
	}

	if err := writeItemsToDB(tx, order); err != nil {
		return err
	}

	return tx.Commit()
}

func writeOrderToDB(tx *sql.Tx, order *models.Order) error {
	_, err := tx.Exec(`
		INSERT INTO orders (
			order_uid,
    		track_number,
    		entry,
    		locale,
    		internal_signature,
    		customer_id,
    		delivery_service,
    		shardkey,
    		sm_id,
    		date_created,
    		oof_shard) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard)
	return err
}

func writeDeliveryToDB(tx *sql.Tx, order *models.Order) error {
	_, err := tx.Exec(`
		INSERT INTO delivery (
			order_uid,
			name,
			phone,
			zip,
			city,
			address,
			region,
			email)
			values ($1, $2, $3, $4, $5, $6, $7, $8)
		`,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email)
	return err
}

func writePaymentToDB(tx *sql.Tx, order *models.Order) error {
	_, err := tx.Exec(`
		INSERT INTO payment (
			order_uid,
			transaction,
			request_id,
			currency,
			provider,
			amount,
			payment_dt,
			bank,
			delivery_cost,
			goods_total,
			custom_fee)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee)
	return err
}

func writeItemsToDB(tx *sql.Tx, order *models.Order) error {
	for _, item := range order.Items {
		_, err := tx.Exec(`
		INSERT INTO item (
			order_uid,
    		chart_id,
			track_number,
			price,
			rid,
			name,
			sale,
			size,
			total_price,
			nm_id,
			brand,
			status)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`,
			order.OrderUID,
			item.ChartID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status)

		return err
	}

	return nil
}

func (s *PGStore) getAllOrdersFromDB() ([]*models.Order, error) {
	rows, err := s.db.Query(`
		SELECT 
			orders.order_uid,
			orders.track_number,
			orders.entry,
			payment.transaction,
			delivery.name
		FROM orders
		JOIN payment ON orders.order_uid = payment.order_uid
		JOIN delivery ON orders.order_uid = delivery.order_uid
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*models.Order{}
	for rows.Next() {
		o := models.Order{
			Delivery: &models.Delivery{},
			Payment:  &models.Payment{},
		}
		err := rows.Scan(&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Payment.Transaction, &o.Delivery.Name)
		if err != nil {
			return nil, err
		}

		items, err := s.getOrderItems(o.OrderUID)
		// хз как обрабатывать
		if err != nil {
			logrus.WithField("runContext", "store.getAllOdersFromDB").Error(err)
		}

		o.Items = items

		orders = append(orders, &o)
	}
	return orders, nil
}

func (s *PGStore) getOrderFromDB(orderUID string) (*models.Order, error) {
	row := s.db.QueryRow(`
		SELECT 
			orders.order_uid,
			orders.track_number,
			orders.entry,
			orders.locale,
			orders.internal_signature,
			orders.customer_id,
			orders.delivery_service,
			orders.shardkey,
			orders.sm_id,
			orders.date_created,
			orders.oof_shard,

			delivery.name,
			delivery.phone, 
			delivery.zip, 
			delivery.city, 
			delivery.address, 
			delivery.region, 
			delivery.email,
			
			payment.transaction,
			payment.request_id,
			payment.currency,
			payment.provider,
			payment.amount,
			payment.payment_dt,
			payment.bank,
			payment.delivery_cost,
			payment.goods_total,
			payment.custom_fee

		FROM orders
			JOIN delivery ON orders.order_uid = delivery.order_uid
			JOIN payment ON orders.order_uid = payment.order_uid
		WHERE
			orders.order_uid = $1
	`, orderUID)

	o := models.Order{
		Delivery: &models.Delivery{},
		Payment:  &models.Payment{},
	}

	err := row.Scan(
		&o.OrderUID, &o.TrackNumber, &o.Entry,
		&o.Locale, &o.InternalSignature, &o.CustomerID,
		&o.DeliveryService, &o.Shardkey, &o.SmID,
		&o.DateCreated, &o.OofShard,

		&o.Delivery.Name, &o.Delivery.Phone,
		&o.Delivery.Zip, &o.Delivery.City, &o.Delivery.Address,
		&o.Delivery.Region, &o.Delivery.Email,

		&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency,
		&o.Payment.Provider, &o.Payment.Amount, &o.Payment.PaymentDT,
		&o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal,
		&o.Payment.CustomFee,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	items, err := s.getOrderItems(orderUID)
	if err != nil {
		return nil, err
	}

	o.Items = items

	return &o, nil
}

func (s *PGStore) getOrderItems(orderUID string) ([]*models.Item, error) {
	rows, err := s.db.Query(`
		SELECT
			chart_id,
			track_number,
			price,
			rid,
			name,
			sale,
			size,
			total_price,
			nm_id,
			brand,
			status
		FROM item
		WHERE order_uid = $1
	`, orderUID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*models.Item{}
	for rows.Next() {
		i := models.Item{}

		err := rows.Scan(&i.ChartID, &i.TrackNumber, &i.Price, &i.RID, &i.Name, &i.Sale, &i.Size,
			&i.TotalPrice, &i.NmID, &i.Brand, &i.Status)
		if err != nil {
			panic(err)
		}

		items = append(items, &i)
	}
	return items, nil
}
