package controllers

import (
	"net/http"

	"github.com/YL-Tan/GoHomeAi/internal/db"
	"github.com/YL-Tan/GoHomeAi/internal/httpresponse"
)

type DeviceController struct {
	Queries *db.Queries
}

func NewDeviceController(q *db.Queries) *DeviceController {
	return &DeviceController{Queries: q}
}

func (c *DeviceController) GetDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := c.Queries.GetDevices(r.Context())
	if err != nil {
		httpresponse.SendJSON(w, http.StatusInternalServerError, false, "Failed to fetch devices", nil, err)
		return
	}
	httpresponse.SendJSON(w, http.StatusOK, true, "Devices fetched successfully", devices, nil)
}
