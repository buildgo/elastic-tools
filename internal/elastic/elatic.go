package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"math"
	"time"
)

type Elastic struct {
	client *elasticsearch.Client
}

type DummyDocument struct {
	Id          string    `json:"id"`
	Batch       int       `json:"batch""`
	Count       int       `json:"count"`
	CreatedTime time.Time `json:"createdTime"`
}

func CreateClient(address []string) (*Elastic, error) {
	cfg := elasticsearch.Config{
		Addresses: address,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Elastic{
		client: client,
	}, nil
}

func Info(elastic *Elastic) {
	response, err := elastic.client.Info()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(response)
}

func InsertDummyDocument(elastic *Elastic, ctx context.Context, index string, batch, size, sleep int) error {
	if batch == 0 {
		batch = math.MaxInt
	}
	for i := 0; i < batch; i++ {
		fmt.Printf("InsertDummyDocument index=%s, batch=%d, size=%d\n", index, i, size)
		for _, doc := range createDummyDocument(batch, size) {
			bdy, err := json.Marshal(doc)
			if err != nil {
				return fmt.Errorf("insert: marshall: %w", err)
			}
			req := esapi.CreateRequest{
				Index:      index,
				DocumentID: doc.Id,
				Body:       bytes.NewReader(bdy),
			}

			//esapi.UpdateRequest{
			//	Index:      index,
			//	DocumentID: doc.Id,
			//	Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc": %s}`, bdy))),
			//}

			ctx, cancel := context.WithTimeout(ctx, time.Duration(1)*time.Second)
			defer cancel()

			res, err := req.Do(ctx, elastic.client)
			if err != nil {
				return fmt.Errorf("insert: request: %w", err)
			}
			defer res.Body.Close()

			if res.StatusCode == 409 {
				return errors.New("conflict")
			}
		}
		time.Sleep(time.Duration(sleep) * time.Second)
	}
	return nil
}

func createDummyDocument(batch, size int) []DummyDocument {
	var docs []DummyDocument
	now := time.Now()
	custom := now.Format("2006_0102_150405")
	for i := 0; i < size; i++ {
		docId := fmt.Sprintf("%s-%d-%d", custom, batch, i)
		doc := DummyDocument{
			Id:          docId,
			Batch:       batch,
			Count:       i,
			CreatedTime: now,
		}
		docs = append(docs, doc)
	}
	return docs
}
