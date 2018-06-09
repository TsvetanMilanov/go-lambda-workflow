package workflow

import (
	"context"
	"reflect"
)

// Context is the AWS Lambda Workflow context.
type Context interface {
	GetLambdaContext() context.Context
	GetLambdaEvent(out interface{}) Error
	GetInjector() Injector
	GetRequestObject(out interface{}) Error
	GetRawResponse(out interface{}) Error
	GetHandlerError() error
	GetRequest() interface{}
	SetResponse(interface{}) Context
	SetRawResponse(interface{}) Context
	SetResponseStatusCode(int) Context
}

type lambdaCtx struct {
	// Set by the builder
	lambdaContext context.Context
	lambdaEvent   interface{}
	injector      Injector
	req           *reflect.Value

	// Set by the user
	response           interface{}
	rawResponse        interface{}
	responseStatusCode int

	handlerErr error
}

func (c *lambdaCtx) SetResponse(res interface{}) Context {
	c.response = res
	return c
}

func (c *lambdaCtx) SetRawResponse(res interface{}) Context {
	c.rawResponse = res
	return c
}

func (c *lambdaCtx) SetResponseStatusCode(code int) Context {
	c.responseStatusCode = code
	return c
}

func (c *lambdaCtx) GetLambdaContext() context.Context {
	return c.lambdaContext
}

func (c *lambdaCtx) GetLambdaEvent(out interface{}) Error {
	return setOutParameterValue(out, c.lambdaEvent, "lambda event")
}

func (c *lambdaCtx) GetRawResponse(out interface{}) Error {
	return setOutParameterValue(out, c.rawResponse, "raw response")
}

func (c *lambdaCtx) GetHandlerError() error {
	return c.handlerErr
}

func (c *lambdaCtx) GetInjector() Injector {
	return c.injector
}

func (c *lambdaCtx) GetRequestObject(out interface{}) Error {
	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr {
		return newErrorWithMessage("the out parameter is not a pointer")
	}

	if c.req != nil {
		reflect.ValueOf(out).Elem().Set(c.req.Elem())
	}

	return nil
}

func (c *lambdaCtx) GetRequest() interface{} {
	return c.req.Interface()
}
