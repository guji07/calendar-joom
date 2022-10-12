package storage

import (
	"calendar/pkg/model"
	"context"

	"github.com/doug-martin/goqu/v9"
)

func (r *Repository) CreateUser(ctx context.Context, user model.User) (int, error) {
	var id int

	_, err := r.storage.Insert(UsersTable).Rows(user).Returning("id").Executor().ScanValContext(ctx, &id)

	return id, err
}

func (r *Repository) IsUserExist(ctx context.Context, userID int) (bool, error) {
	var id int

	exist, err := r.storage.Select("id").From(UsersTable).Where(goqu.Ex{"id": userID}).ScanValContext(ctx, &id)
	if err != nil || id == 0 {
		return false, err
	}
	return exist, err

}
