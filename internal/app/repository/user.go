package repository

import (
	"Backeven/internal/app/ds"
	"fmt"
)

func (r *Repository) RegisterUser(input ds.ChangeUserDTO) error {
	var existingUser ds.User
	err := r.db.Where(`"Login" = ?`, input.Login).First(&existingUser).Error
	if err == nil {
		return fmt.Errorf("логин уже занят")
	}
	user := ds.User{
		Login:    input.Login,
		Password: input.Password,
	}
	return r.db.Create(&user).Error
}

func (r *Repository) LoginUser(login, password string) (*ds.UserDTO, error) {
	var user ds.User
	err := r.db.Where(`"Login" = ? AND "Password" = ?`, login, password).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("неверный логин или пароль")
	}
	return &ds.UserDTO{
		UserId: user.UserId,
		Login:  user.Login,
	}, nil
}
func (r *Repository) UpdateUser(userID int, userUpdate ds.ChangeUserDTO) (*ds.UserDTO, error) {
	var user ds.User
	err := r.db.Where(`"UserID" = ?`, userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	if userUpdate.Login != "" {
		var count int64
		r.db.Where(&ds.User{}).Where(`"Login" = ? AND "UserID" != ?`, userUpdate.Login, user).Count(&count)
		if count > 0 {
			return nil, fmt.Errorf("логин уже занят")
		}
		user.Login = userUpdate.Login
	}
	if userUpdate.Password != "" {
		user.Password = userUpdate.Password
	}
	err = r.db.Save(&user).Error
	if err != nil {
		return nil, err
	}
	return &ds.UserDTO{
		UserId: user.UserId,
		Login:  user.Login,
	}, nil
}

func (r *Repository) LogoutUser(userid int) error {
	return nil
}
func (r *Repository) GetUserByID(userid int) (*ds.UserDTO, error) {
	var user ds.User
	err := r.db.Where(`"UserID" = ?`, userid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &ds.UserDTO{
		UserId: user.UserId,
		Login:  user.Login,
	}, nil
}
