package entities

import ("github.com/google/uuid"
"getapet-backend/internal/models"
)


type UserDB struct {
	ID              uuid.UUID `db:"id"`
	FIO             string    `db:"fio"`
	TelephoneNumber string    `db:"telephone_number"`
	City            string    `db:"city"`
	UserLogin       string    `db:"user_login"`
	UserPassword    string    `db:"user_password"`
	Status          string    `db:"status"`
	UserDescription string    `db:"user_description"`
}


func UserToDomain(entity UserDB) models.User {
	return models.User{
		ID:              entity.ID,
		FIO:             entity.FIO,
		TelephoneNumber: entity.TelephoneNumber,
		City:            entity.City,
		UserLogin:       entity.UserLogin,
		UserPassword:    entity.UserPassword,
		Status:          entity.Status,
		UserDescription: entity.UserDescription,
	}
}

func UserFromDomain(user models.User) UserDB {
	return UserDB{
		ID:              user.ID,
		FIO:             user.FIO,
		TelephoneNumber: user.TelephoneNumber,
		City:            user.City,
		UserLogin:       user.UserLogin,
		UserPassword:    user.UserPassword,
		Status:          user.Status,
		UserDescription: user.UserDescription,
	}
}

func UsersToDomain(entities []UserDB) []models.User {
	users := make([]models.User, len(entities))
	for i, entity := range entities {
		users[i] = UserToDomain(entity)
	}
	return users
}