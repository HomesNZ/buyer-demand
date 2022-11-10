package address

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

type Repo interface {
	IsClaimedByUserID(ctx context.Context, propertyID string, userID int) (bool, error)
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

const isClaimedByUserIDQuery = `
	SELECT true 
	FROM api.property_claims pc
	WHERE pc.property_id = $1
	and pc.user_id = $2
	AND pc.unclaimed_on ISNULL
	LIMIT 1;
`

func (r *repo) IsClaimedByUserID(ctx context.Context, propertyID string, userID int) (bool, error) {
	rows, err := r.db.Query(ctx, isClaimedByUserIDQuery, propertyID, userID)
	if err != nil {
		return false, err
	}
	var claimedByUser null.Bool
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&claimedByUser)
		if err != nil {
			return false, err
		}
	}
	return claimedByUser.ValueOrZero(), nil
}
