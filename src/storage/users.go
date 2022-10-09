package storage

import (
	"context"

	"cryptoColony/src/model"
)

func (r *Repository) CreateUser(ctx context.Context, user model.User) (int, error) {
	var id int

	_, err := r.storage.Insert(UsersTableName).Rows(user).Returning("id").Executor().ScanValContext(ctx, &id)

	return id, err
}
