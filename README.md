# Go_Ai
docker run -p 6333:6333 -p 6334:6334 -d -e QDRANT__SERVICE__GRPC_PORT="6334" -v $(pwd)/qdrant_storage:/qdrant/storage:z qdrant/qdrant

go run . task-number