package com.distributedaddition.serviceb.service;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.io.*;
import java.util.concurrent.locks.ReentrantReadWriteLock;

@Service
public class FileStorage {

    private static final Logger logger = LoggerFactory.getLogger(FileStorage.class);
    private static final String FILE_PATH = "numbers.txt";
    private final ReentrantReadWriteLock lock = new ReentrantReadWriteLock();

    public void addNumber(int number) {
        lock.writeLock().lock();
        try {
            int currentNumber = getLastNumber();
            int newNumber = currentNumber + number;

            try (BufferedWriter writer = new BufferedWriter(new FileWriter(FILE_PATH))) {
                writer.write(String.valueOf(newNumber));
            } catch (IOException e) {
                logger.error("Error writing to file: {}", FILE_PATH, e);
            }
        } finally {
            lock.writeLock().unlock();
        }
    }

    public int getLastNumber() {
        lock.readLock().lock();
        try {
            File file = new File(FILE_PATH);
            if (!file.exists()) {
                return 0;
            }

            try (BufferedReader reader = new BufferedReader(new FileReader(FILE_PATH))) {
                String line = reader.readLine();
                if (line != null && !line.trim().isEmpty()) {
                    return Integer.parseInt(line.trim());
                }
            } catch (IOException | NumberFormatException e) {
                logger.error("Error reading from file: {}", FILE_PATH, e);
            }
            return 0;
        } finally {
            lock.readLock().unlock();
        }
    }
}
