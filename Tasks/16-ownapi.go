package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

func Ownapi() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("ownapi")
	fmt.Println(taskToken, secrets)

	//____Solve_Task____

	type Task struct {
		Code  int    `json:"code"`
		Msg   string `json:"msg"`
		Hint1 string `json:"hint1"`
		Hint2 string `json:"hint2"`
		Hint3 string `json:"hint3"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)

	http.HandleFunc("/question", questionHandler(secrets.OpenaiAPIKey))

	go func() {
		if err := http.ListenAndServe(":16123", nil); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	fmt.Println("Server started on port 16123")
	//
	fmt.Println(secrets.MyRestApi + "/question")
	postBody, _ := json.Marshal(
		map[string]string{
			"answer": secrets.MyRestApi + "/question"},
	)

	APIs.SendAnswer(taskToken, postBody, secrets)
}

type Question struct {
	Question string `json:"question"`
}

func questionHandler(openaiAPIKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var q Question
		if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		fmt.Println(q.Question)

		api := APIs.ClinetOpenAI{
			Client: openai.NewClient(openaiAPIKey),
		}
		ctx := context.Background()
		req := openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{{
				Role:    "user",
				Content: q.Question,
			}},
			MaxTokens: 100,
		}
		resp, err := api.Client.CreateChatCompletion(ctx, req)
		APIs.CheckError(err)

		reply := "{\"reply\":\"" + resp.Choices[0].Message.Content + "\"}"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(reply))
	}
}
