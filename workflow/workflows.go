package workflow

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// APIGatewayProxyWorkflow AWS API GateWay Lambda Proxy request/response workflow.
type APIGatewayProxyWorkflow func(ctx context.Context, evt events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
