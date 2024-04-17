package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sashabaranov/go-openai"
)

func Knowledge() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("knowledge")

	//____Solve_Task____
	type Task struct {
		Code      int    `json:"code"`
		Msg       string `json:"msg"`
		Question  string `json:"question"`
		Database1 string `json:"database #1"`
		Database2 string `json:"database #2"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)

	fmt.Println(task.Msg)
	fmt.Println(task.Question)

	//////////////////

	api := APIs.ClinetOpenAI{
		Client: openai.NewClient(secrets.OpenaiAPIKey),
	}
	ID1 := api.FindFunction("getExchangeRate")
	ID2 := api.FindFunction("getCountryPopulation")
	ID3 := api.FindFunction("generalKnowledge")

	var assistant openai.Assistant
	if ID1 == nil && ID1 == ID2 && ID2 == ID3 {
		exchangeRateParams := APIs.Parameters{
			Type: "object",
			Properties: map[string]APIs.Property{
				"currency": {
					Type:        "string",
					Description: "The currency code like EUR GBP PLN USD",
				},
			},
		}

		getExchangeRateFunc := openai.FunctionDefinition{
			Name:        "getExchangeRate",
			Description: "Get the exchange rate for a specific currency",
			Parameters:  exchangeRateParams,
		}

		countryPopulationParams := APIs.Parameters{
			Type: "object",
			Properties: map[string]APIs.Property{
				"country": {
					Type:        "string",
					Description: "The English name of the country",
				},
			},
		}

		getCountryPopulationFunc := openai.FunctionDefinition{
			Name:        "getCountryPopulation",
			Description: "Get the population of a specific country",
			Parameters:  countryPopulationParams,
		}

		generalKnowledge := openai.FunctionDefinition{
			Name:        "generalKnowledge",
			Description: "General Knowledge",
			Parameters:  nil,
		}

		listFunction := []*openai.FunctionDefinition{}
		listFunction = append(listFunction, &getExchangeRateFunc)
		listFunction = append(listFunction, &getCountryPopulationFunc)
		listFunction = append(listFunction, &generalKnowledge)
		assistant = api.CreateAssistant(openai.GPT3Dot5Turbo0125, "task13", task.Msg, listFunction)

	} else {
		listAssistants := api.ListAssistants()
		for i := range listAssistants {
			if listAssistants[i].ID == *ID1 {
				assistant = listAssistants[i]
				break
			}
		}
	}

	threadReq := openai.ThreadRequest{
		Messages: []openai.ThreadMessage{{
			Role:    openai.ThreadMessageRoleUser,
			Content: task.Question,
		}},
	}
	ctx := context.Background()
	thread, err := api.Client.CreateThread(ctx, threadReq)
	APIs.CheckError(err)
	run, err := api.Client.CreateRun(ctx, thread.ID,
		openai.RunRequest{
			AssistantID:  assistant.ID,
			Instructions: *assistant.Instructions,
		},
	)
	APIs.CheckError(err)

	// Poll for a status that indicates run has finished
	for run.Status != openai.RunStatusRequiresAction {
		run, err = api.Client.RetrieveRun(ctx, run.ThreadID, run.ID)
		if err != nil {
			fmt.Printf("RetrieveRun error: %v\n", err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	selectedFunction := run.RequiredAction.SubmitToolOutputs.ToolCalls[0].Function.Name
	arguments := run.RequiredAction.SubmitToolOutputs.ToolCalls[0].Function.Arguments

	var data map[string]interface{}
	err = json.Unmarshal([]byte(arguments), &data)
	APIs.CheckError(err)
	var postBody []byte

	switch selectedFunction {
	case "getExchangeRate":
		currency := getExchangeRate(fmt.Sprintf("%v", data["currency"]))
		postBody, _ = json.Marshal(
			map[string]float64{
				"answer": currency})
	case "getCountryPopulation":
		country := getCountryPopulation(fmt.Sprintf("%v", data["country"]))
		postBody, _ = json.Marshal(
			map[string]int{
				"answer": country})
	default:
		req := openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0125,
			Messages: []openai.ChatCompletionMessage{{
				Role:    "user",
				Content: task.Question,
			}},
		}
		respGeneral, err := api.Client.CreateChatCompletion(ctx, req)
		APIs.CheckError(err)
		answear := respGeneral.Choices[0].Message.Content
		postBody, _ = json.Marshal(
			map[string]string{
				"answer": answear})
	}
	APIs.SendAnswer(taskToken, postBody, secrets)
}

type Rates struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Mid           float64 `json:"mid"`
}

type NBPResponse struct {
	Table    string  `json:"table"`
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Rates    []Rates `json:"rates"`
}

func getExchangeRate(currency string) float64 {
	url := fmt.Sprintf("https://api.nbp.pl/api/exchangerates/rates/A/%v", currency)

	req, err := http.NewRequest("GET", url, nil)
	APIs.CheckError(err)
	req.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	resp, err := http.DefaultClient.Do(req)
	APIs.CheckResponse(resp, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	APIs.CheckError(err)

	var result NBPResponse
	json.Unmarshal(bytes, &result)

	return result.Rates[0].Mid
}

type Population struct {
	Population int `json:"population"`
}

func getCountryPopulation(country string) int {
	url := fmt.Sprintf("https://restcountries.com/v3.1/name/%v?fields=population", country)
	req, err := http.NewRequest("GET", url, nil)
	APIs.CheckError(err)

	resp, err := http.DefaultClient.Do(req)
	APIs.CheckResponse(resp, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	APIs.CheckError(err)

	var ret Population
	json.Unmarshal(bytes, &ret)

	return ret.Population
}
