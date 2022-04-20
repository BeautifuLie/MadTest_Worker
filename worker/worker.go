package worker

import (
	"io"
	"program/model"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Listener interface {
	StartListen(chan model.Message)
}

type Storage interface {
	UploadMessageTos3(name string, data io.Reader) error
}

type Worker struct {
	listener Listener
	storage  Storage
	wg       sync.WaitGroup
	close    chan struct{}
}

func NewWorker(l Listener, st Storage, closeCh chan struct{}) *Worker {
	w := &Worker{
		listener: l,
		storage:  st,
		close:    closeCh,
	}
	return w
}
func (w *Worker) DoWork(logger *zap.SugaredLogger) {
	logger.Info("DoWork starts")
	msgChan := make(chan model.Message, 100)
	go w.listener.StartListen(msgChan)
	w.wg.Add(5)
	for i := 1; i <= 5; i++ {
		go func(i int) {
			defer w.wg.Done()
			w.processMessage(msgChan, i, logger)
		}(i)
	}
}

func (w *Worker) processMessage(msgs chan model.Message, id int, logger *zap.SugaredLogger) {
	for {
		var msg model.Message

		select {
		case msg = <-msgs:
		case <-w.close:
			return
		}
		logger.Infof("worker %v started a job", id)
		msgBody := msg.GetBody()
		name := time.Now().Format("2017-09-07 17:06:04.000000000")
		err := w.storage.UploadMessageTos3(name, strings.NewReader(msgBody))
		if err != nil {
			logger.Errorw("Upload message to S3 error", err)
		}

		msg.Finalize(err == nil)

		logger.Infof("worker %v finished a job", id)

	}

}
func (w *Worker) Stop() {
	close(w.close)
	w.wg.Wait()
}
