import java.io.*;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.zip.DeflaterOutputStream;

public class ParallelCompression {
    private static final String[] FILES = {
        "../images/random_noise_512x512.bmp",
        "../images/random_noise_1024x1024.bmp",
        "../images/random_noise_2048x2048.bmp"
    };
    private static final int[] WORKERS_LIST = {1, 2, 4, 8};
    private static final int NUM_SAMPLES = 32;
    private static final Lock lock = new ReentrantLock();

    public static void main(String[] args) {
        System.out.println("Experimento de concurrencia: Medición de tiempos en Java");
        for (String filename : FILES) {
            System.out.println("Procesando: " + filename);
            for (int workers : WORKERS_LIST) {
                double[] results = measureExecutionTime(filename, workers, NUM_SAMPLES);
                System.out.printf("Workers: %d | Tiempo medio: %.4f seg | Desviación estándar: %.4f\n", workers, results[0], results[1]);
            }
            System.out.println();
        }
    }

    private static double[] measureExecutionTime(String filename, int numWorkers, int numSamples) {
        byte[] data = loadBMP(filename);
        if (data == null) return new double[]{0, 0};
        
        List<Double> times = new ArrayList<>();
        
        for (int i = 0; i < numSamples; i++) {
            long start = System.nanoTime();
            List<Thread> threads = new ArrayList<>();
            List<Double> durations = new ArrayList<>();
            
            for (int w = 0; w < numWorkers; w++) {
                Thread thread = new Thread(() -> {
                    long tStart = System.nanoTime();
                    compressData(data);
                    double duration = (System.nanoTime() - tStart) / 1e9;
                    lock.lock();
                    try { durations.add(duration); }
                    finally { lock.unlock(); }
                });
                threads.add(thread);
                thread.start();
            }
            
            for (Thread thread : threads) {
                try { thread.join(); } catch (InterruptedException e) { e.printStackTrace(); }
            }
            
            double totalDuration = durations.stream().mapToDouble(Double::doubleValue).sum();
            double elapsedTime = (System.nanoTime() - start) / 1e9;
            times.add(elapsedTime);
        }
        
        return calculateStats(times);
    }

    private static byte[] loadBMP(String filename) {
        try {
            return Files.readAllBytes(Paths.get(filename));
        } catch (IOException e) {
            System.out.println("Error al cargar el archivo: " + filename);
            return null;
        }
    }

    private static void compressData(byte[] data) {
        try (ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
             DeflaterOutputStream deflaterOutputStream = new DeflaterOutputStream(byteArrayOutputStream)) {
            deflaterOutputStream.write(data);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    private static double[] calculateStats(List<Double> times) {
        double sum = times.stream().mapToDouble(Double::doubleValue).sum();
        double mean = sum / times.size();
        
        double varianceSum = times.stream().mapToDouble(t -> Math.pow(t - mean, 2)).sum();
        double stddev = Math.sqrt(varianceSum / times.size());
        
        return new double[]{mean, stddev};
    }
}
