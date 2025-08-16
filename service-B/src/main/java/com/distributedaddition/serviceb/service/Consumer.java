package com.distributedaddition.serviceb.service;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
@AllArgsConstructor
@Slf4j
@Component
public class Consumer {

    private FileStorage fileStorage;

    @KafkaListener(topics = "addition", groupId = "service-b-group")
    public void handleNumber(String message){
        try {
            log.info("Received number form Kafka: {}", message);

            int number = Integer.parseInt(message.trim());
            fileStorage.addNumber(number);

            log.info("Successfully processed the number: {}", number);
        } catch (NumberFormatException e){
            log.error("Failed to parse message as integer: {}", message, e);
        } catch (Exception e) {
            log.error("Error processing Kafka message: {}", message, e);
        }
    }
}
