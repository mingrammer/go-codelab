package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mingrammer/go-codelab/faker"
)

// Sensor is common interface for any sensors
type Sensor interface {
	SendingOutputString() string
	ReceivingOutputString() string
	GenerateSensorData(epsilon float64) Sensor
}

// SensorInfo has common fields for any sensors
type SensorInfo struct {
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	GenTime time.Time `json:"gen_time"`
}

// GyroSensor produces x-y-z axes angle velocity values
type GyroSensor struct {
	SensorInfo
	AngleVelocityX float64 `json:"x_axis_angle_velocity"`
	AngleVelocityY float64 `json:"y_axis_angle_velocity"`
	AngleVelocityZ float64 `json:"z_axis_angle_velocity"`
}

// AccelSensor produces x-y-z axes gravity acceleration values
type AccelSensor struct {
	SensorInfo
	GravityAccX float64 `json:"x_axis_gravity_acceleration"`
	GravityAccY float64 `json:"y_axis_gravity_acceleration"`
	GravityAccZ float64 `json:"z_axis_grativy_acceleration"`
}

// TempSensor produces temperature and humidity values
type TempSensor struct {
	SensorInfo
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

func (s GyroSensor) String() string {
	var result []string

	st := fmt.Sprintf("Measured on %s", s.GenTime)
	result = append(result, st)

	st = fmt.Sprintf("Angle Velocity of X-axis : %f", s.AngleVelocityX)
	result = append(result, st)

	st = fmt.Sprintf("Angle Velocity of Y-axis : %f", s.AngleVelocityY)
	result = append(result, st)

	st = fmt.Sprintf("Angle Velocity of Z-axis : %f\n", s.AngleVelocityZ)
	result = append(result, st)

	return strings.Join(result, "\n")
}

// SendingOutputString of GyroSensor create sending output string with their data
func (s GyroSensor) SendingOutputString() string {
	output := fmt.Sprintf("[%s] SENT : %s : %f, %f, %f",
		s.Type, s.Name, s.AngleVelocityX, s.AngleVelocityY, s.AngleVelocityZ)

	return output
}

// ReceivingOutputString of GyroSensor create receiving output string with their data
func (s GyroSensor) ReceivingOutputString() string {
	output := fmt.Sprintf("[%s] RECEIVED : %s : %f, %f, %f",
		s.Type, s.Name, s.AngleVelocityX, s.AngleVelocityY, s.AngleVelocityZ)

	return output
}

// GenerateSensorData generates randomly GyroSensor data
func (s GyroSensor) GenerateSensorData(epsilon float64) Sensor {
	gyroSensorData := GyroSensor{
		SensorInfo: SensorInfo{
			Name:    "GyroSensor",
			Type:    "VelocitySensor",
			GenTime: time.Now(),
		},
		AngleVelocityX: faker.GenerateAngleVelocity(epsilon),
		AngleVelocityY: faker.GenerateAngleVelocity(epsilon),
		AngleVelocityZ: faker.GenerateAngleVelocity(epsilon),
	}

	return gyroSensorData
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

// SendingOutputString of AccelSensor create sending output string with their data
func (s AccelSensor) SendingOutputString() string {
	output := fmt.Sprintf("[%s] SENT : %s : %f, %f, %f",
		s.Type, s.Name, s.GravityAccX, s.GravityAccY, s.GravityAccZ)

	return output
}

// ReceivingOutputString of AccelSensor create receiving output string with their data
func (s AccelSensor) ReceivingOutputString() string {
	output := fmt.Sprintf("[%s] RECEIVED : %s : %f, %f, %f",
		s.Type, s.Name, s.GravityAccX, s.GravityAccY, s.GravityAccZ)

	return output
}

// GenerateSensorData generates randomly AccelSensor data
func (s AccelSensor) GenerateSensorData(epsilon float64) Sensor {
	accelSensorData := AccelSensor{
		SensorInfo: SensorInfo{
			Name:    "AccelerometerSensor",
			Type:    "VelocitySensor",
			GenTime: time.Now(),
		},
		GravityAccX: faker.GenerateGravityAcceleration(epsilon),
		GravityAccY: faker.GenerateGravityAcceleration(epsilon),
		GravityAccZ: faker.GenerateGravityAcceleration(epsilon),
	}

	return accelSensorData
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

// InlineString of TempSensor create inline string with their data
func (s TempSensor) InlineString() string {
	inlineStr := fmt.Sprintf("[Type : %s, Name : %s, Temp : %f, Humidity : %f]",
		s.Type, s.Name, s.Temperature, s.Humidity)

	return inlineStr
}

// SendingOutputString of TempSensor create sending output string with their data
func (s TempSensor) SendingOutputString() string {
	output := fmt.Sprintf("[%s] SENT : %s : %f, %f",
		s.Type, s.Name, s.Temperature, s.Humidity)

	return output
}

// ReceivingOutputString of TempSensor create receiving output string with their data
func (s TempSensor) ReceivingOutputString() string {
	output := fmt.Sprintf("[%s] RECEIVED : %s : %f, %f",
		s.Type, s.Name, s.Temperature, s.Humidity)

	return output
}

// GenerateSensorData generates randomly TempSensor data
func (s TempSensor) GenerateSensorData(epsilon float64) Sensor {
	tempSensorData := TempSensor{
		SensorInfo: SensorInfo{
			Name:    "TemperatureSensor",
			Type:    "AtmosphericSensor",
			GenTime: time.Now(),
		},
		Temperature: faker.GenerateTemperature(epsilon),
		Humidity:    faker.GenerateHumidity(epsilon),
	}

	return tempSensorData
}
