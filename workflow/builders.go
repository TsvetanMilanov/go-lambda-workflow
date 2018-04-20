package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// APIGWProxyWorkflowBuilder AWS Lambda handler workflow builder.
type APIGWProxyWorkflowBuilder struct {
	httpHandlers map[string]interface{}
}

// AddGetHandler adds the provided handler to the specified path and GET HTTP method.
func (b *APIGWProxyWorkflowBuilder) AddGetHandler(path string, handler interface{}) *APIGWProxyWorkflowBuilder {
	// TODO: Validate handler func.
	b.httpHandlers[b.getHandlerKey(http.MethodGet, path)] = handler
	return b
}

// Build creates the AWS Lambda workflow.
func (b *APIGWProxyWorkflowBuilder) Build() APIGatewayProxyWorkflow {
	return func(ctx context.Context, evt events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		reqMap, err := b.getReqMap(evt)
		if err != nil {
			return nil, err
		}

		reqBytes, err := json.Marshal(reqMap)
		if err != nil {
			return nil, err
		}

		h := b.httpHandlers[b.getHandlerKey(evt.HTTPMethod, evt.Path)]

		hType := reflect.TypeOf(h)
		req := reflect.New(hType.In(1))
		err = json.Unmarshal(reqBytes, req.Interface())
		if err != nil {
			return nil, err
		}

		handlerCtx := &lambdaCtx{lambdaContext: ctx, lambdaEvent: evt}
		in := []reflect.Value{
			reflect.ValueOf(handlerCtx),
			req.Elem(),
		}

		hValue := reflect.ValueOf(h)
		out := hValue.Call(in)
		resErr := out[0].Interface()
		if resErr != nil {
			err, ok := resErr.(error)
			if !ok {
				return nil, errors.New("invalid handler error result")
			}

			return nil, err
		}

		resBody, err := json.Marshal(handlerCtx.response)
		if err != nil {
			return nil, err
		}

		proxyRes := events.APIGatewayProxyResponse{
			StatusCode: handlerCtx.responseStatusCode,
			Body:       string(resBody),
		}
		return &proxyRes, nil
	}
}

func (b *APIGWProxyWorkflowBuilder) getHandlerKey(method, path string) string {
	// TODO: sanitize path.
	return fmt.Sprintf("%s-%s", strings.ToLower(method), path)
}

func (b *APIGWProxyWorkflowBuilder) getReqMap(evt events.APIGatewayProxyRequest) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if len(evt.Body) > 0 {
		err := json.Unmarshal([]byte(evt.Body), result)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range evt.Headers {
		result[k] = v
	}

	for k, v := range evt.PathParameters {
		result[k] = v
	}

	for k, v := range evt.QueryStringParameters {
		result[k] = v
	}

	return result, nil
}
