package repository

import (
	"database/sql"
	"errors"

	"getapet-backend/internal/models"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	const query = `
		INSERT INTO users (fio, telephone_number, city, user_login, user_password, status, user_description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, fio, telephone_number, city, user_login, user_password, status, user_description
	`

	err := r.db.QueryRow(
		query,
		user.FIO,
		user.TelephoneNumber,
		user.City,
		user.UserLogin,
		user.UserPassword,
		user.Status,
		user.UserDescription,
	).Scan(
		&user.ID,
		&user.FIO,
		&user.TelephoneNumber,
		&user.City,
		&user.UserLogin,
		&user.UserPassword,
		&user.Status,
		&user.UserDescription,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	const query = `
		SELECT id, fio, telephone_number, city, user_login, user_password, status, user_description
		FROM users
		ORDER BY fio
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.FIO,
			&user.TelephoneNumber,
			&user.City,
			&user.UserLogin,
			&user.UserPassword,
			&user.Status,
			&user.UserDescription,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	const query = `
		SELECT id, fio, telephone_number, city, user_login, user_password, status, user_description
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FIO,
		&user.TelephoneNumber,
		&user.City,
		&user.UserLogin,
		&user.UserPassword,
		&user.Status,
		&user.UserDescription,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByLogin(login string) (*models.User, error) {
	const query = `
		SELECT id, fio, telephone_number, city, user_login, user_password, status, user_description
		FROM users
		WHERE user_login = $1
	`

	var user models.User
	err := r.db.QueryRow(query, login).Scan(
		&user.ID,
		&user.FIO,
		&user.TelephoneNumber,
		&user.City,
		&user.UserLogin,
		&user.UserPassword,
		&user.Status,
		&user.UserDescription,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(id uuid.UUID, user *models.User) (*models.User, error) {
	const query = `
		UPDATE users
		SET fio = $1, telephone_number = $2, city = $3, user_login = $4, user_password = $5, status = $6, user_description = $7
		WHERE id = $8
		RETURNING id, fio, telephone_number, city, user_login, user_password, status, user_description
	`

	err := r.db.QueryRow(
		query,
		user.FIO,
		user.TelephoneNumber,
		user.City,
		user.UserLogin,
		user.UserPassword,
		user.Status,
		user.UserDescription,
		id,
	).Scan(
		&user.ID,
		&user.FIO,
		&user.TelephoneNumber,
		&user.City,
		&user.UserLogin,
		&user.UserPassword,
		&user.Status,
		&user.UserDescription,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	const query = `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}

	return nil
}
