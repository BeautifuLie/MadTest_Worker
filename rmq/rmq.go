package rmq

import (
	"fmt"
	"program/model"

	"github.com/streadway/amqp"
)

type RabbitFs struct {
	rabbitConn    *amqp.Connection
	rabbitChannel *amqp.Channel
	rabbitQueue   amqp.Queue
}
type RabbitMQMessage struct {
	message amqp.Delivery
}

func (m *RabbitMQMessage) GetBody() string {
	return string(m.message.Body)
}
func (m *RabbitMQMessage) Finalize(success bool) {
	if success {
		m.message.Ack(false)
	} else {
		m.message.Nack(true, false)
	}
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
	err = ch.Qos(1, 0, true)
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

func (r *RabbitFs) StartListen(msgCh chan model.Message) {

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
			msgCh <- &RabbitMQMessage{
				message: msg,
			}

		}

	}()

}
