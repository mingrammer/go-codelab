// This file represents the server for sensors
// It opens communication links (for each sensor) through http protocol
// Port # 8001 for Gyroscope Sensor,
// Port # 8002 for Accelerator Sensor,
// Port # 8003 for Temperature Sensor.
// If the sensor client sends data through appropriate port,
// each http handler is initiated, decoding data to appropriate JSON file.
// Those datas are then logged to each sensor's log file
// (log/Accel.log , log/Temp.log , log/Gyro.log)

package main

// 'modles' package		 : Stores basic sensor information in srtuct form
// 'net/http' package	 : To serve http connection and handle requests
// 'os' package			 : To open files for logging
// 'strings' package	 : To join strings for file path
// 'sync' package		 : To use 'WaitGroup' to hold main thread while goroutines are working

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/mingrammer/go-codelab/models"
)

// String constants that are used in sensor_server.go
// logDir	: Directory name to store logs
// tempLog	: File name to store temperature sensor log
// accelLog	: File name to store accelerator sensor log
// gyroLog	: File name to store gyroscope sensro log

const (
	logDir   = "log"
	tempLog  = "Temp.log"
	accelLog = "Accel.log"
	gyroLog  = "Gyro.log"
)

// This logContent struct is to store data (string) would be written in log file
// content 	: Actual string data, that would be logged in file
// location	: Indicator that where the 'content' should be stored (or written)

type logContent struct {
	content  string
	location string
	sensorName string

}

// Three structs below are to implement ServeHTTP method
// Each handler stores the pointer to data logging channel
// Also, channel is bidirectional in these handlers, which only can store data in channel
// TempHandler 	: Temperature sensor handler to implement ServeHTTP method
// GyroHandler	: Gyroscopte sensor handler to implement ServeHTTP method
// accelHandler : Accelerator sensro handler to implement ServeHTTP method

type TempHandler struct {
	buf chan<- logContent
}

type GyroHandler struct {
	buf chan<- logContent
}

type AccelHandler struct {
	buf chan<- logContent
}

// These methods are to handle request from each port.
// Body of request (which is in JSON format) is decoded and allocated to TempSensor variable
// and string, of which data would be saved in log fileis sent to logging channel
// Note that BOdy of http request CAN NOT BE unmarshalled, because body of request is in array of bytes.

func (m *TempHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.TempSensor

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}
	defer req.Body.Close()

	m.buf <- logContent{content: fmt.Sprintf("%s", data), location: tempLog, sensorName: data.Name}
}

func (m *GyroHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.GyroSensor

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}
	defer req.Body.Close()

	m.buf <- logContent{content: fmt.Sprintf("%s", data), location: gyroLog, sensorName: data.Name}
}

func (m *AccelHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.AccelSensor

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}
	defer req.Body.Close()

	m.buf <- logContent{content: fmt.Sprintf("%s", data), location: accelLog, sensorName: data.Name}
}

// This method loggs the content of sensor data
// This method waits incoming data from 'logContent' channel at range block,
// where ServeHTTP mehthod sends log data
// When data is detected, for/range block immediately processes data
// It checks the location, where the data should be stored,
// and opens file of desired location (by joining string constants)
// Note that channel used in this method is also BIDIRECTIONAL,
// You only can pop the data from channel.

func fileLogger(m <-chan logContent) {

	// When program is initialized, fileLogger checks if 'log' dir exists.
	// If it does not, it creates directory for the first time.
	// If there is problem with creating it, it throws panic and aborts the application.

	dir, _ := os.Open("log")
	dirInfo, _ := dir.Stat()

	if dirInfo == nil {
		err := os.Mkdir("log", os.ModePerm)

		if err != nil {
			log.Fatal("Error creating directory 'log'\n", err)
		}
	}
	dir.Close()

	// This part continuously wait for incoming data through channel

	for i := range m {
		joinee := []string{logDir, i.location}
		filePath := strings.Join(joinee, "/")

		fileHandle, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err != nil {
			log.Fatal("Error Opening File\n", err)
		}

		logger := log.New(fileHandle, "", log.LstdFlags)

		logger.Printf("[%s Data Received]\n%s\n", i.sensorName, i.content)

		defer fileHandle.Close()
	}
}

// This method is for main thread. It starts go application.
// It first creates handler to serve http serve for each sensors
// sync.WaitGroup is to hold main thread while go routines are still working
// Main thread indicates WaitGroup to wait 4 routines to stop.
// After goroutines we need is created, main thread has WaitGroup to wait goroutines.

func main() {
	var wg sync.WaitGroup

	wg.Add(4)

	logBuf := make(chan logContent)
	gyroHander := &GyroHandler{buf: logBuf}
	accelHandler := &AccelHandler{buf: logBuf}
	tempHandler := &TempHandler{buf: logBuf}

	go http.ListenAndServe(":8001", gyroHander)
	go http.ListenAndServe(":8002", accelHandler)
	go http.ListenAndServe(":8003", tempHandler)
	go fileLogger(logBuf)

	wg.Wait()
}
