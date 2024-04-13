package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

func findName(input string) string {
	words := strings.Fields(input)
	if len(words) == 0 {
		return ""
	}

	lastWord := strings.TrimRightFunc(words[len(words)-1], func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	return lastWord

}

func Inprompt() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("inprompt")

	//____Solve_Task____
	type Task struct {
		Code     int      `json:"code"`
		Msg      string   `json:"msg"`
		Input    []string `json:"input"`
		Question string   `json:"question"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)

	name := findName(task.Question)
	knowlade := ""

	for _, input := range task.Input {
		if strings.Contains(input, name) {
			knowlade = input
		}
	}

	var messages []APIs.Message
	systemPrompt := "Just answer the questions based on this information:/\n###/\n"
	systemPrompt += knowlade
	messages = append(messages, APIs.Message{
		Role:    "system",
		Content: systemPrompt,
	})
	messages = append(messages, APIs.Message{
		Role:    "user",
		Content: task.Question,
	})

	response := APIs.Completions(secrets.OpenaiAPIKey, "gpt-3.5-turbo-0125", messages, 0, nil)
	result := response.Choices[0].Message.Content
	fmt.Println(messages)
	fmt.Println(response)

	postBody, _ := json.Marshal(map[string]string{
		"answer": result,
	})
	APIs.SendAnswer(taskToken, postBody, secrets)
}
