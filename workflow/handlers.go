package workflow

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// APIGWProxyHandler is AWS API Gateway handler function.
type APIGWProxyHandler func(ctx context.Context, evt events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

var defaultAPIGWProxyHandler = func(ctx context.Context, evt events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
}
