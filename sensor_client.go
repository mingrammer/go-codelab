package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mingrammer/go-codelab/faker"
	"github.com/mingrammer/go-codelab/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(3)

	gyroTicker := time.NewTicker(500 * time.Millisecond)
	accelTicker := time.NewTicker(500 * time.Millisecond)
	tempTicker := time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case <-gyroTicker.C:
				epsilon := 5.0
				gyroSensorData := models.GyroSensor{
					Sensor: models.Sensor{
						Name:    "GyroSensor",
						Type:    "VelocitySensor",
						GenTime: time.Now(),
					},
					AngleVelocityX: faker.GenerateAngleVelocity(epsilon),
					AngleVelocityY: faker.GenerateAngleVelocity(epsilon),
					AngleVelocityZ: faker.GenerateAngleVelocity(epsilon),
				}
				url := getRequestServerUrl(8001)

				fmt.Println(gyroSensorData)

				sendJsonSensorData(url, gyroSensorData)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-accelTicker.C:
				epsilon := 10.0
				accelSensorData := models.AccelSensor{
					Sensor: models.Sensor{
						Name:    "AccelerometerSensor",
						Type:    "VelocitySensor",
						GenTime: time.Now(),
					},
					GravityAccX: faker.GenerateGravityAcceleration(epsilon),
					GravityAccY: faker.GenerateGravityAcceleration(epsilon),
					GravityAccZ: faker.GenerateGravityAcceleration(epsilon),
				}
				url := getRequestServerUrl(8002)

				fmt.Println(accelSensorData)

				sendJsonSensorData(url, accelSensorData)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-tempTicker.C:
				tempEpsilon := 2.0
				humidityEpsilon := 1.5
				tempSensorData := models.TempSensor{
					Sensor: models.Sensor{
						Name:    "TemperatureSensor",
						Type:    "AtmosphericSensor",
						GenTime: time.Now(),
					},
					Temperature: faker.GenerateTemperature(tempEpsilon),
					Humidity:    faker.GenerateHumidity(humidityEpsilon),
				}
				url := getRequestServerUrl(8003)

				fmt.Println(tempSensorData)

				sendJsonSensorData(url, tempSensorData)
			}
		}
	}()

	wg.Wait()
}

func getRequestServerUrl(port int) string {
	urlComponents := []string{"http://127.0.0.1", strconv.Itoa(port)}

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
