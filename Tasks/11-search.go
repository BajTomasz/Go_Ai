package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	qdrant "github.com/qdrant/go-client/qdrant"
)

type Source struct {
	Title     string `json:"title"`
	Url       string `json:"url"`
	Info      string `json:"info"`
	Date      string `json:"date"`
	Embedding []float32
}

func Search() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("search")

	//____Solve_Task____
	type Task struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Question string `json:"question"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)

	fmt.Println(task.Msg)
	fmt.Println(task.Question)

	// Contact the server
	collectionName := "collection_search"
	qdrantAPI, _ := APIs.NewQdrantAPI(collectionName)
	collectionExists := qdrantAPI.FindCollection(collectionName)
	if !collectionExists {
		qdrantAPI.CreateCollection()
	}

	// Get Information
	sourceUrl := APIs.GetUrl(task.Msg)[0]
	sources := getSourcesList(sourceUrl)

	// Make embedding
	var titles []string
	for _, source := range sources {
		titles = append(titles, source.Title)
	}

	embeddingObjects := APIs.Embeddings(secrets.OpenaiAPIKey, "text-embedding-3-small", titles).Data

	for i := range sources {
		sources[i].Embedding = embeddingObjects[i].Embedding
	}

	// Upsert embedding in to qdrant
	var upsertPoints []*qdrant.PointStruct
	for index, source := range sources {
		upsertPoints = append(upsertPoints, &qdrant.PointStruct{
			Id:      &qdrant.PointId{PointIdOptions: &qdrant.PointId_Num{Num: uint64(index)}},
			Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: source.Embedding}}},
			Payload: map[string]*qdrant.Value{
				"url": {Kind: &qdrant.Value_StringValue{StringValue: source.Url}},
			},
		})
	}

	qdrantAPI.Upsert(upsertPoints)

	// Embedding question
	embeddedQuestion := APIs.Embeddings(secrets.OpenaiAPIKey, "text-embedding-3-small", []string{task.Question}).Data[0].Embedding

	// Find source
	answear := qdrantAPI.SearchClosestVector(embeddedQuestion)
	fmt.Println(answear.Score)
	answearSource := sources[answear.GetId().GetNum()]

	postBody, _ := json.Marshal(map[string]string{
		"answer": answearSource.Url,
	})
	APIs.SendAnswer(taskToken, postBody, secrets)
}

func getSourcesList(url string) []Source {
	response, err := http.Get(url)
	APIs.CheckError(err)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	APIs.CheckError(err)

	articles := []Source{}
	err = json.Unmarshal(body, &articles)
	APIs.CheckError(err)

	return articles
}
