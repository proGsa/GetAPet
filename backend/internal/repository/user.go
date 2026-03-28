package repository

import (
	"database/sql"
	"errors"

	"getapet-backend/internal/models"
	entities "getapet-backend/internal/repository/models"
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

	userDB := entities.UserFromDomain(*user)
	err := r.db.QueryRow(
		query,
		userDB.FIO,
		userDB.TelephoneNumber,
		userDB.City,
		userDB.UserLogin,
		userDB.UserPassword,
		userDB.Status,
		userDB.UserDescription,
	).Scan(
		&userDB.ID,
		&userDB.FIO,
		&userDB.TelephoneNumber,
		&userDB.City,
		&userDB.UserLogin,
		&userDB.UserPassword,
		&userDB.Status,
		&userDB.UserDescription,
	)
	if err != nil {
		return nil, err
	}

	userDomain := entities.UserToDomain(userDB)
	return &userDomain, nil
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

	userEntities := make([]entities.UserDB, 0)
	for rows.Next() {
		var userDB entities.UserDB
		if err := rows.Scan(
			&userDB.ID,
			&userDB.FIO,
			&userDB.TelephoneNumber,
			&userDB.City,
			&userDB.UserLogin,
			&userDB.UserPassword,
			&userDB.Status,
			&userDB.UserDescription,
		); err != nil {
			return nil, err
		}
		userEntities = append(userEntities, userDB)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities.UsersToDomain(userEntities), nil
}


func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	const query = `
		SELECT id, fio, telephone_number, city, user_login, user_password, status, user_description
		FROM users
		WHERE id = $1
	`

	var user entities.UserDB
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

	userdomain := entities.UserToDomain(user)
	return &userdomain, nil
}


func (r *UserRepository) GetByLogin(login string) (*models.User, error) {
	const query = `
		SELECT id, fio, telephone_number, city, user_login, user_password, status, user_description
		FROM users
		WHERE user_login = $1
	`

	var userDB entities.UserDB
	err := r.db.QueryRow(query, login).Scan(
		&userDB.ID,
		&userDB.FIO,
		&userDB.TelephoneNumber,
		&userDB.City,
		&userDB.UserLogin,
		&userDB.UserPassword,
		&userDB.Status,
		&userDB.UserDescription,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	userDomain := entities.UserToDomain(userDB)
	return &userDomain, nil
}

func (r *UserRepository) Update(id uuid.UUID, user *models.User) (*models.User, error) {
	const query = `
		UPDATE users
		SET fio = $1, telephone_number = $2, city = $3, user_login = $4, user_password = $5, status = $6, user_description = $7
		WHERE id = $8
		RETURNING id, fio, telephone_number, city, user_login, user_password, status, user_description
	`

	userDB := entities.UserFromDomain(*user)
	err := r.db.QueryRow(
		query,
		userDB.FIO,
		userDB.TelephoneNumber,
		userDB.City,
		userDB.UserLogin,
		userDB.UserPassword,
		userDB.Status,
		userDB.UserDescription,
		id,
	).Scan(
		&userDB.ID,
		&userDB.FIO,
		&userDB.TelephoneNumber,
		&userDB.City,
		&userDB.UserLogin,
		&userDB.UserPassword,
		&userDB.Status,
		&userDB.UserDescription,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	userDomain := entities.UserToDomain(userDB)
	return &userDomain, nil
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
