package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
)

// BaseWorkflow is the base lambda workflow.
type BaseWorkflow struct {
	bootstrap   Bootstrap
	preActions  []PreAction
	postActions []PostAction
}

func (w *BaseWorkflow) createContext(ctx context.Context, evt interface{}) *lambdaCtx {
	var injector Injector
	if !reflect.ValueOf(w.bootstrap).IsNil() {
		injector = w.bootstrap()
	}

	return &lambdaCtx{lambdaContext: ctx, lambdaEvent: evt, injector: injector}
}

func (w *BaseWorkflow) invokeHandler(awsContext context.Context, evt interface{}, evtBytes []byte, handler interface{}) (*lambdaCtx, error) {
	// Create handler workflow context and register the dependencies
	// in the bootstrap if there are any.
	hContext := w.createContext(awsContext, evt)

	// Execute Pre Actions.
	for _, a := range w.preActions {
		err := a(hContext)
		if err != nil {
			return hContext, err
		}
	}

	in := []reflect.Value{
		reflect.ValueOf(hContext),
	}

	// Add request parameter to the handler input if there is
	// input parameter.
	hValue := reflect.ValueOf(handler)
	if hValue.Type().NumIn() > 1 {
		req, err := w.getHandlerInputFromEvent(handler, evtBytes)
		if err != nil {
			return hContext, err
		}

		in = append(in, req)
	}

	out := hValue.Call(in)
	for _, a := range w.postActions {
		err := a(hContext)
		if err != nil {
			return hContext, err
		}
	}

	resErr := out[0].Interface()
	if resErr != nil {
		err, ok := resErr.(error)
		if !ok {
			return hContext, errors.New("invalid handler error result")
		}

		return hContext, err
	}

	return hContext, nil
}

func (w *BaseWorkflow) getHandlerInputFromEvent(handler interface{}, evt []byte) (reflect.Value, error) {
	hType := reflect.TypeOf(handler)
	inputType := hType.In(1)
	// The inputValue will always be have type *inputType.
	var inputValue reflect.Value
	if inputType.Kind() == reflect.Ptr {
		// If the input parameter is ptr we need to create reflect.Value
		// which is *inputType not **inputType.
		inputValue = reflect.New(inputType.Elem())
	} else {
		inputValue = reflect.New(inputType)
	}

	err := json.Unmarshal(evt, inputValue.Interface())
	if err != nil {
		return reflect.Value{}, err
	}

	if inputType.Kind() == reflect.Ptr {
		return inputValue, nil
	}

	return inputValue.Elem(), nil
}

// APIGatewayProxyWorkflow AWS API Gateway Lambda Proxy request/response workflow.
type APIGatewayProxyWorkflow struct {
	*BaseWorkflow
	httpHandlers map[string]interface{}
}

// GetLambdaHandler returns AWS API Gateway Proxy Lambda handler.
func (w *APIGatewayProxyWorkflow) GetLambdaHandler() APIGWProxyHandler {
	return func(ctx context.Context, evt events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		reqBytes, err := w.getReqBytes(evt)
		if err != nil {
			return nil, err
		}

		h := w.httpHandlers[getHandlerKey(evt.HTTPMethod, evt.Path)]
		hContext, err := w.BaseWorkflow.invokeHandler(ctx, evt, reqBytes, h)
		if err != nil {
			return nil, err
		}

		resBody, err := json.Marshal(hContext.response)
		if err != nil {
			return nil, err
		}

		proxyRes := events.APIGatewayProxyResponse{
			StatusCode: hContext.responseStatusCode,
			Body:       string(resBody),
		}
		return &proxyRes, nil
	}
}

func (w *APIGatewayProxyWorkflow) getReqBytes(evt events.APIGatewayProxyRequest) ([]byte, error) {
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

	return json.Marshal(result)
}
