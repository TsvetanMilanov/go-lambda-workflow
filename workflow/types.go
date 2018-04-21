package workflow

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// PreAction is action which will be executed before the
// handler function. The dependencies in the injector will be
// registered before the action is invoked.
type PreAction func(c Context) error

// PostAction is action which will be executed after the
// handler function.
type PostAction func(c Context) error

// APIGWProxyHandler is AWS API Gateway handler function.
type APIGWProxyHandler func(ctx context.Context, evt events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
