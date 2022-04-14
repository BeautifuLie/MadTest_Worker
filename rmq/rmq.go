package rmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitFs struct {
	rabbitConn    *amqp.Connection
	rabbitChannel *amqp.Channel
	rabbitQueue   amqp.Queue
}

func NewRabbitStorage(url string) (*RabbitFs, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	// defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// defer ch.Close()

	q, err := ch.QueueDeclare(
		"jokes-message",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	// err = ch.Qos(1, 0, true)
	if err != nil {
		return nil, err
	}
	c := &RabbitFs{
		rabbitConn:    conn,
		rabbitChannel: ch,
		rabbitQueue:   q,
	}
	return c, nil
}

func (r *RabbitFs) StartListen(msgCh chan string) {

	msgs, err := r.rabbitChannel.Consume(
		r.rabbitQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Printf("failed to fetch  message %v", err)
	}
	go func() {
		for msg := range msgs {
			msgCh <- string(msg.Body)
			msg.Ack(true)
		}
	}()

}

// func MsgsToString(msgs interface{}) chan string {

// }

// func (r *RabbitFs) Worker(msgs chan string, id int) {

// 	for msg := range msgs {

// 		fmt.Printf("worker %v started a job\n", id)

// 		name := time.Now().Format("2017-09-07 17:06:04.000000000")
// 		s3s, err := NewS3Storage()
// 		if err != nil {
// 			return
// 		}
// 		s3s.UploadMessageTos3(name, strings.NewReader(msg))
// 		if err != nil {
// 			return
// 		}

// 		// err = msg.Ack(false)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	return
// 		// }
// 		fmt.Printf("worker %v finished a job\n", id)

// 	}

// }
