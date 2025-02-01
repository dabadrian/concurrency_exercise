import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;
import java.util.Arrays;
import javax.imageio.ImageIO;

public class LZ77Java {

    public static void lz77Compress(byte[] data, int workerID, byte[][] results, Object lock) throws InterruptedException {
        // Simulación de compresión LZ77
        Thread.sleep(data.length / 100);
        synchronized (lock) {
            results[workerID] = data; // En un caso real, esto sería la salida comprimida
        }
    }

    public static void compressBMP(String filePath, int numWorkers) throws IOException, InterruptedException {
        BufferedImage img = ImageIO.read(new File(filePath));
        int width = img.getWidth();
        int height = img.getHeight();
        byte[] pixelData = new byte[width * height * 3];

        for (int y = 0; y < height; y++) {
            for (int x = 0; x < width; x++) {
                int rgb = img.getRGB(x, y);
                pixelData[(y * width + x) * 3] = (byte) ((rgb >> 16) & 0xFF);
                pixelData[(y * width + x) * 3 + 1] = (byte) ((rgb >> 8) & 0xFF);
                pixelData[(y * width + x) * 3 + 2] = (byte) (rgb & 0xFF);
            }
        }

        byte[][] results = new byte[numWorkers][];
        Object lock = new Object();

        Thread[] threads = new Thread[numWorkers];
        for (int i = 0; i < numWorkers; i++) {
            final int workerID = i; // Hacemos una variable final
            final int start = workerID * pixelData.length / numWorkers;
            final int end = (workerID + 1) * pixelData.length / numWorkers;
            final byte[] segment = Arrays.copyOfRange(pixelData, start, end);

            threads[i] = new Thread(() -> {
                try {
                    lz77Compress(segment, workerID, results, lock);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            });

            threads[i].start();
        }

        for (Thread thread : threads) {
            thread.join();
        }
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        String[] imageFiles = {"random_512.bmp", "random_1024.bmp", "random_2048.bmp"};
        int[] numWorkersList = {4, 8};

        for (String filePath : imageFiles) {
            for (int numWorkers : numWorkersList) {
                long startTime = System.currentTimeMillis();
                compressBMP(filePath, numWorkers);
                System.out.println("Tiempo de ejecución para " + filePath + " con " + numWorkers + " workers: " + (System.currentTimeMillis() - startTime) + " ms");
            }
        }
    }
}
