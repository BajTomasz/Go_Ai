package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func Gnome() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("gnome")
	fmt.Println(taskToken, secrets)

	//____Solve_Task____

	type Task struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Hint string `json:"hint"`
		Url  string `json:"url"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)

	fmt.Println(task.Msg)
	fmt.Println(task.Hint)
	fmt.Println(task.Url)

	//////////////////

	api := APIs.ClinetOpenAI{
		Client: openai.NewClient(secrets.OpenaiAPIKey),
	}

	mess := []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: task.Msg + " " + task.Hint},
		{
			Role: "user",
			MultiContent: []openai.ChatMessagePart{{
				Type:     "image_url",
				ImageURL: &openai.ChatMessageImageURL{URL: task.Url},
			}}},
	}

	ctx := context.Background()
	reqImage := openai.ChatCompletionRequest{
		Model:    openai.GPT4Turbo,
		Messages: mess,
	}

	respImage, err := api.Client.CreateChatCompletion(ctx, reqImage)
	APIs.CheckError(err)

	fmt.Printf("%v", respImage.Choices[0].Message)

	postBody, _ := json.Marshal(
		map[string]string{
			"answer": respImage.Choices[0].Message.Content})

	APIs.SendAnswer(taskToken, postBody, secrets)

}
