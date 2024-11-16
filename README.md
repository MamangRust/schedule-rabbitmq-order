### Schedule RabbitMQ Order


Schedule RabbitMQ Order is a project that demonstrates the use of Apache Kafka for scheduling and processing orders. It consists of several components, including order producer, order processor, email service, and Kafka broker.

## Features

- Rabbit-mq order processing system
- Order scheduling and processing
- Email notification on order completion

## Running the project

### Setup Email in Ethereal

`https://ethereal.email/`

```
emailUser := ""
emailPassword := ""
emailServer := "smtp.ethereal.email"
emailPort := "587"
```

### Build and run

```bash
docker-compose up -d --build
```

### Test Curl

```
    curl -X POST -H "Content-Type: application/json
" -d '{"id": 235, "status": "processed"}' http://172.20.0.4:5000/placeOrder
```
