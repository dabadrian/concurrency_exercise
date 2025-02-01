# Comparación de Compresión Concurrente en Go y Java

## Lenguaje Go
Este repositorio contiene dos implementaciones para evaluar el rendimiento de compresión de imágenes BMP usando múltiples hilos de ejecución. Se han implementado versiones en Go y Java para comparar su desempeño en entornos concurrentes.

Cada implementación:

- Carga un archivo BMP en memoria.

- Utiliza múltiples hilos (workers) para comprimir los datos con el algoritmo Deflate (LZ77).

- Mide el tiempo de ejecución para diferentes cantidades de hilos.

- Calcula la media y desviación estándar de los tiempos de ejecución.

### 1. Requisitos

- Clonar el proyecto y abrirlo en un editor de código fuente como visual studio code.

#### Para Go

- Instalar Go: https://go.dev/dl/

- Asegurarse de que los archivos BMP estén en ../images/

#### Para Java

- Instalar JDK 17+

- Configurar JAVA_HOME en las variables de entorno

- Asegurarse de que los archivos BMP estén en ../images/



### 2. Ejecución

#### Go

```code
cd go-concurrency-lz77
go run main.go
```

#### Java

```code
cd java-concurrency-lz77
javac ParallelCompression.java
java ParallelCompression
```


### 3. Arquitectura de las Implementaciones

#### Go

- Usa sync.WaitGroup para la sincronización de hilos.

- Utiliza canales (chan) para recolectar los tiempos de ejecución de cada worker.

- Mide el tiempo total por iteración y luego calcula la media y desviación estándar.

#### Java

- Usa Thread y ReentrantLock para la sincronización.

- Almacena los tiempos en una List<Double> protegida por un lock.

- Mide el tiempo total por iteración y calcula estadísticas.

### 4. Resultados y Comparación
Ambas implementaciones ejecutan los mismos experimentos:

- Archivos BMP: 512x512, 1024x1024, 2048x2048.

- Workers: 1, 2, 4, 8, 16.

- Número de muestras por prueba: 32.

Se recomienda analizar los tiempos de ejecución para evaluar las diferencias entre los modelos de concurrencia de Go y Java.