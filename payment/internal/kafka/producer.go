package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/maksroxx/DeliveryService/payment/internal/models"
)

type Producer struct {
	producer sarama.SyncProducer
	topics   []string
}

func NewProducer(cfg ConfigProducer) (*Producer, error) {
	saramaConfig := NewProducerConfig()

	producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}
	return &Producer{
		producer: producer,
		topics:   cfg.Topic,
	}, nil
}

func (p *Producer) PaymentMessage(payment models.Payment, userID string) error {
	paymentJson, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	headers := []sarama.RecordHeader{
		{
			Key:   []byte("User-ID"),
			Value: []byte(userID),
		},
	}

	msg := &sarama.ProducerMessage{
		Headers: headers,
		Key:     sarama.StringEncoder(payment.PackageID),
		Topic:   p.topics[0],
		Value:   sarama.ByteEncoder(paymentJson),
	}
	_, _, err = p.producer.SendMessage(msg)
	return err
}

func (p *Producer) PaymentAucitonMessage(payment models.Payment, userID string) error {
	paymentJson, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	headers := []sarama.RecordHeader{
		{
			Key:   []byte("User-ID"),
			Value: []byte(userID),
		},
	}

	msg := &sarama.ProducerMessage{
		Headers: headers,
		Key:     sarama.StringEncoder(payment.PackageID),
		Topic:   p.topics[1],
		Value:   sarama.ByteEncoder(paymentJson),
	}
	_, _, err = p.producer.SendMessage(msg)
	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
