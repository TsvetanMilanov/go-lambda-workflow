package workflow

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
)

// BaseWorkflow is the base lambda workflow.
type BaseWorkflow struct {
	bootstrap   Bootstrap
	preActions  []Action
	postActions []Action
}

func (w *BaseWorkflow) createContext(ctx context.Context, evt interface{}) *lambdaCtx {
	var injector Injector
	if !reflect.ValueOf(w.bootstrap).IsNil() {
		injector = w.bootstrap()
	}

	return &lambdaCtx{lambdaContext: ctx, lambdaEvent: evt, injector: injector}
}

func (w *BaseWorkflow) invokeHandler(awsContext context.Context, evt interface{}, evtBytes []byte, handler interface{}) (*lambdaCtx, Error) {
	// Create handler workflow context and register the dependencies
	// in the bootstrap if there are any.
	hContext := w.createContext(awsContext, evt)

	// Execute Pre Actions.
	err := w.executeActions(hContext, w.preActions)
	if err != nil {
		return hContext, err
	}

	in := []reflect.Value{
		reflect.ValueOf(hContext),
	}

	// Add request parameter to the handler input if there is
	// input parameter.
	hValue := reflect.ValueOf(handler)
	if reflect.TypeOf(handler).NumIn() > 1 {
		req, err := w.getHandlerInputFromEvent(handler, evtBytes)
		if err != nil {
			return hContext, err
		}

		in = append(in, req)
	}

	// Invoke the provided handler.
	out := hValue.Call(in)

	// Execute Post Action.
	err = w.executeActions(hContext, w.postActions)
	if err != nil {
		return hContext, err
	}

	resErr := out[0].Interface()
	if resErr != nil {
		err, ok := resErr.(error)
		if !ok {
			return hContext, newErrorWithMessage("invalid handler error result")
		}

		return hContext, newError(err)
	}

	return hContext, nil
}

func (w *BaseWorkflow) getHandlerInputFromEvent(handler interface{}, evt []byte) (reflect.Value, Error) {
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
		return reflect.Value{}, newError(err)
	}

	if inputType.Kind() == reflect.Ptr {
		return inputValue, nil
	}

	return inputValue.Elem(), nil
}

func (w *BaseWorkflow) executeActions(c Context, actions []Action) Error {
	for _, a := range w.postActions {
		err := a(c)
		if err != nil {
			return newError(err)
		}
	}

	return nil
}

// APIGatewayProxyWorkflow AWS API Gateway Lambda Proxy request/response workflow.
type APIGatewayProxyWorkflow struct {
	*BaseWorkflow
	httpHandlers map[string]interface{}
}

// GetLambdaHandler returns AWS API Gateway Proxy Lambda handler.
func (w *APIGatewayProxyWorkflow) GetLambdaHandler() APIGWProxyHandler {
	return func(ctx context.Context, evt events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		h, ok := w.httpHandlers[getHandlerKey(evt.HTTPMethod, evt.Path)]
		if !ok {
			return defaultAPIGWProxyHandler(ctx, evt)
		}

		var reqBytes []byte
		var err error
		hType := reflect.TypeOf(h)
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

		hContext, err := w.BaseWorkflow.invokeHandler(ctx, evt, reqBytes, h)
		if err != nil {
			return nil, err
		}

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
