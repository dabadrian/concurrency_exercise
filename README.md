# Ejercicios de Concurrencia con Go y Java

## Lenguaje Go
Este repositorio contiene tres ejercicios que exploran el uso de concurrencia en Go para las siguientes aplicaciones: Compresión de archivos y multiplicación de matrices. Se presentan diferentes niveles de optimización para entender el impacto de la concurrencia en el rendimiento:  

### 1. Instalación del ambiente para go

- Descargar Go: https://go.dev/dl/
- Clonar el proyecto y abrirlo en un editor de código fuente como visual studio code.

### 2. go-compression-workers

#### Descripción:
Ejercicio de compresión de imágenes BMP utilizando concurrencia, pero sin demostrar una mejora significativa en rendimiento debido a la naturaleza del procesamiento.

####  Ejecución:
```code
cd go-compression-workers
go run main.go
```

### 3. go-matrix-concurrency
#### Descripción:
Ejercicio de multiplicación de matrices utilizando concurrencia, donde cada celda de la matriz resultante se computa en una goroutine separada. Sin embargo, la sobrecarga de comunicación y sincronización no permite observar beneficios claros.

####  Ejecución:
```code
cd go-matrix-concurrency
go run main.go
```
####  Diferencias clave respecto a la versión optimizada:
- Cada celda de la matriz resultado se calcula individualmente en una tarea separada.
- Alta sobrecarga en la gestión de tareas y comunicación entre goroutines.
- Uso intensivo de canales para recolección de resultados.
- No se aprovecha la localidad de caché.

### 4. go-matrix-concurrency-optimized
#### Descripción:
Optimización del ejercicio anterior, mejorando el uso de concurrencia mediante el procesamiento por bloques en lugar de celdas individuales. Se reduce la sobrecarga y se mejora el acceso a memoria.

####  Ejecución:
```code
cd go-matrix-concurrency-optimized
go run main.go
```
####  Mejoras respecto a la versión anterior:

- Se procesan bloques de la matriz en lugar de celdas individuales.
- Menor comunicación entre goroutines y reducción del uso de canales.
- Mejor aprovechamiento de caché y menor latencia en memoria.
- Reducción de sobrecarga en la gestión de tareas.

