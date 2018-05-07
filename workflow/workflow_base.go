package workflow

import (
	"context"
	"encoding/json"
	"reflect"
)

// BaseWorkflow is the base lambda workflow.
type BaseWorkflow struct {
	bootstrap   Bootstrap
	preActions  []Action
	postActions []Action
}

// InvokeHandler invokes the provided handler.
func (w *BaseWorkflow) InvokeHandler(awsContext context.Context, evt interface{}, evtBytes []byte, hData *handlerData) (Context, Error) {
	// Add request parameter to the handler input if there is
	// input parameter.
	req, err := w.getReqParamIfAny(hData.handler, evtBytes)
	if err != nil {
		return w.createContext(awsContext, evt, nil), err
	}

	// Create handler workflow context and register the dependencies
	// in the bootstrap if there are any.
	hContext := w.createContext(awsContext, evt, req)

	// Execute Pre Actions.
	err = w.executeActions(hContext, w.preActions)
	if err != nil {
		return hContext, err
	}

	in := []reflect.Value{
		reflect.ValueOf(hContext),
	}

	if req != nil {
		in = append(in, *req)
	}

	hValue := reflect.ValueOf(hData.handler)
	// Execute Pre Handler Actions.
	err = w.executeActions(hContext, hData.preActions)
	if err != nil {
		return hContext, err
	}

	// Invoke the provided handler.
	out := hValue.Call(in)

	// Execute Post Handler Actions.
	err = w.executeActions(hContext, hData.postActions)
	if err != nil {
		return hContext, err
	}

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

func (w *BaseWorkflow) createContext(ctx context.Context, evt interface{}, req *reflect.Value) *lambdaCtx {
	var injector Injector
	if !reflect.ValueOf(w.bootstrap).IsNil() {
		injector = w.bootstrap()
	}

	return &lambdaCtx{lambdaContext: ctx, lambdaEvent: evt, injector: injector, req: req}
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

func (w *BaseWorkflow) getReqParamIfAny(handler interface{}, evtBytes []byte) (*reflect.Value, Error) {
	if reflect.TypeOf(handler).NumIn() > 1 {
		req, err := w.getHandlerInputFromEvent(handler, evtBytes)
		if err != nil {
			return nil, err
		}

		return &req, nil
	}

	return nil, nil
}
