package chat

import (
	"bytes"
	"encoding/binary"
	"log/slog"
	"strings"
)

type Chat struct {
	ID        *int64   `json:"id"`
	Model     int      `json:"model"`
	Context   []int    `json:"context"`
	Owner     string   `json:"owner"`
	Title     string   `json:"title"`
	Questions []string `json:"questions"`
	Responses []string `json:"responses"`
}

const (
	Separator = "(╯°o°）╯︵ ┻━┻"
)

func (c *Chat) QuestionsToString() string {
	return strings.Join(c.Questions, Separator)
}

func (c *Chat) ResponsesToString() string {
	return strings.Join(c.Responses, Separator)
}
func (c *Chat) SetQuestions(questions string) {
	c.Questions = strings.Split(questions, Separator)
}

func (c *Chat) SetResponses(responses string) {
	c.Responses = strings.Split(responses, Separator)
}

func (c *Chat) ContextToDB() []byte {
	result := make([]byte, len(c.Context)*4)
	for i, v := range c.Context {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(v))
		copy(result[i*4:], b)
	}
	return result
}

func (c *Chat) ContextFromDB(db []byte) {
	buf := bytes.NewBuffer(db)

	var intArray []int32

	err := binary.Read(buf, binary.LittleEndian, &intArray)
	if err != nil {
		slog.Error("Failed to read from db", "error", err)
		return
	}

	c.Context = int32ToInt32Array(intArray)
}

func int32ToInt32Array(arr []int32) []int {
	result := make([]int, len(arr))
	for i, v := range arr {
		result[i] = int(v)
	}
	return result
}
