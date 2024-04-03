package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func whisper() {
	var resp bytes.Buffer
	taskToken, resp, secrets := downloadTask("whisper")
	fmt.Println(taskToken, resp.String(), secrets)

	//____Solve_Task____

	type Task struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Hint string `json:"hint"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)

	re := regexp.MustCompile(`https?://\S+`)
	voiceUrl := re.FindString(task.Msg)

	voiceResp, err := http.Get(voiceUrl)
	checkResponse(voiceResp, err)
	defer voiceResp.Body.Close()
	bytes, err := io.ReadAll(voiceResp.Body)

	response := transcriptions(secrets.OpenaiAPIKey, "whisper-1", bytes)
	result := response.Teskt

	postBody, _ := json.Marshal(map[string]string{
		"answer": result,
	})

	sendAnswer(taskToken, postBody, secrets)
}
