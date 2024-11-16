package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Order struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func main() {
	setupRabbitMQConsumer()
}

func setupRabbitMQConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"processed_orders", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare a queue:", err)
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	log.Println("Waiting for messages. To exit press CTRL+C")

	for msg := range msgs {
		var order Order
		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Println("Error unmarshalling order:", err)
			continue
		}
		sendEmail(order)
	}
}

func sendEmail(order Order) {
	to := "customer@example.com"
	subject := "Order Processed"

	// HTML and CSS template
	emailTemplate := `
	<!DOCTYPE html>
	<html lang="en">
	
	<head>
	    <meta charset="UTF-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <style>
	        body {
	            font-family: 'Arial', sans-serif;
	            background-color: #f4f4f4;
	            margin: 0;
	            padding: 0;
	            text-align: center;
	        }
	
	        .container {
	            max-width: 600px;
	            margin: 0 auto;
	            background-color: #fff;
	            padding: 20px;
	            border-radius: 10px;
	            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
	            margin-top: 20px;
	        }
	
	        h1 {
	            color: #333;
	        }
	
	        p {
	            color: #555;
	        }
	
	        .cta-button {
	            display: inline-block;
	            padding: 10px 20px;
	            background-color: #007BFF;
	            color: #fff;
	            text-decoration: none;
	            border-radius: 5px;
	            margin-top: 20px;
	        }
	
	        .footer {
	            margin-top: 20px;
	            color: #777;
	        }
	    </style>
	</head>
	
	<body>
	    <div class="container">
	        <h1>Order Processed</h1>
	        <p>Your order with ID: <strong>{{.ID}}</strong> has been successfully processed.</p>
	        <a href="#" class="cta-button">View Details</a>
	    </div>
	
	    <div class="footer">
	        <p>This is an automated email. Please do not reply.</p>
	    </div>
	</body>
	
	</html>
	`

	// Create a new template and parse the HTML
	tmpl, err := template.New("emailTemplate").Parse(emailTemplate)
	if err != nil {
		log.Println("Error parsing email template:", err)
		return
	}

	// Create a buffer to store the rendered HTML
	var buffer bytes.Buffer

	// Execute the template, passing the order data
	err = tmpl.Execute(&buffer, order)
	if err != nil {
		log.Println("Error executing email template:", err)
		return
	}

	// Compose the email body
	emailBody := fmt.Sprintf("Subject: %s\n", subject)
	emailBody += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	emailBody += buffer.String()

	// Email configuration
	emailUser := "terrence.kuhn44@ethereal.email"
	emailPassword := "2B24nePbvQxrfzhwFY"
	emailServer := "smtp.ethereal.email"
	emailPort := "587"

	// Send the email
	auth := smtp.PlainAuth("", emailUser, emailPassword, emailServer)
	err = smtp.SendMail(emailServer+":"+emailPort, auth, emailUser, []string{to}, []byte(emailBody))
	if err != nil {
		log.Println("Error sending email:", err)
	} else {
		log.Println("Email sent successfully")
	}
}

func createEmailBody(order Order) string {
	templateStr := `
	Subject: {{.Subject}}
	MIME-version: 1.0;
	Content-Type: text/html; charset="UTF-8";
 
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Order Confirmation</title>
	</head>
	<body>
		<h1>Order Processed</h1>
		<p>Your order with ID: <strong>{{.ID}}</strong> has been successfully processed.</p>
	</body>
	</html>`

	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		log.Fatal("Error parsing email template:", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, map[string]interface{}{"Subject": "Order Processed", "ID": order.ID})
	if err != nil {
		log.Fatal("Error executing email template:", err)
	}

	return body.String()
}
