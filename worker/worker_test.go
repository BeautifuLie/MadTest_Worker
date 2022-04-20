package worker

import (
	"program/logging"
	"program/mockmessages"
	"program/model"
	"program/rmq"
	"program/storage"

	"testing"
)

func TestWorker(t *testing.T) {
	logger := logging.InitZapLog()
	s3stor, err := storage.NewS3Storage()
	if err != nil {
		logger.Errorw("Error during connect to AWS services", "error", err)
	}
	rabbitstor, err := rmq.NewRabbitStorage("amqp://denys:lafazan@localhost:5672/")
	if err != nil {
		logger.Errorw("Error during connect to RabbitMQ broker", "error", err)
	}
	closeCh := make(chan struct{})
	w := NewWorker(rabbitstor, s3stor, closeCh)
	mockCh := make(chan model.Message, 10)
	msgChan := mockmessages.MockMessages(mockCh)

	for i := 1; i <= len(msgChan); i++ {
		go w.processMessage(msgChan, i, logger)
	}
}
