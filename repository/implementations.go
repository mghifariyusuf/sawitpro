package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/SawitProRecruitment/UserService/handler/models"
	"github.com/SawitProRecruitment/UserService/repository/entity"
)

func (r *Repository) CreateUser(ctx context.Context, req *models.RegisterUserRequest) (int64, error) {
	var lastInsertID int64

	err := r.Db.QueryRowContext(ctx,
		"INSERT INTO users (phone_number, full_name, password) VALUES ($1, $2, $3) RETURNING id",
		req.PhoneNumber, req.FullName, req.Password).
		Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func (r *Repository) GetUser(ctx context.Context, filter *entity.UserFilter) (*entity.UserData, error) {
	user := new(entity.UserData)

	var condition []string
	query := "SELECT id, phone_number, full_name, password FROM users"

	if filter.ID != nil {
		condition = append(condition, fmt.Sprintf(" WHERE id = %d", *filter.ID))
	}

	if filter.PhoneNumber != nil {
		condition = append(condition, fmt.Sprintf(" WHERE phone_number = '%s'", *filter.PhoneNumber))
	}

	query += strings.Join(condition, " AND ")

	err := r.Db.QueryRowContext(ctx, query).
		Scan(&user.ID, &user.PhoneNumber, &user.FullName, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return user, err
	}

	return user, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, userID int64, req *models.UpdateUserProfileRequest) error {
	if req.FullName == nil && req.PhoneNumber == nil {
		return nil
	}

	query := "UPDATE users SET"

	if req.FullName != nil {
		query += fmt.Sprintf(" full_name = '%s',", *req.FullName)
	}

	if req.PhoneNumber != nil {
		query += fmt.Sprintf(" phone_number = '%s',", *req.PhoneNumber)
	}

	query = query[:len(query)-1]

	query += fmt.Sprintf(" WHERE id = %d", userID)

	result, err := r.Db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}

func (r *Repository) IncLogin(ctx context.Context, userID int64) error {
	_, err := r.Db.ExecContext(ctx,
		"UPDATE users SET successful_login = successful_login + 1 WHERE id = $1",
		userID)

	return err
}
