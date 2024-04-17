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

type PeopleData struct {
	Name        string `json:"imie"`
	Surname     string `json:"nazwisko"`
	Age         int    `json:"wiek"`
	Description string `json:"o_mnie"`
	FavHero     string `json:"ulubiona_postac_z_kapitana_bomby"`
	FavSeries   string `json:"ulubiony_serial"`
	FavMovie    string `json:"ulubiony_film"`
	FavCoulor   string `json:"ulubiony_kolor"`
	Embedding   []float32
}

func People() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("people")

	//____Solve_Task____
	type Task struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Data     string `json:"data"`
		Question string `json:"question"`
		Hint1    string `json:"hint1"`
		Hint2    string `json:"hint2"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)

	fmt.Println(task.Msg)
	fmt.Println(task.Question)

	// Contact the server
	collectionName := "collection_people"
	qdrantAPI, _ := APIs.NewQdrantAPI(collectionName)
	collectionExists := qdrantAPI.FindCollection(collectionName)

	if !collectionExists {
		qdrantAPI.CreateCollection()
	}

	// Get Information
	sources := getPeopleList(task.Data)

	// Make embedding
	var information []string

	for _, source := range sources {
		txt := source.Name + " " + source.Surname + " " + source.Description
		information = append(information, txt)
	}

	embeddingObjects := APIs.Embeddings(secrets.OpenaiAPIKey, "text-embedding-3-small", information).Data

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
				"info": {Kind: &qdrant.Value_StringValue{StringValue: source.Description}},
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
		"answer": answearSource.Description,
	})

	APIs.SendAnswer(taskToken, postBody, secrets)

}

func getPeopleList(url string) []PeopleData {
	response, err := http.Get(url)
	APIs.CheckError(err)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	APIs.CheckError(err)

	articles := []PeopleData{}
	err = json.Unmarshal(body, &articles)
	APIs.CheckError(err)

	return articles
}
