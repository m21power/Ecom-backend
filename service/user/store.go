package user

import (
	"database/sql"
	"fmt"

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
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = ScanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil

}

func (s *Store) GetUserByID(id int) (*types.User, error) {

	rows, err := s.db.Query("SELECT * FROM user WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = ScanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil

}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO user(firstName,lastName,email,password) VALUES(?,?,?,?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
func ScanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	u := new(types.User)
	err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err

	}
	return u, nil
}
