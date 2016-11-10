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

type TempHandler struct {
}

type GyroHandler struct {
}

type AccelHandler struct {
}

func (m *TempHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.TempSensor

	fileHandle, err := os.OpenFile("./log/Temp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error Opening Temp.log")
	}

	logger := log.New(fileHandle, "", log.LstdFlags)
	defer fileHandle.Close()

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Somethink wrong")
	}

	defer req.Body.Close()
	logger.Printf("[%s of %s Received]\n%s", data.Name, data.Type, data)
	log.Println("Temperature Data Incoming")
}

func (m *GyroHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.GyroSensor

	fileHandle, err := os.OpenFile("./log/Gyro.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
		log.Fatal("Error Opening Gyro.log")
	}

	logger := log.New(fileHandle, "", log.LstdFlags)
	defer fileHandle.Close()

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}

	defer req.Body.Close()
	logger.Printf("[%s of %s Received]\n%s", data.Name, data.Type, data)
	log.Println("Gyro Sensor Data Incoming")
}

func (m *AccelHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var data models.AccelSensor

	fileHandle, err := os.OpenFile("./log/Accel.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error Opening Accel.log")
	}

	logger := log.New(fileHandle, "", log.LstdFlags)
	defer fileHandle.Close()

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Something wrong")
	}

	defer req.Body.Close()
	logger.Printf("[%s of %s Received]\n%s", data.Name, data.Type, data)
	log.Println("Accelorator Sensor Data Incoming")
}

func main() {
	var wg sync.WaitGroup

	wg.Add(4)

	gyroHander := &GyroHandler{}
	accelHandler := &AccelHandler{}
	tempHandler := &TempHandler{}
	// buf := make(chan interface{})

	go http.ListenAndServe(":8001", gyroHander)
	go http.ListenAndServe(":8002", accelHandler)
	go http.ListenAndServe(":8003", tempHandler)

	wg.Wait()
}
