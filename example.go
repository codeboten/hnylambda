package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"reflect"

	"github.com/aws/aws-lambda-go/lambda"
	hnylambda "github.com/codeboten/hnylambda/handler"
	"github.com/honeycombio/beeline-go"
)

type weatherRequestEvent struct {
	City string `json:"city"`
}

func (w *weatherRequestEvent) Marshal(input map[string]interface{}) {
	if city, ok := input["city"]; ok {
		if reflect.ValueOf(city).Kind() == reflect.String {
			w.City = city.(string)
		}
	}
}

// Handler is the lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, input map[string]interface{}) (hnylambda.Response, error) {
	event := weatherRequestEvent{}
	event.Marshal(input)

	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{
		"city":    event.City,
		"weather": "weather is fine",
	})
	if err != nil {
		return hnylambda.Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := hnylambda.Response{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	beeline.Init(beeline.Config{
		WriteKey: os.Getenv("HONEYCOMB_KEY"),
		Dataset:  os.Getenv("HONEYCOMB_DATASET"),
	})
	lambda.Start(hnylambda.Middleware(Handler))
}
