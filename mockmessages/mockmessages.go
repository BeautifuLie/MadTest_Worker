package mockmessages

import (
	"encoding/json"
	"fmt"
	"program/model"
)

func MockMessages(mockCh chan string) chan string {
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
		mockCh <- string(out)
	}
	return mockCh

}
