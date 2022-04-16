package psql

import (
	"context"
	"fmt"

	"github.com/nevajno-kto/without-logo-auth/internal/entity"
	"github.com/nevajno-kto/without-logo-auth/pkg/postgres"
)

type UsersRepo struct {
	*postgres.Postgres
}

func NewClientsRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}

func (r *UsersRepo) GetUser(ctx context.Context, c entity.User) (entity.User, error) {

	sqlStatement := fmt.Sprintf("SELECT * FROM users%s WHERE phone = $1", c.Organization)

	rows, err := r.Pool.Query(ctx, sqlStatement, c.Phone)
	if err != nil {
		return entity.User{}, fmt.Errorf("psql - user - GetUser - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	e := entity.User{}
	for rows.Next() {
		err = rows.Scan(&e.Id, &e.Name, &e.Phone, &e.Password)
		if err != nil {
			return entity.User{}, fmt.Errorf("psql - user - GetUser - rows.Scan: %w", err)
		}
	}

	return e, nil
}

func (r *UsersRepo) InsertUser(ctx context.Context, c entity.User) error {

	sqlStatement := fmt.Sprintf("INSERT INTO users%s VALUES ( DEFAULT, $1, $2, $3 )", c.Organization)

	_, err := r.Pool.Exec(ctx,
		sqlStatement,
		c.Name, c.Phone, c.Password)

	if err != nil {
		return fmt.Errorf("psql - user - InsertUser - r.Pool.Exec: %w", err)
	}

	return nil
}
