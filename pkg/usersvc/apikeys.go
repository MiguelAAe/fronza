package usersvc

import (
	"context"
	"time"
)

func UpsertApiKey(APIKey []byte, userID int64) error {
	apiKey := ApiKeys{
		APIKeyOwner:      userID,
		APIKey:           APIKey,
		LastModifiedTime: time.Now(),
	}
	_, err := db.NewInsert().
		Model(&apiKey).
		On("CONFLICT (api_key_owner) DO UPDATE").Set("api_key = ?, last_modified_time = ?", APIKey, apiKey.LastModifiedTime).
		Exec(context.TODO())

	return err
}

func GetAPIKey(userID int64) ([]byte, error) {
	var apiKey ApiKeys

	err := db.NewSelect().Model(&apiKey).Where("api_key_owner = ?", userID).Scan(context.TODO())
	if err != nil {
		return nil, err
	}

	return apiKey.APIKey, nil
}
