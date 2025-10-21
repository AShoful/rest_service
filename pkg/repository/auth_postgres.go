package repository

import "rest/models"

type AuthPostgres struct {
	db any
}

func NewAuthPostgres(db any) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {

	return 0, nil
}

func (r *AuthPostgres) GetUser(username, password string) {

}
