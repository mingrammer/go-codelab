package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mingrammer/go-codelab/models"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	gyroTicker := time.NewTicker(500 * time.Millisecond)
	accelTicker := time.NewTicker(500 * time.Millisecond)
	tempTicker := time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case <-gyroTicker.C:
				url := getRequestServerUrl(8001)
				gyroSensorData := models.GyroSensor{
					Sensor: models.Sensor{
						Name:    "GyroSensor",
						Type:    "VelocitySensor",
						GenTime: time.Now(),
					},
					AngleVelocityX: 32.54,
					AngleVelocityY: 35.12,
					AngleVelocityZ: 61.23,
				}

				fmt.Println(gyroSensorData)

				sendJsonSensorData(url, gyroSensorData)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-accelTicker.C:
				url := getRequestServerUrl(8002)
				accelSensorData := models.AccelSensor{
					Sensor: models.Sensor{
						Name:    "AccelerometerSensor",
						Type:    "VelocitySensor",
						GenTime: time.Now(),
					},
					GravityAccX: 41.31,
					GravityAccY: 81.36,
					GravityAccZ: 46.19,
				}

				fmt.Println(accelSensorData)

				sendJsonSensorData(url, accelSensorData)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-tempTicker.C:
				url := getRequestServerUrl(8003)
				tempSensorData := models.TempSensor{
					Sensor: models.Sensor{
						Name:    "TemperatureSensor",
						Type:    "AtmosphericSensor",
						GenTime: time.Now(),
					},
					Temperature: 84.13,
					Humidity:    76.12,
				}

				fmt.Println(tempSensorData)

				sendJsonSensorData(url, tempSensorData)
			}
		}
	}()
}

func getRequestServerUrl(port int) string {
	urlComponents := []string{"http://127.0.0.1", string(port)}

	return strings.Join(urlComponents, ":")
}

func sendJsonSensorData(url string, sensorValues interface{}) {
	jsonBytes, err := json.Marshal(sensorValues)
	if err != nil {
		log.Fatal("Error occurs when mashaling the gyro sensor values")
	}

	buff := bytes.NewBuffer(jsonBytes)

	resp, err := http.Post(url, "application/json", buff)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Error occurs when request the post data")

	}
}
