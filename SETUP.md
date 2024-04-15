# Setup

## Running the services

### 1. cloning the Repository
Clone the repository to your local machine:

```bash
git clone https://github.com/sobystanly/tucows-interview
cd tucows-interview
```
### 2. Running Docker Compose
- Start PostgreSQL and RabbitMQ using Docker Compose

```bash
docker compose up -d
```
### 3.Building and Running the services
Order Management Service
- Navigate to the `order-management` directory.

```bash
cd order-management
```
Build and run order-management service using the Makefile:

```bash
make run
```
Payment Processing Service
- Navigate to the `payment-processing` directory.

```bash
cd ../payment-processing
```
Build and run the Payment Processing service using the Makefile:

```bash
make run
```

### Accessing Services
Once the services are running, you can access the APIs:
- order-management service: http://localhost:8001
- payment-processing service: http://localhost:8002

### Cleaning Up
To stop PostgreSQL and RabbitMQ and remove containers

```bash
docker compose down
```

### Additional Information
Makefile Commands
- `make build`: Build the binaries
- `make run`: Build and run the services
- `make clean`: Clean up binaries.
- `make test`: Run tests.

The documentation for each service can be found at respective readme.md of both services.

[order-management](order-management/README.md)
[payment-processing](payment-processing/README.md)

### Communication Flow Between Order and Payment Services

```mermaid
graph LR;
    subgraph "Order Service"
        order1[Order Service 1]
        order2[Order Service 2]
        order3[Order Service 3]
    end
    subgraph "RabbitMQ"
        rabbitMQ[RabbitMQ]
    end
    subgraph "Payment Service"
        payment1[Payment Service 1]
        payment2[Payment Service 2]
        payment3[Payment Service 3]
    end
    subgraph "Load Balancer"
        lb1((LB1))
    end
    subgraph "PostgreSQL"
        primary[Primary DB]
        secondary1[Secondary DB 1]
        secondary2[Secondary DB 2]
    end
    lb1 -->|Route traffic to| order1;
    lb1 -->|Route traffic to| order2;
    lb1 -->|Route traffic to| order3;
    order1 -->|Emits order events| rabbitMQ;
    order2 -->|Emits order events| rabbitMQ;
    order3 -->|Emits order events| rabbitMQ;
    rabbitMQ -->|Subscribes to payment status events| order1;
    rabbitMQ -->|Subscribes to payment status events| order2;
    rabbitMQ -->|Subscribes to payment status events| order3;
    rabbitMQ -->|Publishes order events| payment1;
    rabbitMQ -->|Publishes order events| payment2;
    rabbitMQ -->|Publishes order events| payment3;
    payment1 -->|Emits payment status events| rabbitMQ;
    payment2 -->|Emits payment status events| rabbitMQ;
    payment3 -->|Emits payment status events| rabbitMQ;
    order1 -->|Writes| primary;
    order2 -->|Writes| primary;
    order3 -->|Writes| primary;
    order1 -->|Reads| secondary1;
    order2 -->|Reads| secondary2;
    order3 -->|Reads| secondary1;
    primary -->|Replicates| secondary1;
    primary -->|Replicates| secondary2;
```