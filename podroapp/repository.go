package podroapp

import (
	"context"
	"interview/adapter"
)

type DB struct {
	conn *adapter.SQLDB
}

func NewDB(conn *adapter.SQLDB) *DB {
	return &DB{
		conn: conn,
	}
}

func (db *DB) GetProvidersWeaklyReport(ctx context.Context) ([]Report, error) {
	query := `
        SELECT 
            provider_id, 
            AVG(EXTRACT(EPOCH FROM (delivered_at - picked_up_at)) / 3600) AS avg_delivery_time
        FROM 
            orders
        WHERE 
            created_at >= NOW() - INTERVAL '7 days'
        GROUP BY 
            provider_id
        ORDER BY 
            avg_delivery_time DESC
    `

	rows, err := db.conn.Conn.Query(query)
	if err != nil {
		return []Report{}, err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var r Report
		if err := rows.Scan(&r.Provider, &r.Average); err != nil {
			return []Report{}, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}

func (db *DB) UpdateOrders(ctx context.Context, orders []Order) error {
	for _, o := range orders {
		stmt := `
            UPDATE orders 
            SET provider_id = $1, customer_id = $2, customer_name = $3, customer_phone = $4, customer_address = $5, 
                recipient_phone = $6, recipient_name = $7, recipient_address = $8, status = $9, created_at = $10, 
                picked_up_at = $11, delivered_at = $12, updated_at = $13 
            WHERE id = $14
        `
		_, err := db.conn.Conn.Exec(stmt, o.ProviderID, o.CustomerID, o.CustomerName, o.CustomerPhone, o.CustomerAddress, o.RecipientPhone, o.RecipientName, o.RecipientAddress, o.Status, o.CreatedAt, o.PickedUpAt, o.DeliveredAt, o.UpdatedAt, o.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
func (db *DB) GetOrders(ctx context.Context) ([]Order, error) {
	// note : pending must be in condition to update to ProviderSeen so i removed !Pending
	rows, err := db.conn.Conn.Query("SELECT * FROM orders  WHERE status != 'Delivered'")
	if err != nil {
		return []Order{}, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.ProviderID, &o.CustomerID,
			&o.CustomerName, &o.CustomerPhone, &o.CustomerAddress, &o.RecipientPhone,
			&o.RecipientName, &o.RecipientAddress, &o.Status, &o.CreatedAt, &o.PickedUpAt,
			&o.DeliveredAt, &o.UpdatedAt); err != nil {
			return []Order{}, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
func (db *DB) GetProviders(ctx context.Context) ([]Provider, error) {
	rows, err := db.conn.Conn.Query("SELECT * FROM providers")
	if err != nil {
		return []Provider{}, err
	}
	defer rows.Close()

	var providers []Provider
	for rows.Next() {
		var p Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.API, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return []Provider{}, err
		}
		providers = append(providers, p)
	}
	return providers, nil
}
