package driver

type CreateDriverRequest struct {
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	VehicleInfo     string `json:"vehicle_info" binding:"required"`
	CurrentLocation string `json:"current_location" binding:"required"`
}

type UpdateDriverProfileRequest struct {
	VehicleInfo     string `json:"vehicle_info" binding:"required"`
	CurrentLocation string `json:"current_location" binding:"required"`
}

type UpdateDriverRequest struct {
	Column string      `json:"column" binding:"required"`
	Value  interface{} `json:"value" binding:"required"`
}

func (r *CreateDriverRequest) ToDriver() *Driver {
	return &Driver{
		FullName:        r.FullName,
		Email:           r.Email,
		VehicleInfo:     r.VehicleInfo,
		CurrentLocation: r.CurrentLocation,
		Available:       true,
	}
}
