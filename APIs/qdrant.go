package APIs

import (
	"context"
	"time"

	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QdrantAPI struct {
	CollectionsClient qdrant.CollectionsClient
	CollectionName    string
	Connection        *grpc.ClientConn
}

func NewQdrantAPI(collectionName string) (*QdrantAPI, error) {
	addr := "localhost:6334"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	CheckError(err)
	collectionsClient := qdrant.NewCollectionsClient(conn)

	//ctx, _ := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	return &QdrantAPI{
		CollectionsClient: collectionsClient,
		CollectionName:    collectionName,
		Connection:        conn,
	}, nil
}

func (api *QdrantAPI) CreateCollection() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := api.CollectionsClient.Create(ctx, &qdrant.CreateCollection{
		CollectionName: api.CollectionName,
		VectorsConfig: &qdrant.VectorsConfig{
			Config: &qdrant.VectorsConfig_Params{
				Params: &qdrant.VectorParams{
					Size:     1536,
					Distance: qdrant.Distance_Cosine,
				},
			},
		},
	})

	CheckError(err)
}

func (api *QdrantAPI) DeleteCollection() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	api.CollectionsClient.Delete(ctx, &qdrant.DeleteCollection{CollectionName: api.CollectionName})
}

func (api *QdrantAPI) ListCollection() []string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	listCollections, err := api.CollectionsClient.List(ctx, &qdrant.ListCollectionsRequest{})
	CheckError(err)
	var listCollectionsName []string
	for _, collectionDescription := range listCollections.GetCollections() {
		listCollectionsName = append(listCollectionsName, collectionDescription.GetName())
	}
	return listCollectionsName
}

func (api *QdrantAPI) FindCollection(collectionName string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	collectionExists, err := api.CollectionsClient.CollectionExists(ctx, &qdrant.CollectionExistsRequest{
		CollectionName: collectionName,
	})
	CheckError(err)
	return collectionExists.Result.Exists
}

func (api *QdrantAPI) Upsert(upsertPoints []*qdrant.PointStruct) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	waitUpsert := true
	pointsClient := qdrant.NewPointsClient(api.Connection)

	_, err := pointsClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: api.CollectionName,
		Wait:           &waitUpsert,
		Points:         upsertPoints,
	})
	CheckError(err)
}

func (api *QdrantAPI) SearchClosestVector(queryVector []float32) *qdrant.ScoredPoint {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pointsClient := qdrant.NewPointsClient(api.Connection)
	searchRequest := &qdrant.SearchPoints{
		CollectionName: api.CollectionName,
		Vector:         queryVector,
		Limit:          1,
	}
	searchResponse, err := pointsClient.Search(ctx, searchRequest)
	CheckError(err)

	return searchResponse.Result[0]
}
