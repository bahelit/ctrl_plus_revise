package chat

import (
	"fmt"
	"log/slog"
	"strconv"
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

func IntSliceToString(nums []int) string {
	return fmt.Sprint(nums)
}

func StringToIntSlice(s string) ([]int, error) {
	nums := make([]int, 0, len(stringSplitter(s)))
	for _, v := range stringSplitter(s) {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}
	return nums, nil
}

func stringSplitter(s string) []string {
	var res []string
	for _, v := range s {
		if v == ' ' {
			continue
		}
		res = append(res, string(v))
	}
	return res
}

func (c *Chat) ContextToDB() []byte {
	stringContext := IntSliceToString(c.Context)
	return []byte(stringContext[1 : len(stringContext)-1])
}

func (c *Chat) ContextFromDB(db []byte) {
	ints, err := StringToIntSlice(string(db))
	if err != nil {
		slog.Error("Failed to parse chat's context", "error", err.Error())
		return
	}
	c.Context = ints
}
