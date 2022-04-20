package mockmessages

import (
	"encoding/json"
	"fmt"
	"program/model"
)

type MockMessage struct {
	message string
}

func (m *MockMessage) GetBody() string {
	return m.message
}
func (m *MockMessage) Finalize(success bool) {
	if success {
		fmt.Println("1")
	} else {
		fmt.Println("2")
	}
}
func MockMessages(mockCh chan model.Message) chan model.Message {
	joke1 := model.Joke{
		Title: "title",
		Body:  "body",
		Score: 1,
		ID:    "123321",
	}
	joke2 := model.Joke{
		Title: "title",
		Body:  "body",
		Score: 1,
		ID:    "123322",
	}
	joke3 := model.Joke{
		Title: "title",
		Body:  "body",
		Score: 1,
		ID:    "123323",
	}
	jokes := []model.Joke{joke1, joke2, joke3}

	for _, msg := range jokes {
		out, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
		}
		mockCh <- &MockMessage{
			message: string(out),
		}
	}
	return mockCh

}
