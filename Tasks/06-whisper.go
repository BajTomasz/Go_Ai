package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

func Whisper() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("whisper")

	//____Solve_Task____
	type Task struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Hint string `json:"hint"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)

	re := regexp.MustCompile(`https?://\S+`)
	voiceUrl := re.FindString(task.Msg)

	voiceResp, err := http.Get(voiceUrl)
	APIs.CheckResponse(voiceResp, err)
	defer voiceResp.Body.Close()
	bytes, _ := io.ReadAll(voiceResp.Body)

	response := APIs.Transcriptions(secrets.OpenaiAPIKey, "whisper-1", bytes)
	result := response.Teskt

	postBody, _ := json.Marshal(map[string]string{
		"answer": result,
	})
	APIs.SendAnswer(taskToken, postBody, secrets)
}
