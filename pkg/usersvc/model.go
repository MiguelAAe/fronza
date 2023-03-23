package usersvc

import (
	"time"

	"github.com/uptrace/bun"
)

func Init(dbconn *bun.DB) {
	db = dbconn
}

var db *bun.DB

type User struct {
	ID             int64     `bun:"id,pk,autoincrement"`
	Name           string    `bun:"name"`
	Email          string    `bun:"email,pk"`
	PhoneNumber    string    `bun:"phone_number"`
	HashedPassword []byte    `bun:"hashed_password"`
	Role           string    `bun:"user_role"`
	CreateTime     time.Time `bun:"create_time"`
}

type ApiKeys struct {
	APIKeyOwner      int64     `bun:"api_key_owner,pk"`
	APIKey           []byte    `bun:"api_key"`
	LastModifiedTime time.Time `bun:"last_modified_time"`
}

type DriverState struct {
	CourierID int64 `bun:"courier_id"`
	Active    bool  `bun:"active"`
}

type DriverLocation struct {
	CourierID   int64  `bun:"courier_id"`
	Coordinates string `bun:"coordinates"`
}
