package usersvc

import (
	"context"
	"fmt"
)

func GetDriverStates() ([]DriverState, error) {
	driverStates := &[]DriverState{}

	err := db.NewSelect().Model(driverStates).Scan(context.TODO())

	if driverStates == nil {
		return nil, nil
	}
	return *driverStates, err
}

func GetDriverLocation(driverID int64) (string, error) {
	var resp DriverLocation
	err := db.NewSelect().Model(&resp).Where("courier_id = ?", driverID).Scan(context.TODO())
	if err != nil {
		return "", err
	}

	if resp.Coordinates == "" {
		return "", fmt.Errorf("empty location")
	}

	return resp.Coordinates, nil
}

func SaveDriverLocation(driverID int64, coordinates string) error {
	var resp DriverLocation

	resp.Coordinates = coordinates
	resp.CourierID = driverID
	_, err := db.NewInsert().
		Model(&resp).
		On("CONFLICT (courier_id) DO UPDATE").Set("coordinates = ?", coordinates).
		Exec(context.TODO())

	return err
}

func SaveDriverStatus(driverID int64, active bool) error {
	var resp DriverState

	resp.Active = active
	resp.CourierID = driverID
	_, err := db.NewInsert().
		Model(&resp).
		On("CONFLICT (courier_id) DO UPDATE").Set("active = ?", active).
		Exec(context.TODO())

	return err
}
