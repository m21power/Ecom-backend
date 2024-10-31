package product

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/m21power/Ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateProduct(p *types.Product) error {
	p.CreatedAt = time.Now().UTC()
	res, err := s.db.Exec("INSERT INTO products (name, description, image, price,quantity, createdAt) VALUES(?,?,?,?,?,?)", p.Name, p.Description, p.Image, p.Price, p.Quantity, p.CreatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	var products = make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)

	}
	return products, nil
}

func (s *Store) GetProductByIDs(productIDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(", ?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)
	// convert productids to interface
	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		pr, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *pr)
	}
	return products, nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	p := new(types.Product)
	err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, price = ?, image = ?, description = ?, quantity = ? WHERE id = ?", product.Name, product.Price, product.Image, product.Description, product.Quantity, product.ID)
	if err != nil {
		return err
	}
	return nil

}
