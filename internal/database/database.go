package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	"WB_GO_L0/internal/entity"

	_ "github.com/jackc/pgx/v5/stdlib"    // Import the pgx driver
	_ "github.com/joho/godotenv/autoload" // Load the .env file
	"github.com/redis/go-redis/v9"
)

// Service represents a service that interacts with a database.
type Service interface {
	Close() error
	GetOrder(id string) (entity.Order, error)
	GetOrderByUID(uid string) (entity.Order, error)
	GetOrdersPlain() ([]string, error)
	SaveOrderPlain(order string) error
}

type service struct {
	db    *sql.DB
	cashe *redis.Client
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	radress    = os.Getenv("REDIS_ADDRESS")
	rport      = os.Getenv("REDIS_PORT")
	dbInstance *service
	ctx        = context.Background()
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	cache := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", radress, rport),
	})
	dbInstance = &service{
		db:    db,
		cashe: cache,
	}
	err = dbInstance.restoreCache()
	if err != nil {
		log.Fatal(err)
	}
	return dbInstance
}

// Close closes the database connection.
func (s *service) Close() error {
	slog.Info("Disconnected from database", "db name", database)
	return s.db.Close()
}

// SaveOrderPlain saves an order to the database as a json.
func (s *service) SaveOrderPlain(order string) error {
	var id string
	var uid string
	err := s.db.QueryRow("INSERT INTO orders_plain (order_json) VALUES ($1) RETURNING id, order_json->'order_uid'", order).Scan(&id, &uid)
	if err != nil {
		return err
	}
	// save to cashe
	err = s.cashe.Set(ctx, id, order, 0).Err()
	if err != nil {
		return err
	}
	err = s.cashe.Set(ctx, uid, id, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetOrder(id string) (entity.Order, error) {
	var order entity.Order
	orderStr, err := s.cashe.Get(ctx, id).Result()
	if err != nil {
		return order, err
	}
	err = json.Unmarshal([]byte(orderStr), &order)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (s *service) GetOrderByUID(uid string) (entity.Order, error) {
	var id string
	id, err := s.cashe.Get(ctx, uid).Result()
	if err != nil {
		return entity.Order{}, err
	}
	return s.GetOrder(id)
}

// GetOrdersPlain returns all orders as slice of strings from the database
func (s *service) GetOrdersPlain() ([]string, error) {
	var orders []string
	rows, err := s.db.Query("SELECT order_json FROM orders_plain")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order string
		err := rows.Scan(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// RestoreCache restores the cache from the database
func (s *service) restoreCache() error {
	slog.Debug("Restoring cache from database")
	rows, err := s.db.Query("SELECT id, order_json->'order_uid', order_json FROM orders_plain")
	if err != nil {
		return err
	}

	for rows.Next() {
		var id string
		var uid string
		var order string
		err := rows.Scan(&id, &uid, &order)
		if err != nil {
			return err
		}
		err = s.cashe.Set(ctx, id, order, 0).Err()
		if err != nil {
			return err
		}
		err = s.cashe.Set(ctx, uid, id, 0).Err()
		if err != nil {
			return err
		}
	}
	slog.Debug("Cache restored from database")
	return nil
}
