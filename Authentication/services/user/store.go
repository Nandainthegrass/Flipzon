package user

import (
	"database/sql"
	"fmt"

	"github.com/Nandainthegrass/Flipzon/Authentication/types"
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
	u := &types.User{}
	err := s.db.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.Phone,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			return nil, nil
		}
		fmt.Println(err)
		return nil, fmt.Errorf("error retrieving user")
	}

	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	u := &types.User{}
	err := s.db.QueryRow("SELECT * FROM users WHERE ID = ?", id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.Phone,
	)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Store) CreateUser(user types.User) error {
	stmt, err := s.db.Prepare(InsertUser)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.Name, user.Email, user.Password, user.Phone)
	if err != nil {
		return err
	}
	fmt.Printf("User created successfully")
	return nil
}
