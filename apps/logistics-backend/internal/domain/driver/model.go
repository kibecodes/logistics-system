package driver

import (
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type Driver struct {
	ID              uuid.UUID      `db:"id" json:"id"`
	FullName        string         `db:"full_name" json:"full_name"`
	Email           string         `db:"email" json:"email"`
	VehicleInfo     string         `db:"vehicle_info" json:"vehicle_info"`
	CurrentLocation postgis.PointS `db:"current_location" json:"current_location"`
	Available       bool           `db:"available" json:"available"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
}
