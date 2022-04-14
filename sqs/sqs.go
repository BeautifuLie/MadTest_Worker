package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type AwsSQS struct {
	sqsSess   *sqs.SQS
	bucketSQS string
	urlQueue  string
}

func NewSqsStorage() (*AwsSQS, error) {

	// session, err := session.NewSession(
	// 	&aws.Config{
	// 		Region: &region,
	// 		Credentials: credentials.NewStaticCredentials(
	// 			id,
	// 			secret,
	// 			token,
	// 		),
	// 	})
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqsSess := sqs.New(session)
	conn := &AwsSQS{

		sqsSess:   sqsSess,
		bucketSQS: *aws.String("jokes-sqs-messages"),
		urlQueue:  *aws.String("https://sqs.eu-central-1.amazonaws.com/333746971525/JokesQueueSend"),
	}
	return conn, nil
}

func (a *AwsSQS) StartListen(msgCh chan string) {
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

			msgCh <- *msg.Body
			a.DeleteMsg(*msg.ReceiptHandle)

		}

	}

}

// func (a *AwsSQS) Worker(msgCh chan string, id int) {

// 	for msg := range msgCh {

// 		fmt.Printf("worker %v started a job\n", id)

// 		name := time.Now().Format("2017-09-07 17:06:04.000000000")
// 		// err := UploadMessageTos3(name, strings.NewReader(msg))
// 		s3s, err := NewS3Storage()
// 		if err != nil {
// 			return
// 		}
// 		s3s.UploadMessageTos3(name, strings.NewReader(msg))
// 		// err = a.DeleteMsg(*msg.ReceiptHandle)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	return
// 		// }

// 		fmt.Printf("worker %v finished a job\n", id)

// 	}

// }
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
