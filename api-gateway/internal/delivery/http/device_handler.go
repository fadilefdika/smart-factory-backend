package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/smartfactory/api-gateway/internal/domain"
)

type DeviceHandler struct {
	deviceUsecase domain.DeviceUsecase
}

func NewDeviceHandler(r *gin.Engine, us domain.DeviceUsecase) {
	handler := &DeviceHandler{
		deviceUsecase: us,
	}

	api := r.Group("/api/v1")
	{
		api.GET("/devices", handler.ListDevices)
		api.POST("/devices", handler.RegisterDevice)
		api.GET("/devices/:id/status", handler.GetDeviceStatus)
	}
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
	devices, err := h.deviceUsecase.ListDevices("")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *DeviceHandler) RegisterDevice(c *gin.Context) {
	var device domain.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.deviceUsecase.RegisterDevice(&device); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fake MQTT credentials returned according to PRD
	c.JSON(http.StatusCreated, gin.H{
		"message": "Device registered successfully",
		"device":  device,
		"mqtt_credentials": gin.H{
			"username": device.ID.String(),
			"password": "auto-generated-secret-key",
		},
	})
}

func (h *DeviceHandler) GetDeviceStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	device, err := h.deviceUsecase.GetDeviceStatus(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}
