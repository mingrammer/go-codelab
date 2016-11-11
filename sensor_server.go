package main

import (
	"encoding/json"
	"fmt"
	"github.com/mingrammer/go-codelab/models"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	logDir   = "log"
	tempLog  = "Temp.log"
	accelLog = "Accel.log"
	gyroLog  = "Gyro.log"
)

type logContent struct {
	content  string
	location string
}

type TempHandler struct {
	buf chan logContent
}

type GyroHandler struct {
	buf chan logContent
}

type AccelHandler struct {
	buf chan logContent
}

func (m *TempHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.TempSensor

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Somethink wrong")
	}
	defer req.Body.Close()

	m.buf <- logContent{content: fmt.Sprintf("%s", data), location: tempLog}
}

func (m *GyroHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.GyroSensor

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}
	defer req.Body.Close()

	m.buf <- logContent{content: fmt.Sprintf("%s", data), location: gyroLog}
}

func (m *AccelHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.AccelSensor

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}
	defer req.Body.Close()

	m.buf <- logContent{content: fmt.Sprintf("%s", data), location: accelLog}
}

func fileLogger(m chan logContent) {

	for i := range m {
		joinee := []string{logDir, i.location}
		filePath := strings.Join(joinee, "/")
		fileHandle, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err != nil {
			log.Fatal("Error Occured Opening File\n%s", err)
		}

		logger := log.New(fileHandle, "", log.LstdFlags)

		switch i.location {
		case gyroLog:
			logger.Printf("[GyroSensor Data Received]\n%s\n", i.content)
		case accelLog:
			logger.Printf("[AccelSensor Data Received]\n%s\n", i.content)
		case tempLog:
			logger.Printf("[TempSensor Data Received]\n%s\n", i.content)
		}

		defer fileHandle.Close()
	}
}

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
