package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultDialAttempts = 10
	defaultDialInterval = 3 * time.Second
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQ(url, queueName string) (*RabbitMQ, error) {
	conn, err := dialWithRetry(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		queue:   queueName,
	}, nil
}

func dialWithRetry(url string) (*amqp.Connection, error) {
	var lastErr error
	for attempt := 1; attempt <= defaultDialAttempts; attempt++ {
		conn, err := amqp.Dial(url)
		if err == nil {
			if attempt > 1 {
				log.Printf("Connected to RabbitMQ after %d attempts", attempt)
			}
			return conn, nil
		}

		lastErr = err
		log.Printf("RabbitMQ connection attempt %d/%d failed: %v", attempt, defaultDialAttempts, err)
		if attempt < defaultDialAttempts {
			time.Sleep(defaultDialInterval)
		}
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", defaultDialAttempts, lastErr)
}

func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *RabbitMQ) Publish(ctx context.Context, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = r.channel.PublishWithContext(
		ctx,
		"",
		r.queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to queue %s", r.queue)
	return nil
}

func (r *RabbitMQ) Consume(ctx context.Context, handler func([]byte) error) error {
	msgs, err := r.channel.Consume(
		r.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Started consuming messages from queue %s", r.queue)

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer context cancelled, stopping...")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}

			log.Printf("Received message from queue %s", r.queue)

			err := handler(msg.Body)
			if err != nil {
				log.Printf("Error handling message: %v", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
				log.Println("Message acknowledged")
			}
		}
	}
}
