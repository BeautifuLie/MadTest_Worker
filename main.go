package main

import (
	"fmt"
	"log"
	"os"
	"program/storage"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
	_ "gocloud.dev/blob/s3blob"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error during load environments", err)
	}
	awsstor, err := storage.NewAwsStorage(
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"")
	if err != nil {
		log.Fatal("Error during connect to AWS services", "error", err)
	}
	msgCh := make(chan *sqs.Message, 100)

	go awsstor.GetMsg(msgCh)

	for i := 1; i <= 5; i++ {
		go awsstor.Worker(msgCh, i)
	}

	c := make(chan struct{})
	<-c
}
