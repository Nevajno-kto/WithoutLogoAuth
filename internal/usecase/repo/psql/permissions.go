package psql

import (
	"context"
	"fmt"

	"github.com/nevajno-kto/without-logo-auth/internal/entity"
	"github.com/nevajno-kto/without-logo-auth/pkg/postgres"
)

type PemissionsRepo struct {
	*postgres.Postgres
}

func NewPemissionsRepo(pg *postgres.Postgres) *PemissionsRepo {
	return &PemissionsRepo{pg}
}

func (r *PemissionsRepo) InsertPermissionForUser(ctx context.Context, u entity.User, p entity.Permission) error {
	sqlStatement := fmt.Sprintf(`
	INSERT INTO users_permissions%s
	VALUES ( DEFAULT, $1, $2 )`, u.Organization)

	_, err := r.Pool.Exec(ctx,
		sqlStatement,
		u.Id, p.Id)

	if err != nil {
		return fmt.Errorf("psql - auth - InsertPermissionForUser - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *PemissionsRepo) GetClientPermission(ctx context.Context, u entity.User) (entity.Permission, error) {
	var permission entity.Permission

	sqlStatement := fmt.Sprintf("SELECT * FROM permissions%s WHERE id = 1", u.Organization)

	rows, err := r.Pool.Query(ctx, sqlStatement)

	if err != nil {
		return permission, fmt.Errorf("psql - auth - GetClientPermission - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&permission.Id)
		if err != nil {
			return permission, fmt.Errorf("psql - auth - GetClientPermission - rows.Scan: %w", err)
		}
	}

	return permission, nil
}

func (r *PemissionsRepo) GetAdminPermission(ctx context.Context, u entity.User) (entity.Permission, error) {
	var permission entity.Permission

	sqlStatement := fmt.Sprintf("SELECT * FROM permissions%s WHERE id = 2", u.Organization)

	rows, err := r.Pool.Query(ctx, sqlStatement)

	if err != nil {
		return permission, fmt.Errorf("psql - auth - GetAdminPermission - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&permission.Id)
		if err != nil {
			return permission, fmt.Errorf("psql - auth - GetAdminPermission - rows.Scan: %w", err)
		}
	}

	return permission, nil
}
