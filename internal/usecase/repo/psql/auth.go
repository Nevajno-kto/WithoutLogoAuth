package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/nevajno-kto/without-logo-auth/internal/entity"
	"github.com/nevajno-kto/without-logo-auth/pkg/postgres"
)

type AuthRepo struct {
	*postgres.Postgres
}

func NewAuthRepo(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}

func (r *AuthRepo) InsertAuthCode(ctx context.Context, c entity.User, sign int, code int) error {

	sqlStatement := `
	INSERT INTO auth
	VALUES ( DEFAULT, $1, $2, $3, $4, $5 )
	`

	_, err := r.Pool.Exec(ctx,
		sqlStatement,
		c.Phone, c.Organization, code, time.Now().Unix(), sign)

	if err != nil {
		return fmt.Errorf("psql - auth - InsertAuthCode - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *AuthRepo) GetAuthCode(ctx context.Context, c entity.User, sign int) (int, int64, error) {

	sqlStatement := `
	SELECT code, request_time
	FROM auth 
	WHERE typeofsign = $1 AND phone = $2 AND organization = $3`

	rows, err := r.Pool.Query(ctx, sqlStatement, sign, c.Phone, c.Organization)

	if err != nil {
		return 0, 0, fmt.Errorf("psql - auth - GetAuthCode - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var code int = 0
	var request_time int64 = 0

	for rows.Next() {
		err = rows.Scan(&code, &request_time)
		if err != nil {
			return code, request_time, fmt.Errorf("psql - auth - GetAuthCode - rows.Scan: %w", err)
		}
	}

	return code, request_time, nil
}

func (r *AuthRepo) UpdateAuthCode(ctx context.Context, c entity.User, sign int, code int) error {

	sqlStatement := `
	UPDATE auth
	SET phone = $1, organization = $2, code = $3, request_time = $4
	WHERE typeofsign = $5 AND phone = $1 AND organization = $2`

	_, err := r.Pool.Exec(ctx,
		sqlStatement,
		c.Phone, c.Organization, code, time.Now().Unix(), sign)

	if err != nil {
		return fmt.Errorf("psql - auth - UpdateAuthCode - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *AuthRepo) DeleteAuthCode(ctx context.Context, c entity.User, sign int) error {
	sqlStatement := `
	DELETE FROM auth
	WHERE typeofsign = $1 AND phone = $2 AND organization = $3
	`

	_, err := r.Pool.Exec(ctx,
		sqlStatement,
		sign, c.Phone, c.Organization)

	if err != nil {
		return fmt.Errorf("psql - auth - DeleteAuthCode - r.Pool.Exec: %w", err)
	}

	return nil
}
