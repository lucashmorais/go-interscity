package measurement

import "time"

var MeasurementArrays [][]time.Duration

func InitMeasurementArrays(numArrays int, arraySize int) {
	count = 0
	MeasurementArrays = make([][]time.Duration, 0, numArrays)
	for i := 0; i < numArrays; i++ {
		MeasurementArrays = append(MeasurementArrays, make([]time.Duration, 0, arraySize))
	}
}

func StartMeasuringRequest() time.Time {
	return time.Now()
}

var count int

func StopMeasuringRequest(startTime time.Time, requestType int) {
	elapsedTime := time.Since(startTime)
	count++
	if count%1000 == 0 {
		println("[StopMeasuringRequest]: Request latency (us): ", elapsedTime.Microseconds())
	}
	MeasurementArrays[requestType] = append(MeasurementArrays[requestType], elapsedTime)
}
