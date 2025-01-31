package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	MatrixSize = 200 // Tamaño de la matriz cuadrada NxN
	Samples    = 32  // Número de ejecuciones por configuración
)

// Multiplicación de matrices con concurrencia
type Task struct {
	row   int
	col   int
	value int
}

func worker(jobs <-chan Task, results chan<- Task, wg *sync.WaitGroup) {
	for task := range jobs {
		results <- Task{task.row, task.col, task.value}
	}
	wg.Done()
}

func multiplyMatrices(A, B [][]int, workers int) [][]int {
	N := len(A)
	C := make([][]int, N)
	for i := range C {
		C[i] = make([]int, N)
	}

	jobs := make(chan Task, N*N)
	results := make(chan Task, N*N)
	var wg sync.WaitGroup

	// Lanzamos los workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Enviar tareas
	go func() {
		for i := 0; i < N; i++ {
			for j := 0; j < N; j++ {
				sum := 0
				for k := 0; k < N; k++ {
					sum += A[i][k] * B[k][j]
				}
				jobs <- Task{i, j, sum}
			}
		}
		close(jobs)
	}()

	// Recibir resultados
	go func() {
		wg.Wait()
		close(results)
	}()

	for task := range results {
		C[task.row][task.col] = task.value
	}

	return C
}

func generateMatrix(N int) [][]int {
	mat := make([][]int, N)
	for i := range mat {
		mat[i] = make([]int, N)
		for j := range mat[i] {
			mat[i][j] = rand.Intn(10)
		}
	}
	return mat
}

func meanAndStdDev(times []float64) (float64, float64) {
	sum := 0.0
	for _, t := range times {
		sum += t
	}
	mean := sum / float64(len(times))

	variance := 0.0
	for _, t := range times {
		variance += (t - mean) * (t - mean)
	}
	stdDev := sqrt(variance / float64(len(times)))

	return mean, stdDev
}

func sqrt(value float64) float64 {
	if value == 0 {
		return 0
	}
	z := value / 2
	for i := 0; i < 20; i++ {
		z -= (z*z - value) / (2 * z)
	}
	return z
}

func main() {
	rand.Seed(time.Now().UnixNano())
	workersList := []int{1, 2, 4, 8, 16}
	fmt.Println("Experimento de concurrencia: Medición de tiempos")

	for _, workers := range workersList {
		times := make([]float64, Samples)
		//fmt.Printf("Workers: %d\n", workers)

		for i := 0; i < Samples; i++ {
			A := generateMatrix(MatrixSize)
			B := generateMatrix(MatrixSize)
			start := time.Now()
			_ = multiplyMatrices(A, B, workers)
			times[i] = time.Since(start).Seconds()
		}

		mean, stdDev := meanAndStdDev(times)
		fmt.Printf("Workers: %d | Tiempo medio: %.4f seg | Desviación estándar: %.4f\n", workers, mean, stdDev)
	}
}
