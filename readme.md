# Order Management Service

## Overview

The **Order Management Service** is a scalable backend application designed to handle order data efficiently using modern technologies. It integrates with **Apache Kafka** for message processing, uses **PostgreSQL** for data persistence, and incorporates **BigCache** for fast in-memory caching. The service is structured with clean architecture principles, making it easy to maintain and extend.

### Features

- **Kafka Integration**: Processes incoming messages from the `orders` topic in Kafka.
- **PostgreSQL Database**: Stores detailed order data with transactional integrity.
- **BigCache**: Provides high-performance caching for frequently accessed orders.
- **Graceful Shutdown**: Ensures the service shuts down cleanly without data loss.
- **Docker Support**: Comes with a `docker-compose` configuration for Kafka.

---

## Getting Started

## Running the Application

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/order-management-service.git
   cd order-management-service
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```

2.  Install dependencies:
    ```bash
    go mod tidy
    ```

3. Start Kafka with Docker Compose

     ```bash
     docker-compose up -d
     ```
4. Initialize PostgreSQL Database. Put the PG_URL in the .env
5. Start the server
    ```bash
   go run cmd/main.go
   ```

## Testing

```bash
go test ./... -v
```

