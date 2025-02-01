package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
	"time"
)

// Función para comprimir datos usando Deflate (LZ77)
func compressData(data []byte, workerID int, wg *sync.WaitGroup, resultChan chan float64) {
	defer wg.Done()

	start := time.Now()

	var buf bytes.Buffer
	writer, _ := flate.NewWriter(&buf, flate.BestSpeed)
	_, _ = writer.Write(data)
	_ = writer.Close()

	duration := time.Since(start).Seconds()
	resultChan <- duration
}

// Función para cargar un archivo BMP en memoria
func loadBMP(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Función para medir tiempo de ejecución con N workers y múltiples samples
func measureExecutionTime(filename string, numWorkers int, numSamples int) (float64, float64) {
	data, err := loadBMP(filename)
	if err != nil {
		fmt.Println("Error al cargar el archivo BMP:", err)
		return 0, 0
	}

	var times []float64
	for i := 0; i < numSamples; i++ {
		start := time.Now()

		var wg sync.WaitGroup
		resultChan := make(chan float64, numWorkers)

		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go compressData(data, w, &wg, resultChan)
		}

		wg.Wait()
		close(resultChan)

		var totalDuration float64
		for t := range resultChan {
			totalDuration += t
		}

		duration := time.Since(start).Seconds()
		times = append(times, duration)
	}

	mean, stddev := calculateStats(times)
	return mean, stddev
}

// Función para calcular media y desviación estándar
func calculateStats(data []float64) (mean float64, stddev float64) {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean = sum / float64(len(data))

	var varianceSum float64
	for _, v := range data {
		varianceSum += (v - mean) * (v - mean)
	}
	stddev = math.Sqrt(varianceSum / float64(len(data)))
	return
}

func main() {
	files := []string{
		"../images/random_noise_512x512.bmp",
		"../images/random_noise_1024x1024.bmp",
		"../images/random_noise_2048x2048.bmp",
	}
	numSamples := 32                 // Número de repeticiones por experimento
	workersList := []int{1, 2, 4, 8} // Diferentes números de workers

	fmt.Println("Experimento de concurrencia: Medición de tiempos")
	for _, filename := range files {
		fmt.Println("\nArchivo:", filename)
		for _, workers := range workersList {
			mean, stddev := measureExecutionTime(filename, workers, numSamples)
			fmt.Printf("Workers: %d | Tiempo medio: %.4f seg | Desviación estándar: %.4f\n", workers, mean, stddev)
		}
	}
}
