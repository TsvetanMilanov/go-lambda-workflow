package workflow

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// APIGatewayProxyWorkflow AWS API GateWay Lambda Proxy request/response workflow.
type APIGatewayProxyWorkflow func(ctx context.Context, evt events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

// NewAPIGWProxyWorkflowBuilder creates new AWS API Gateway Proxy workflow builder.
func NewAPIGWProxyWorkflowBuilder() *APIGWProxyWorkflowBuilder {
	return &APIGWProxyWorkflowBuilder{httpHandlers: make(map[string]interface{})}
}
