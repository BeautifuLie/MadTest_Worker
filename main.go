package main

import (
	"os"
	"os/signal"
	"program/logging"
	"program/sqs"
	"program/storage"
	"program/worker"
	"syscall"

	"github.com/joho/godotenv"

	_ "gocloud.dev/blob/s3blob"
)

func main() {
	logger := logging.InitZapLog()
	sch := make(chan os.Signal, 1)
	err := godotenv.Load(".env")
	if err != nil {
		logger.Errorw("Error during load environments", err)
	}

	s3stor, err := storage.NewS3Storage()
	if err != nil {
		logger.Errorw("Error during connect to AWS services", "error", err)
	}
	sqsstor, err := sqs.NewSqsStorage()
	if err != nil {
		logger.Errorw("Error during connect to AWS services", "error", err)
	}
	// rabbitstor, err := rmq.NewRabbitStorage(os.Getenv("RABBIT_MQ_URI"))
	// if err != nil {
	// 	logger.Errorw("Error during connect to RabbitMQ broker", "error", err)
	// }
	closeCh := make(chan struct{})
	w := worker.NewWorker(sqsstor, s3stor, closeCh)
	w.DoWork(logger)
	signal.Notify(sch, os.Interrupt, syscall.SIGTERM)

	sig := <-sch
	logger.Infof("got signal:%v", sig)
	w.Stop()
}
