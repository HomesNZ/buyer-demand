package address

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repo interface {
	AllSuburbIDs(ctx context.Context) ([]int, error)
}

func New(db *pgxpool.Pool) (Repo, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &repo{db}, nil
}

type repo struct {
	db *pgxpool.Pool
}

const allSuburbIDsQuery = `
	SELECT id FROM homes_external_data.nz_localities
	WHERE type = 'SUBURB' OR type = 'LOCALITY';
`

func (r *repo) AllSuburbIDs(ctx context.Context) ([]int, error) {
	var result []int
	rows, err := r.db.Query(ctx, allSuburbIDsQuery)
	if err != nil {
		return nil, errors.Wrap(err, "db.Query")
	}
	defer rows.Close()
	for rows.Next() {
		var suburbID int
		err = rows.Scan(&suburbID)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}

		result = append(result, suburbID)
	}
	return result, nil
}
