package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
)

type AppInput struct {
	Tool string `json:"tool"`
	Desc string `json:"desc"`
	Date string `json:"date,omitempty"`
}

func Tools() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("tools")

	//____Solve_Task____

	type Task struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Hint     string `json:"hint"`
		Todo     string `json:"example for ToDo"`
		Calendar string `json:"example for Calendar"`
		Question string `json:"question"`
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
	ID1 := api.FindFunction("ToDo")
	ID2 := api.FindFunction("Calendar")

	var assistant openai.Assistant
	if ID1 == nil && ID1 == ID2 {
		toDoParams := APIs.Parameters{
			Type: "object",
			Properties: map[string]APIs.Property{
				"desc": {
					Type:        "string",
					Description: "Task that needs to be added to todo list",
				},
			},
		}

		addToList := openai.FunctionDefinition{
			Name:        "ToDo",
			Description: "Add task to todo list",
			Parameters:  toDoParams,
		}

		calendarParams := APIs.Parameters{
			Type: "object",
			Properties: map[string]APIs.Property{
				"desc": {
					Type:        "string",
					Description: "Event that needs to be added to calendar",
				},
				"date": {
					Type:        "string",
					Description: task.Hint,
				},
			},
		}

		addToCalendar := openai.FunctionDefinition{
			Name:        "Calendar",
			Description: "Add event to calendar",
			Parameters:  calendarParams,
		}

		listFunction := []*openai.FunctionDefinition{}
		listFunction = append(listFunction, &addToList)
		listFunction = append(listFunction, &addToCalendar)
		assistant = api.CreateAssistant(openai.GPT3Dot5Turbo0125, "task14", task.Msg+" "+"Today is: "+time.Now().Format("2006-01-02"), listFunction)
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

	var response map[string]string
	err = json.Unmarshal([]byte(arguments), &response)
	APIs.CheckError(err)

	data := AppInput{
		Tool: selectedFunction,
		Desc: response["desc"],
		Date: response["date"],
	}
	postBody, _ := json.Marshal(
		map[string]AppInput{
			"answer": data})

	APIs.SendAnswer(taskToken, postBody, secrets)
}
