FROM eclipse-temurin:21-jdk-alpine AS builder

WORKDIR /app

COPY . .

RUN chmod +x ./mvnw && ./mvnw clean package -DskipTests

FROM eclipse-temurin:21-jre-alpine

WORKDIR /app

COPY --from=builder /app/target/*.jar app.jar

# Expose port
EXPOSE 8080

# Run application
CMD ["java", "-jar", "app.jar"]