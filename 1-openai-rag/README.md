# 1 - OpenAI RAG 

In this chapter we use completely OpenAI to provide a retrieve augmented generation process.
The code is split in two parts: indexing and retrieval. 

### Requirements

* Go 1.25+

### indexing

1. Go to OpenAI dashboard and create a vector storage and a api key.
2. Rename .env.template to .env adding the information required.
3. Run `go run cmd/cli/main.go <path>` to upload the file that is going to be used as source of truth. OpenAI will internally split it chunks and generate the embeddings in the vector database.

### retrieval

1. Run `go run cmd/server/main.go`
2. Api will be running on http://localhost:8080
3. Then you can interact with the endpoint to get answers based on your document.
```sh
curl  -X POST \
  'http://localhost:8080/ask' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "message": "Your question"
}'
```