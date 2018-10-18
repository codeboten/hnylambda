package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/honeycombio/beeline-go"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
type Request events.APIGatewayProxyRequest

type HoneycombLambdaFunction func(ctx context.Context, event map[string]interface{}) (Response, error)

// addRequestProperties adds a variety of details about the lambda request
func addRequestProperties(ctx context.Context) {
	ctx, span := beeline.StartSpan(ctx, "addRequestProperties")
	defer span.Send()
	span.AddField("function_name", lambdacontext.FunctionName)
	span.AddField("function_version", lambdacontext.FunctionVersion)
}

// Middleware will wrap our lambda handler function to create a trace for it
func Middleware(fn func(ctx context.Context, event map[string]interface{}) (Response, error)) func(ctx context.Context, event map[string]interface{}) (Response, error) {
	return func(ctx context.Context, event map[string]interface{}) (Response, error) {
		ctx, span := beeline.StartSpan(ctx, "hnylambda.Middleware")
		span.AddTraceField("platform", "aws")

		defer beeline.Flush(ctx)
		defer span.Send()

		addRequestProperties(ctx)

		resp, err := fn(ctx, event)
		if err != nil {
			span.AddField("lambda.error", err)
		}

		span.AddField("response.status_code", resp.StatusCode)
		return resp, err
	}
}
