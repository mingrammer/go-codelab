package models

import "time"

type sensor struct {
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	GenTime time.Time `json:"gen_time"`
}

// Struct for Gyro Sensor
type GyroSensor struct {
	sensor
	AngleVelocityX float32 `json:"x_axis_angle_velocity"`
	AngleVelocityY float32 `json:"y_axis_angle_velocity"`
	AngleVelocityZ float32 `json:"z_axis_angle_velocity"`
}

// Struct for Accelerometer Sensor
type AccelSensor struct {
	sensor
	GravityAccX float32 `json:"x_axis_gravity_acceleration"`
	GravityAccY float32 `json:"y_axis_gravity_acceleration"`
	GravityAccZ float32 `json:"z_axis_grativy_acceleration"`
}

// Struct for Temperature Sensor
type TempSensor struct {
	sensor
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}
