package backend

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

//go:embed messages.json
var messagesJSON []byte

const indexNameMessageData = "message_data"

type esService struct {
	client *elasticsearch.Client
}

func NewElasticsearch() *esService {
	es, err := makeElasticsearchClient()
	if err != nil {
		panic(fmt.Errorf("cannot connect to elasticsearch: %w", err))
	}

	s := &esService{
		client: es,
	}

	err = s.initElasticsearch()
	if err != nil {
		panic(fmt.Errorf("init elasticsearch: %w", err))
	}

	return s
}

func makeElasticsearchClient() (*elasticsearch.Client, error) {
	url := os.Getenv("ELASTICSEARCH_URL")

	log.Printf("Connecting to Elasticsearch at %s", url)

	// uses ELASTICSEARCH_URL env var. If not set - localhpst:9200 is used
	return elasticsearch.NewDefaultClient()
}

func (s *esService) initElasticsearch() error {
	err := s.makeElasticsearchIndex()
	if err != nil {
		return fmt.Errorf("make elasticsearch index: %w", err)
	}

	err = s.indexMessages(context.Background())
	if err != nil {
		return fmt.Errorf("index messages: %w", err)
	}

	return nil
}

func (s *esService) makeElasticsearchIndex() error {
	_, err := s.client.Indices.Create(
		indexNameMessageData,
		s.client.Indices.Create.WithBody(strings.NewReader(`{
			"settings": {
				"number_of_shards": 1
			},
			"mappings": {
				"properties": {
					"data": {
						"type": "text"
					}
				}
			}
		}`)),

		s.client.Indices.Create.WithPretty(),
	)

	return err
}

func (s *esService) indexMessages(ctx context.Context) error {
	var messages []Message
	err := json.Unmarshal(messagesJSON, &messages)
	if err != nil {
		return fmt.Errorf("unmarshal messages: %w", err)
	}

	for _, msg := range messages {
		err = s.indexMessage(ctx, msg)
		if err != nil {
			return fmt.Errorf("index message %s: %w", msg.ID, err)
		}

		log.Printf("message %s has been indexed", msg.ID)
		log.Println()
	}

	return nil
}

type indexRequestMessageData struct {
	Data string `json:"data"`
}

func (s *esService) indexMessage(ctx context.Context, msg Message) error {
	buf, err := json.Marshal(indexRequestMessageData{Data: msg.Data})
	if err != nil {
		return fmt.Errorf("marshal index request: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexNameMessageData,
		DocumentID: msg.ID,
		Body:       bytes.NewReader(buf),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("make index request: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index message %s response error: %s", msg.ID, res.String())
	}

	log.Printf("message %s has been indexed", msg.ID)

	return nil
}

func (s esService) Search(ctx context.Context, data string) error {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"data": data,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex("test"),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Printf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
			log.Println()
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)

	return nil
}
