package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sashabaranov/go-openai"
)

func Ownapipro() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("ownapipro")

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

	ctx := context.Background()
	api := APIs.ClinetOpenAI{
		Client: openai.NewClient(secrets.OpenaiAPIKey),
	}
	instruction := "I'm assistant for remembering facts from user and providing him answers."
	assistant, err := api.Client.CreateAssistant(ctx, openai.AssistantRequest{
		Model:        openai.GPT3Dot5Turbo,
		Instructions: &instruction,
	})
	APIs.CheckError(err)

	thread, err := api.Client.CreateThread(ctx, openai.ThreadRequest{})
	APIs.CheckError(err)
	http.HandleFunc("/answer", questionHandlerPro(&api, &thread))
	go func() {
		if err := http.ListenAndServeTLS(":16123", "drafts/cert_key/cert.pem", "drafts/cert_key/key.pem", nil); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	fmt.Println("Server started")
	postBody, _ := json.Marshal(
		map[string]string{
			"answer": secrets.MyRestApi + "/answer"},
	)

	APIs.SendAnswer(taskToken, postBody, secrets)

	api.Client.DeleteThread(ctx, thread.ID)
	api.Client.DeleteAssistant(ctx, assistant.ID)
}

func questionHandlerPro(api *APIs.ClinetOpenAI, thread *openai.Thread) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check POST req
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var q Question
		if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		///////Add question to thread
		ctx := context.Background()
		newMessage := openai.MessageRequest{
			Role:    "user",
			Content: q.Question,
		}
		_, err := api.Client.CreateMessage(ctx, thread.ID, newMessage)
		APIs.CheckError(err)

		assistants, err := api.Client.ListAssistants(ctx, nil, nil, nil, nil)
		APIs.CheckError(err)
		run, err := api.Client.CreateRun(ctx, thread.ID,
			openai.RunRequest{
				AssistantID: *assistants.FirstID,
			},
		)
		APIs.CheckError(err)

		// Poll for a status that indicates run has finished
		for run.Status != openai.RunStatusCompleted {
			run, err = api.Client.RetrieveRun(ctx, run.ThreadID, run.ID)
			if err != nil {
				fmt.Printf("RetrieveRun error: %v\n", err)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}

		APIs.CheckError(err)
		order := "desc"
		allMess, err := api.Client.ListMessage(ctx, thread.ID, nil, &order, nil, nil)
		APIs.CheckError(err)

		answer := allMess.Messages[0].Content[0].Text.Value
		reply := "{\"reply\":\"" + answer + "\"}"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(reply))
	}
}
