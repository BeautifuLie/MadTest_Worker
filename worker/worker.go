package worker

import (
	"io"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Listener interface {
	StartListen(chan string)
}

type Storage interface {
	UploadMessageTos3(name string, data io.Reader) error
}

type Worker struct {
	listener Listener
	storage  Storage
}

func NewWorker(l Listener, st Storage) *Worker {
	w := &Worker{
		listener: l,
		storage:  st,
	}
	return w
}
func (w *Worker) DoWork(logger *zap.SugaredLogger) {
	logger.Info("DoWork starts")
	msgChan := make(chan string, 100)
	go w.listener.StartListen(msgChan)

	for i := 1; i <= 5; i++ {
		go w.processMessage(msgChan, i, logger)
	}
}

func (w *Worker) processMessage(msgs chan string, id int, logger *zap.SugaredLogger) {
	for msg := range msgs {
		logger.Infof("worker %v started a job", id)
		name := time.Now().Format("2017-09-07 17:06:04.000000000")

		err := w.storage.UploadMessageTos3(name, strings.NewReader(msg))
		if err != nil {
			logger.Errorw("Upload message to S3 error", err)
		}
		logger.Infof("worker %v finished a job", id)
	}

}
