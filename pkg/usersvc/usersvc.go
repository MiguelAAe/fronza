package usersvc

import (
	"context"
	"time"
)

const (
	Admin    = "admin"
	Customer = "customer"
	Courier  = "courier"
)

func CreateUser(user User) error {
	user.CreateTime = time.Now().UTC()
	_, err := db.NewInsert().Model(&user).Exec(context.TODO())
	return err
}

func GetUser(email string) (User, error) {
	var user User
	err := db.NewSelect().Model(&user).Where("email = ?", email).Scan(context.TODO())
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserByID(userID int64) (User, error) {
	var user User
	err := db.NewSelect().Model(&user).Where("id = ?", userID).Scan(context.TODO())
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func UpdatePassword(email string, hashedPassword []byte) error {
	_, err := db.NewUpdate().Model(&User{}).Where("email = ?", email).Set("hashed_password = ?", hashedPassword).Exec(context.TODO())
	return err
}
