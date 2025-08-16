# Distributed Addition Service

A distributed system consisting of two microservices communicating via **gRPC** and **Kafka**.  
The system allows adding numbers through **Service A** and processing/storing the results in **Service B**.

---

## Architecture

![System Architecture](https://github.com/user-attachments/assets/23d58829-b6a4-4aea-bf51-c4b5cd92ad6f)

- **Service A**: gRPC server that receives numbers and publishes them to Kafka.
- **Kafka**: Message broker for reliable communication between services.
- **Service B**: A Spring Boot microservice that contains a Kafka consumer.  
  It listens for incoming messages from Kafka, processes the numbers, and stores the results in a file.
- **PostgreSQL**: Used as outbox storage for implementing the outbox pattern.
- **Prometheus**: For metrics collection and monitoring.
