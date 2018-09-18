package train

import (
	"database/sql"

	pb "github.com/GoingFast/gotrains/train/protobuf"
)

type (
	store struct {
		*sql.DB
	}

	trainStore interface {
		createTrain(*pb.Train) error
		createRoute(*pb.Route) error
		trainExists(string) (bool, error)
		getTrains() ([]*pb.Train, error)
	}
)

func newTrainStore(db *sql.DB) store {
	return store{db}
}

func (s store) trainExists(bn string) (bool, error) {
	var t pb.Train
	err := s.QueryRow("SELECT brandname FROM trains WHERE brandname = $1", bn).Scan(&t.Brandname)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s store) createTrain(t *pb.Train) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO trains (brandname, brandlogo, brandfeatures) VALUES ($1, $2, $3)")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(t.Brandname, t.Brandlogo, t.Brandfeatures)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s store) createRoute(t *pb.Route) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO routes (brandname, routes_from, routes_to, price, date, time) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(t.Brandname, t.From, t.To, t.Price, t.Date, t.Time)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s store) getTrains() ([]*pb.Train, error) {
	var trains []*pb.Train

	rows, err := s.Query("SELECT brandname, brandlogo, brandfeatures FROM trains")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var train pb.Train
		err := rows.Scan(&train.Brandname, &train.Brandlogo, &train.Brandfeatures)
		if err != nil {
			return nil, err
		}

		trains = append(trains, &train)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return trains, nil
}
