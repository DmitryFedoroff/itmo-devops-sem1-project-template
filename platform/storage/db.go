package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"project_sem/platform/config"
	"strings"
)

type Repository interface {
	InsertProducts(products []Product) error
	GetAllProductsFiltered(start, end string, min, max string) ([]Product, error)
	GetAllProducts() ([]Product, error)
	GetStats() (int, int, float64, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(cfg config.DatabaseSettings) (Repository, error) {
	log.Println("connecting to the database...")
	const sslModeDisable = "disable"

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		sslModeDisable,
	)

	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err = database.Ping(); err != nil {
		return nil, err
	}
	log.Printf("successfully connected to database '%s'\n", cfg.Name)
	return &repository{db: database}, nil
}

func (r *repository) InsertProducts(products []Product) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertProductsQuery)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, product := range products {
		_, err := stmt.Exec(product.ID, product.Name, product.Category, product.Price, product.CreateDate)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *repository) GetAllProducts() ([]Product, error) {
	rows, err := r.db.Query(getAllProductsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0, 10)

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.CreateDate)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) GetAllProductsFiltered(start, end string, min, max string) ([]Product, error) {
	var (
		conditions []string
		args       []interface{}
	)
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(baseFilteredQuery)

	argIndex := 1

	if start != "" && end != "" {
		conditions = append(conditions, fmt.Sprintf("create_date BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, start, end+fullDayEndSuffix)
		argIndex += 2
	}
	if min != "" && max != "" {
		conditions = append(conditions, fmt.Sprintf("price BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, min, max)
		argIndex += 2
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		for i, cond := range conditions {
			if i > 0 {
				queryBuilder.WriteString(" AND ")
			}
			queryBuilder.WriteString(cond)
		}
	}

	queryBuilder.WriteString(" ORDER BY id")

	rows, err := r.db.Query(queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0, 10)

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.CreateDate)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) GetStats() (int, int, float64, error) {
	var totalItems, totalCategories int
	var totalPrice float64
	err := r.db.QueryRow(getStatsQuery).Scan(&totalItems, &totalCategories, &totalPrice)
	if err != nil {
		return 0, 0, 0, err
	}
	return totalItems, totalCategories, totalPrice, nil
}
