package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const chunkSize = 4096 // Tamaño del bloque en bytes

type Chunk struct {
	Index int
	Data  []byte
}

type CompressedChunk struct {
	Index int
	Data  []byte
}

func readFileChunks(filePath string) ([]Chunk, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunks []Chunk
	index := 0

	for {
		buffer := make([]byte, chunkSize)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		chunks = append(chunks, Chunk{Index: index, Data: buffer[:n]})
		index++
	}

	return chunks, nil
}

func compressWorker(id int, wg *sync.WaitGroup, jobs <-chan Chunk, results chan<- CompressedChunk, times chan<- time.Duration) {
	defer wg.Done()

	for chunk := range jobs {
		start := time.Now()

		var buf bytes.Buffer
		writer, err := flate.NewWriter(&buf, flate.BestCompression)
		if err != nil {
			fmt.Println("Error creando el escritor de compresión:", err)
			continue
		}

		_, err = writer.Write(chunk.Data)
		if err != nil {
			fmt.Println("Error al comprimir el chunk:", err)
			continue
		}

		writer.Close()
		results <- CompressedChunk{Index: chunk.Index, Data: buf.Bytes()}

		// Guardar tiempo de compresión
		times <- time.Since(start)
	}
}

func writeCompressedFile(outputPath string, compressedChunks []CompressedChunk) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, chunk := range compressedChunks {
		_, err := file.Write(chunk.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	inputFile := "sample.bmp"
	outputFile := "output.deflate"
	numWorkers := 16 // Parametrización del número de workers

	startTime := time.Now()

	// Leer archivo en chunks
	chunks, err := readFileChunks(inputFile)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	// Canales para comunicación
	jobs := make(chan Chunk, len(chunks))
	results := make(chan CompressedChunk, len(chunks))
	times := make(chan time.Duration, len(chunks))

	// WaitGroup para sincronización
	var wg sync.WaitGroup

	// Iniciar goroutines de compresión
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go compressWorker(i, &wg, jobs, results, times)
	}

	// Enviar chunks a los workers
	for _, chunk := range chunks {
		jobs <- chunk
	}
	close(jobs)

	// Esperar a que terminen los workers
	wg.Wait()
	close(results)
	close(times)

	// Recoger resultados en orden
	compressedChunks := make([]CompressedChunk, 0, len(chunks))
	var totalTime time.Duration

	for compressedChunk := range results {
		compressedChunks = append(compressedChunks, compressedChunk)
	}
	for t := range times {
		totalTime += t
	}

	// Escribir archivo comprimido
	err = writeCompressedFile(outputFile, compressedChunks)
	if err != nil {
		fmt.Println("Error al escribir el archivo comprimido:", err)
		return
	}

	elapsed := time.Since(startTime)

	fmt.Println("Compresión completada.")
	fmt.Println("Archivo guardado en:", outputFile)
	fmt.Printf("Tiempo total de compresión: %v\n", elapsed)
	fmt.Printf("Tiempo promedio por bloque: %v\n", totalTime/time.Duration(len(chunks)))
	fmt.Printf("Número de goroutines usadas: %d\n", numWorkers)
}
