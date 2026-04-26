package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
)

type askRequest struct {
	Message string `json:"message"`
}

type askResponse struct {
	Answer string `json:"answer"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	apiKey := os.Getenv("OPENAI_KEY")
	vectorStoreID := os.Getenv("OPENAI_VECTOR_DB_ID")

	if apiKey == "" || vectorStoreID == "" {
		log.Fatal("OPENAI_KEY and OPENAI_VECTOR_DB_ID must be set")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	http.HandleFunc("POST /ask", func(w http.ResponseWriter, r *http.Request) {
		var req askRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Message == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := client.Responses.New(context.Background(), responses.ResponseNewParams{
			Model: openai.ChatModelGPT4oMini,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: param.NewOpt(req.Message),
			},
			Tools: []responses.ToolUnionParam{
				responses.ToolParamOfFileSearch([]string{vectorStoreID}),
			},
		})
		if err != nil {
			log.Printf("openai error: %v", err)
			http.Error(w, "failed to query the model", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(askResponse{Answer: resp.OutputText()})
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
