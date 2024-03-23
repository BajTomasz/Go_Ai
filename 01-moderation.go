package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func moderation(openaiAPIKey string, input string) bool {
	//start := time.Now()
	jsonInput := fmt.Sprintf(`{"input": "%s"}`, input)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/moderations", strings.NewReader(jsonInput))
	checkError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	client := http.Client{}
	respModer, err := client.Do(req)
	defer respModer.Body.Close()
	checkResponse(respModer, err)

	var moderation Moderation
	json.NewDecoder(respModer.Body).Decode(&moderation)
	//fmt.Println(moderation)
	//fmt.Println("%.2fs", time.Since(start).Seconds())
	return moderation.Results[0].Flagged
}

func main() {
	var resp bytes.Buffer
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)
	fmt.Println("Code:", task.Code)
	fmt.Println("Msg:", task.Msg)
	fmt.Println("Input:", task.Input)

	var results []int
	for i := range task.Input[:] {
		if moderation(secrets.OpenaiAPIKey, task.Input[i]) {
			results = append(results, 1)
		} else {
			results = append(results, 0)
		}
	}

	postBody, _ := json.Marshal(map[string][]int{
		"answer": results,
	})

	sendAnswer(taskToken, postBody, secrets)
}
