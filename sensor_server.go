package main

import (
	"encoding/json"
	"fmt"
	"github.com/mingrammer/go-codelab/models"
	"log"
	"net/http"
	"os"
	"sync"
)

type logContent struct {
	content  string
	location string
}

type logBuffer struct {
	buffer chan logContent
	mux    sync.Mutex
}

type TempHandler struct {
	buf *logBuffer
}

type GyroHandler struct {
	buf *logBuffer
}

type AccelHandler struct {
	buf *logBuffer
}

func (m *TempHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.TempSensor

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Somethink wrong")
	}

	defer req.Body.Close()

	m.buf.mux.Lock()
	m.buf.buffer <- logContent{content: fmt.Sprintf("%s", data), location: "Temp"}
	defer m.buf.mux.Unlock()
}

func (m *GyroHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.GyroSensor

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}

	defer req.Body.Close()
	m.buf.mux.Lock()
	m.buf.buffer <- logContent{content: fmt.Sprintf("%s", data), location: "Gyro"}
	defer m.buf.mux.Unlock()
}

func (m *AccelHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.AccelSensor

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}

	defer req.Body.Close()

	m.buf.mux.Lock()
	m.buf.buffer <- logContent{content: fmt.Sprintf("%s", data), location: "Accel"}
	defer m.buf.mux.Unlock()
}

func (m *logBuffer) fileLogger() {
	tempLog, tempErr := os.OpenFile("./log/Temp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	gyroLog, gyroErr := os.OpenFile("./log/Gyro.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	accelLog, accelErr := os.OpenFile("./log/Accel.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if tempErr != nil || gyroErr != nil || accelErr != nil {
		fmt.Println("Error Opening File. See Below for Reason")
		log.Fatal(tempErr)
		log.Fatal(gyroErr)
		log.Fatal(accelErr)
	}

	tempLogger := log.New(tempLog, "", log.LstdFlags)
	gyroLogger := log.New(gyroLog, "", log.LstdFlags)
	accelLogger := log.New(accelLog, "", log.LstdFlags)

	defer tempLog.Close()
	defer gyroLog.Close()
	defer accelLog.Close()

	for i := range m.buffer {
		switch i.location {
		case "Gyro":
			gyroLogger.Printf("[GyroSensor Data Received]\n%s\n", i.content)
		case "Accel":
			accelLogger.Printf("[AccelSensor Data Received]\n%s\n", i.content)
		case "Temp":
			tempLogger.Printf("[TempSensor Data Received]\n%s\n", i.content)
		}
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(4)

	logBuf := &logBuffer{buffer: make(chan logContent)}
	gyroHander := &GyroHandler{buf: logBuf}
	accelHandler := &AccelHandler{buf: logBuf}
	tempHandler := &TempHandler{buf: logBuf}

	go http.ListenAndServe(":8001", gyroHander)
	go http.ListenAndServe(":8002", accelHandler)
	go http.ListenAndServe(":8003", tempHandler)
	go logBuf.fileLogger()

	wg.Wait()
}
