package APIs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Secrets struct {
	Url          string `json:"url"`
	Token        string `json:"token"`
	OpenaiAPIKey string `json:"openaiAPIKey"`
	MyRestApi    string `json:"myRestApi,omitempty"`
}

type Handshake struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

type Moderation struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Results []Result `json:"results"`
}

type Result struct {
	Flagged        bool           `json:"flagged"`
	Categories     Categories     `json:"categories"`
	CategoryScores CategoryScores `json:"category_scores"`
}

type Categories struct {
	Sexual           bool `json:"sexual"`
	Hate             bool `json:"hate"`
	Harassment       bool `json:"harassment"`
	SelfHarm         bool `json:"self-harm"`
	SexualMinors     bool `json:"sexual/minors"`
	HateThreatening  bool `json:"hate/threatening"`
	ViolenceGraphic  bool `json:"violence/graphic"`
	SelfHarmIntent   bool `json:"self-harm/intent"`
	SelfHarmInstr    bool `json:"self-harm/instructions"`
	HarassmentThreat bool `json:"harassment/threatening"`
	Violence         bool `json:"violence"`
}

type CategoryScores struct {
	Sexual           float64 `json:"sexual"`
	Hate             float64 `json:"hate"`
	Harassment       float64 `json:"harassment"`
	SelfHarm         float64 `json:"self-harm"`
	SexualMinors     float64 `json:"sexual/minors"`
	HateThreatening  float64 `json:"hate/threatening"`
	ViolenceGraphic  float64 `json:"violence/graphic"`
	SelfHarmIntent   float64 `json:"self-harm/intent"`
	SelfHarmInstr    float64 `json:"self-harm/instructions"`
	HarassmentThreat float64 `json:"harassment/threatening"`
	Violence         float64 `json:"violence"`
}

func CheckResponse(resp *http.Response, err error) {
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("ERROR: %v\n", err)
		fmt.Println("Status HTTP:", resp.Status)
		b, _ := io.ReadAll(resp.Body)
		fmt.Printf("%s", b)
		os.Exit(1)
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func DownloadTask(task string) (string, bytes.Buffer, Secrets) {
	//____Handshake____
	secretsFile, err := os.ReadFile("drafts/secrets.json")
	CheckError(err)

	var secrets Secrets
	json.Unmarshal(secretsFile, &secrets)

	urlHandshake := secrets.Url + "token/" + task
	fmt.Println("Send post to " + urlHandshake)
	respHS, err := http.Post(urlHandshake, "application/json", strings.NewReader("{\"apikey\":\""+secrets.Token+"\"}"))
	CheckResponse(respHS, err)

	var hs Handshake
	dec := json.NewDecoder(respHS.Body)
	dec.DisallowUnknownFields()
	dec.Decode(&hs)

	//____Take_Task____
	urlTask := secrets.Url + "task/" + hs.Token
	fmt.Println("Send GET to " + urlTask)
	respTask, err := http.Get(urlTask)
	CheckResponse(respTask, err)

	var buf bytes.Buffer
	io.Copy(&buf, respTask.Body)
	//fmt.Println(buf.String())
	return hs.Token, buf, secrets
}

func SendAnswer(taskToken string, answer []byte, secrets Secrets) {
	//____Answer____
	urlAnswer := secrets.Url + "answer/" + taskToken
	fmt.Println("Send post to " + urlAnswer)
	respAnswer, err := http.Post(urlAnswer, "application/json", bytes.NewReader(answer))
	CheckResponse(respAnswer, err)
	b, _ := io.ReadAll(respAnswer.Body)
	fmt.Printf("%s", b)
}

func GetUrl(text string) []string {
	re := regexp.MustCompile(`\b(?:https?://|www\.)\S+\b`)
	matches := re.FindAllString(text, -1)
	return matches
}
