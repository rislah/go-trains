package user

import (
	"database/sql"

	pb "github.com/GoingFast/gotrains/user/protobuf"
)

type (
	userStore interface {
		findUserByEmail(string) (pb.User, bool, error)
		findUserByUsername(string) (pb.User, bool, error)
		createUser(*pb.User, string) error
	}

	store struct {
		*sql.DB
	}
)

func newUserStore(db *sql.DB) store {
	return store{db}
}

func (s store) findUserByEmail(email string) (pb.User, bool, error) {
	var u pb.User
	err := s.QueryRow("SELECT password FROM users WHERE email = $1", email).Scan(&u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, false, nil
		}
		return u, false, err
	}
	return u, true, nil
}

func (s store) findUserByUsername(uname string) (pb.User, bool, error) {
	var u pb.User
	err := s.QueryRow("SELECT username, password, verified, role FROM users WHERE username = $1", uname).Scan(&u.Username, &u.Password, &u.Verified, &u.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, false, nil
		}
		return u, false, err
	}
	return u, true, nil
}

func (s store) createUser(u *pb.User, vid string) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO users (username, firstname, lastname, email, password, verificationid) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(u.Username, u.Firstname, u.Lastname, u.Email, u.Password, vid)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
