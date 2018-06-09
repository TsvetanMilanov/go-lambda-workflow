package workflow

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
)

// APIGatewayProxyWorkflow AWS API Gateway Lambda Proxy request/response workflow.
type APIGatewayProxyWorkflow struct {
	*BaseWorkflow
	httpHandlers              map[string]*handlerData
	parameterizedHTTPHandlers map[string]*handlerData
}

// GetLambdaHandler returns AWS API Gateway Proxy Lambda handler.
func (w *APIGatewayProxyWorkflow) GetLambdaHandler() APIGWProxyHandler {
	return func(ctx context.Context, evt events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		hData, ok := w.httpHandlers[getHandlerKey(evt.HTTPMethod, evt.Path)]
		if !ok {
			return defaultAPIGWProxyHandler(ctx, evt)
		}

		var reqBytes []byte
		var err error
		hType := reflect.TypeOf(hData.handler)
		// Get event bytes only if the handler has input parameter.
		if hType.NumIn() > 1 {
			// Use directly the event body if the handler input parameter
			// has type string or []byte.
			handlerInputType := hType.In(1)
			if handlerInputType.Kind() == reflect.Struct {
				reqBytes, err = w.getReqBytes(evt)
				if err != nil {
					return nil, err
				}
			} else {
				reqBytes = []byte(evt.Body)
			}
		}

		c, err := w.BaseWorkflow.InvokeHandler(ctx, evt, reqBytes, hData)
		if err != nil {
			return nil, err
		}

		hContext := c.(*lambdaCtx)
		// Handle Raw response.
		if hContext.rawResponse != nil {
			if r, ok := hContext.rawResponse.(events.APIGatewayProxyResponse); ok {
				return &r, nil
			} else if r, ok := hContext.rawResponse.(*events.APIGatewayProxyResponse); ok {
				return r, nil
			} else {
				return nil, newErrorWithMessage("invalid raw response")
			}
		}

		// Handle response body.
		var resBytes []byte
		var mErr error
		if hContext.response != nil {
			resBytes, mErr = json.Marshal(hContext.response)
			if mErr != nil {
				return nil, newError(mErr)
			}
		}

		var resBody string
		if len(resBytes) > 0 {
			resBody = string(resBytes)
		}

		proxyRes := events.APIGatewayProxyResponse{
			StatusCode: hContext.responseStatusCode,
			Body:       resBody,
		}
		return &proxyRes, nil
	}
}

func (w *APIGatewayProxyWorkflow) getReqBytes(evt events.APIGatewayProxyRequest) ([]byte, Error) {
	input := make(map[string]interface{})
	if len(evt.Body) > 0 {
		err := json.Unmarshal([]byte(evt.Body), &input)
		if err != nil {
			return nil, newError(err)
		}
	}

	for k, v := range evt.Headers {
		input[k] = v
	}

	for k, v := range evt.PathParameters {
		input[k] = v
	}

	for k, v := range evt.QueryStringParameters {
		input[k] = v
	}

	res, err := json.Marshal(input)
	return res, newError(err)
}

func (w *APIGatewayProxyWorkflow) getHandler(evt events.APIGatewayProxyRequest) *handlerData {
	return nil
}
