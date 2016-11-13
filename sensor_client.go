package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mingrammer/go-codelab/models"
)

// worker struct has necessary information to run sensor data producing job
type worker struct {
	ticker      *time.Ticker
	sensor      models.Sensor
	sensorError float64
	serverPort  int
}

// counter is for counting the number of sending request
type counter struct {
	mutex sync.Mutex
	n     int
}

func (c *counter) count() {
	c.mutex.Lock()
	c.n++
	c.mutex.Unlock()
}

func (c *counter) value() int {
	return c.n
}

// sensorWorker will be ran in goroutine to produce (random) sensor data and sends it to server
func sensorWorker(done <-chan struct{}, w worker, c *counter) {
	for {
		select {
		case <-done:
			return
		case <-w.ticker.C:
			sensorData := w.sensor.GenerateSensorData(w.sensorError)
			url := getRequestServerURL(w.serverPort)

			fmt.Println(sensorData.SendingOutputString())

			sendJSONSensorData(url, sensorData)

			c.count()
		}
	}
}

func main() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	var sendCounter counter
	var wg sync.WaitGroup
	const numWorkers = 3

	done := make(chan struct{})

	wg.Add(numWorkers)

	workerList := [numWorkers]worker{
		worker{
			ticker:      time.NewTicker(500 * time.Millisecond),
			sensor:      models.GyroSensor{},
			sensorError: 4.0,
			serverPort:  8001,
		},
		worker{
			ticker:      time.NewTicker(500 * time.Millisecond),
			sensor:      models.AccelSensor{},
			sensorError: 12.0,
			serverPort:  8002,
		},
		worker{
			ticker:      time.NewTicker(2 * time.Second),
			sensor:      models.TempSensor{},
			sensorError: 1.5,
			serverPort:  8003,
		},
	}

	for _, w := range workerList {
		go func(w worker) {
			sensorWorker(done, w, &sendCounter)
			wg.Done()
		}(w)
	}

	go func() {
		for {
			if sendCounter.value() > 100 {
				close(done)
				return
			}
		}
	}()

	wg.Wait()

	fmt.Printf("\n[Count : %d] Sending is stopped.\n", sendCounter.value())
}

// getRequestServerURL just concatenate the server url and port number to get complete server url
func getRequestServerURL(port int) string {
	urlComponents := []string{"http://127.0.0.1", strconv.Itoa(port)}

	return strings.Join(urlComponents, ":")
}

// sendJSONSensorData sends the sensor data as JSON type to server
func sendJSONSensorData(url string, sensorValues interface{}) {
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
