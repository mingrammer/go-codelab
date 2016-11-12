package faker

import (
	"math/rand"
	"time"
)

// getRandValueInRange creates random value in (range size - error, range size + error)
func getRandValueInRange(rangeSize int, epsilon float64) float64 {
	rand.Seed(time.Now().UTC().UnixNano())

	maxDataWithError := float64(rangeSize) + 2*epsilon

	dataInRange := rand.Float64()*maxDataWithError - epsilon

	return dataInRange
}

// GenerateAngleVelocity generates random value in (0 - epsilon, 360 + epsilon)
func GenerateAngleVelocity(epsilon float64) float64 {
	return getRandValueInRange(360, epsilon)
}

// GenerateGravityAcceleration generates random value in (0 - epsilon, 1023 + epsilon)
func GenerateGravityAcceleration(epsilon float64) float64 {
	return getRandValueInRange(1023, epsilon)
}

// GenerateTemperature generates random value in (0 - epsilon, 104 + epsilon)
func GenerateTemperature(epsilon float64) float64 {
	return getRandValueInRange(104, epsilon)
}

// GenerateHumidity generates random value in (0 - epsilon, 100 + epsilon)
func GenerateHumidity(epsilon float64) float64 {
	return getRandValueInRange(100, epsilon)
}
