package main

import (
	"encoding/json"
	"github.com/mingrammer/go-codelab/models"
	"fmt"
	"net/http"
	"sync"
	"log"
	"os"
)

type logBuffer struct {
	buffer chan interface{}
	mux sync.Mutex
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

	fileHandle, err := os.OpenFile("./log/Temp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error Opening Temp.log")
	}

	// logger := log.New(fileHandle, "", log.LstdFlags)
	defer fileHandle.Close()

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Somethink wrong")
	}

	defer req.Body.Close()
	// log.Println("Temperature Data Incoming")
	m.buf.mux.Lock()
	defer m.buf.mux.Unlock()
	fmt.Println("Locked at Temp")
	m.buf.buffer <- data
	fmt.Println("Allocated at Temp")
}

func (m *GyroHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.GyroSensor

	fileHandle, err := os.OpenFile("./log/Gyro.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
		log.Fatal("Error Opening Gyro.log")
	}

	// logger := log.New(fileHandle, "", log.LstdFlags)
	defer fileHandle.Close()

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}

	defer req.Body.Close()
	// log.Println("Gyro Sensor Data Incoming")
	m.buf.mux.Lock()
	defer m.buf.mux.Unlock()
	fmt.Println("Locked at Gyro")
	m.buf.buffer <- data
	fmt.Println("Allocated at Gyro")
}

func (m *AccelHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.AccelSensor

	fileHandle, err := os.OpenFile("./log/Accel.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error Opening Accel.log")
	}

	// logger := log.New(fileHandle, "", log.LstdFlags)
	defer fileHandle.Close()

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}

	defer req.Body.Close()
	// log.Println("Accelorator Sensor Data Incoming")

	m.buf.mux.Lock()
	defer m.buf.mux.Unlock()
	fmt.Println("Locked at Accel")
	m.buf.buffer <- data
	fmt.Println("Allocated at Accel")

}

func fileLogger(buf *logBuffer){
	for i := range buf.buffer {
		fmt.Println("Locked at fileLogger")
		switch v := i.(type) {
		case models.GyroSensor:
			fmt.Println("[%s of %s Received]\n%s", v.Name, v.Type, v)
		case models.AccelSensor:
			fmt.Println("[%s of %s Received]\n%s", v.Name, v.Type, v)
		case models.TempSensor:
			fmt.Println("[%s of %s Received]\n%s", v.Name, v.Type, v)
		}
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(4)

	logBuf := &logBuffer{buffer : make(chan interface{})}
	gyroHander := &GyroHandler{buf : logBuf}
	accelHandler := &AccelHandler{buf : logBuf}
	tempHandler := &TempHandler{buf : logBuf}

	go http.ListenAndServe(":8001", gyroHander)
	go http.ListenAndServe(":8002", accelHandler)
	go http.ListenAndServe(":8003", tempHandler)
	go fileLogger(logBuf)

	wg.Wait()
}
