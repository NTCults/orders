CREATE TABLE IF NOT EXISTS orders_db.orders (
    order_uid               varchar(30) primary key,
    track_number            varchar(30),
    entry                   varchar(30),
    locale                  varchar(30),
    internal_signature      varchar(30),
    customer_id             varchar(30),
    delivery_service        varchar(30),
    shardkey                varchar(30),
    sm_id                   integer,
    date_created            timestamp,
    oof_shard               varchar(30)
);

CREATE TABLE IF NOT EXISTS orders_db.delivery (
    order_uid               varchar(30),
	name                    varchar(30),
	phone                   varchar(30),
	zip                     varchar(30),
	city                    varchar(30),
	address                 varchar(30),
	region                  varchar(30),
	email                   varchar(30),
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);

CREATE TABLE IF NOT EXISTS orders_db.payment (
    order_uid               varchar(30),
	transaction             varchar(30),
	request_id              varchar(30),
	currency                varchar(30),
	provider                varchar(30),
	amount                  integer,
	payment_dt              integer,
	bank                    varchar(30),
	delivery_cost           integer,
	goods_total             integer,
	custom_fee              integer,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);

CREATE TABLE IF NOT EXISTS orders_db.item (
    order_uid       varchar(30),
    chart_id        integer,
	track_number    varchar(30),
	price           integer,
	rid             varchar(30),
	name            varchar(30),
	sale            integer,
	size            varchar(30),
	total_price     integer,
	nm_id           integer,
	brand           varchar(30),
	status          varchar(30),
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);
