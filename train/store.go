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
		trainExists(string) (bool, error)
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
