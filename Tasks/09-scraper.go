package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Scraper() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("scraper")

	//____Solve_Task____
	type Task struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Input    string `json:"input"`
		Question string `json:"question"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", task.Input, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	var textResp *http.Response
	var bytes []byte
	for {
		fmt.Println("Send GET to " + task.Input)
		textResp, err = client.Do(req)
		if textResp.StatusCode == 200 {
			defer textResp.Body.Close()
			bytes, err = io.ReadAll(textResp.Body)
			APIs.CheckError(err)
			break
		} else {
			fmt.Printf("ERROR: %v\n", err)
			fmt.Println("Status HTTP:", textResp.Status)
			b, _ := io.ReadAll(textResp.Body)
			fmt.Printf("%s", b)
			time.Sleep(3 * time.Second)
		}

	}
	APIs.CheckResponse(textResp, err)

	fmt.Println(string(bytes))
	APIs.CheckError(err)

	fmt.Println(task)

	var messages []APIs.Message
	messages = append(messages, APIs.Message{
		Role:    "system",
		Content: string(bytes) + "/\n" + "###" + "/\n" + task.Msg,
	})
	messages = append(messages, APIs.Message{
		Role:    "user",
		Content: task.Question,
	})

	var max_tokens int64 = 100
	response := APIs.Completions(secrets.OpenaiAPIKey, "gpt-3.5-turbo-0125", messages, max_tokens, nil)
	result := response.Choices[0].Message.Content
	fmt.Println(response)

	postBody, _ := json.Marshal(map[string]string{
		"answer": result,
	})
	APIs.SendAnswer(taskToken, postBody, secrets)
}
