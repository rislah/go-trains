package user

import (
	"database/sql"

	pb "github.com/GoingFast/gotrains/user/protobuf"
)

type (
	userStore interface {
		findUserByEmail(string) (pb.User, bool, error)
		findUserByUsername(string) (pb.User, bool, error)
		getUsers() ([]*pb.User, error)
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

func (s store) getUsers() ([]*pb.User, error) {
	var users []*pb.User

	rows, err := s.Query("SELECT username, firstname, lastname, email, role, uuid FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user pb.User
		err = rows.Scan(&user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Role, &user.Uuid)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}
