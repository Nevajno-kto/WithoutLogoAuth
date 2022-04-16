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

func (r *AuthRepo) InsertSignUpCode(ctx context.Context, c entity.User, code int) error {

	sqlStatement := `
	INSERT INTO requestsignup
	VALUES ( DEFAULT, $1, $2, $3, $4 )
	`

	_, err := r.Pool.Exec(ctx,
		sqlStatement,
		c.Phone, c.Organization, code, time.Now().Unix())

	if err != nil {
		return fmt.Errorf("psql - auth - InsertSignUpCode - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *AuthRepo) GetSignUpCode(ctx context.Context, c entity.User) (int, int64, error) {

	sqlStatement := `
	SELECT code, request_time
	FROM requestsignup 
	WHERE phone = $1 AND organization = $2`

	rows, err := r.Pool.Query(ctx, sqlStatement, c.Phone, c.Organization)

	if err != nil {
		return 0, 0, fmt.Errorf("psql - auth - GetSignUpCode - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var code int = 0
	var request_time int64 = 0

	for rows.Next() {
		err = rows.Scan(&code, &request_time)
		if err != nil {
			return code, request_time, fmt.Errorf("psql - auth - GetSignUpCode - rows.Scan: %w", err)
		}
	}

	return code, request_time, nil
}

func (r *AuthRepo) UpdateSignUpCode(ctx context.Context, c entity.User, code int) error {

	sqlStatement := `
	UPDATE requestsignup
	SET phone = $1, organization = $2, code = $3, request_time = $4
	WHERE phone = $1 AND organization = $2`

	_, err := r.Pool.Exec(ctx,
		sqlStatement,
		c.Phone, c.Organization, code, time.Now().Unix())

	if err != nil {
		return fmt.Errorf("psql - auth - UpdateSignUpCode - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *AuthRepo) DeleteSignUpCode(ctx context.Context, c entity.User) error {
	sqlStatement := `
	DELETE FROM requestsignup
	WHERE phone = $1 AND password = $2 AND organization = $3
	`

	_, err := r.Pool.Exec(ctx,
		sqlStatement,
		c.Phone, c.Password, c.Organization)

	if err != nil {
		return fmt.Errorf("psql - auth - DeleteSignUpCode - r.Pool.Exec: %w", err)
	}

	return nil
}
