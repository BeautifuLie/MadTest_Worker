package sqs

import (
	"fmt"
	"program/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type AwsSQS struct {
	sqsSess   *sqs.SQS
	bucketSQS string
	urlQueue  string
}
type SQSMessage struct {
	sqsSess *sqs.SQS
	message *sqs.Message
}

func (m *SQSMessage) GetBody() string {
	return *m.message.Body
}
func (m *SQSMessage) Finalize(success bool) {
	_, err := m.sqsSess.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String("https://sqs.eu-central-1.amazonaws.com/333746971525/JokesQueueSend"),
		ReceiptHandle: m.message.ReceiptHandle,
	})
	if err != nil {
		fmt.Printf("delete error - %v", err)
	}
}
func NewSqsStorage() (*AwsSQS, error) {

	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqsSess := sqs.New(session)
	conn := &AwsSQS{

		sqsSess:   sqsSess,
		bucketSQS: "jokes-sqs-messages",
		urlQueue:  "https://sqs.eu-central-1.amazonaws.com/333746971525/JokesQueueSend",
	}
	return conn, nil
}

func (a *AwsSQS) StartListen(msgCh chan model.Message) {
	for {
		msgResult, err := a.sqsSess.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            &a.urlQueue,
			MaxNumberOfMessages: aws.Int64(10),
			WaitTimeSeconds:     aws.Int64(10),
		})

		if err != nil {
			fmt.Printf("failed to fetch sqs message %v", err)
			continue
		}
		if len(msgResult.Messages) <= 0 {
			fmt.Printf("no messages in queue\n")
			continue
		}

		for _, msg := range msgResult.Messages {
			msgCh <- &SQSMessage{
				sqsSess: a.sqsSess,
				message: msg,
			}
		}

	}

}

func (a *AwsSQS) DeleteMsg(messageHandle string) error {

	_, err := a.sqsSess.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &a.urlQueue,
		ReceiptHandle: &messageHandle,
	})

	if err != nil {
		return err
	}
	return nil
}
