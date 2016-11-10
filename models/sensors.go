package models

import (
	"time"
	"fmt"
	"strings"
)

// Sensor has common fields for any sensors
type Sensor struct {
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	GenTime time.Time `json:"gen_time"`
}

// GyroSensor produces x-y-z axes angle velocity values
type GyroSensor struct {
	Sensor
	AngleVelocityX float32 `json:"x_axis_angle_velocity"`
	AngleVelocityY float32 `json:"y_axis_angle_velocity"`
	AngleVelocityZ float32 `json:"z_axis_angle_velocity"`
}

func (s GyroSensor) String() string {
	var result []string

	st := fmt.Sprintf("Measured on %s", s.GenTime)
	result = append(result, st)
	
	st = fmt.Sprintf("Angle Velocity of X-axis : %f", s.AngleVelocityX)
	result = append(result, st)

	st = fmt.Sprintf("Angle Velocity of Y-axis : %f", s.AngleVelocityY)
	result = append(result, st)

	st = fmt.Sprintf("Anglue Velocity of Z-axis : %f\n", s.AngleVelocityZ)
	result = append(result, st)

	return strings.Join(result, "\n")
}

// AccelSensor produces x-y-z axes gravity acceleration values
type AccelSensor struct {
	Sensor
	GravityAccX float32 `json:"x_axis_gravity_acceleration"`
	GravityAccY float32 `json:"y_axis_gravity_acceleration"`
	GravityAccZ float32 `json:"z_axis_grativy_acceleration"`
}

func (s AccelSensor) String() string {
	var result []string

	st := fmt.Sprintf("Measured on %s", s.GenTime)
	result = append(result, st)
	
	st = fmt.Sprintf("Gravitational Velocity of X-axis : %f", s.GravityAccX)
	result = append(result, st)

	st = fmt.Sprintf("Gravitational Velocity of Y-axis : %f", s.GravityAccY)
	result = append(result, st)

	st = fmt.Sprintf("Gravitational Velocity of Z-axis : %f\n", s.GravityAccZ)
	result = append(result, st)

	return strings.Join(result, "\n")
}

// TempSensor produces temperature and humidity values
type TempSensor struct {
	Sensor
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}

func (s TempSensor) String() string {
	var result []string

	st := fmt.Sprintf("Measured on %s", s.GenTime)
	result = append(result, st)
	
	st = fmt.Sprintf("Temperature : %f", s.Temperature)
	result = append(result, st)

	st = fmt.Sprintf("Humidity : %f\n", s.Humidity)
	result = append(result, st)

	return strings.Join(result, "\n")
}

