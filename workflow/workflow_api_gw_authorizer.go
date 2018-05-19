package workflow

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// APIGatewayAuthorizerWorkflow AWS API Gateway Authorizer workflow.
type APIGatewayAuthorizerWorkflow struct {
	*BaseWorkflow
	handler *handlerData
}

// GetLambdaHandler returns AWS API Gateway Authorizer Lambda handler.
func (w *APIGatewayAuthorizerWorkflow) GetLambdaHandler() APIGWAuthorizerHandler {
	return func(ctx context.Context, evt events.APIGatewayCustomAuthorizerRequest) (*events.APIGatewayCustomAuthorizerResponse, error) {
		reqBytes, err := json.Marshal(evt)
		if err != nil {
			return nil, err
		}

		c, err := w.BaseWorkflow.InvokeHandler(ctx, evt, reqBytes, w.handler)
		if err != nil {
			return nil, err
		}

		hContext := c.(*lambdaCtx)

		var res interface{}
		if hContext.rawResponse != nil {
			res = hContext.rawResponse
		} else if hContext.response != nil {
			res = hContext.response
		}

		// Handle the response.
		if res != nil {
			if r, ok := res.(events.APIGatewayCustomAuthorizerResponse); ok {
				return &r, nil
			} else if r, ok := res.(*events.APIGatewayCustomAuthorizerResponse); ok {
				return r, nil
			} else {
				return nil, newErrorWithMessage("invalid response")
			}
		}

		return nil, nil
	}
}
