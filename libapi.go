package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Secrets struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

type Handshake struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

func checkResponse(resp *http.Response, err error) {
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("ERROR: %v\n", err)
		fmt.Println("Status HTTP:", resp.Status)
		os.Exit(1)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func downloadTask() bytes.Buffer {
	//____Handshake____
	secretsFile, err := os.ReadFile("drafts/secrets.json")
	checkError(err)

	var secrets Secrets
	json.Unmarshal(secretsFile, &secrets)

	urlHandshake := secrets.Url + "token/helloapi"
	fmt.Println("Send post to " + urlHandshake)
	respHS, err := http.Post(urlHandshake, "application/json", strings.NewReader("{\"apikey\":\""+secrets.Token+"\"}"))
	checkResponse(respHS, err)

	var hs Handshake
	dec := json.NewDecoder(respHS.Body)
	dec.DisallowUnknownFields()
	dec.Decode(&hs)
	fmt.Println("Code:", hs.Code, "Msg:", hs.Msg, "Token:", hs.Token)

	//____Take_Task____
	urlTask := secrets.Url + "task/" + hs.Token
	fmt.Println("Send GET to " + urlTask)
	respTask, err := http.Get(urlTask)
	checkResponse(respTask, err)

	var buf bytes.Buffer
	io.Copy(&buf, respTask.Body)
	fmt.Println(buf.String())
	return buf
}

func sendAnswer(answer *strings.Reader) {
	//____Answer____
	secretsFile, err := os.ReadFile("drafts/secrets.json")
	checkError(err)

	var secrets Secrets
	json.Unmarshal(secretsFile, &secrets)

	urlAnswer := secrets.Url + "answer/" + secrets.Token
	fmt.Println("Send post to " + urlAnswer)
	respAnswer, err := http.Post(urlAnswer, "application/json", answer)
	checkResponse(respAnswer, err)
	b, err := io.ReadAll(respAnswer.Body)
	fmt.Printf("%s", b)
}
