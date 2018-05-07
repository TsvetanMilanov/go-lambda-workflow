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

func (c *lambdaCtx) GetLambdaEvent(out interface{}) (err Error) {
	defer func() {
		if r := recover(); r != nil {
			err = newErrorWithMessage("cannot get lambda event, reson is: %s", r)
		}
	}()

	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr {
		return newErrorWithMessage("the out parameter must be a pointer")
	}

	if !outValue.Elem().CanSet() {
		return newErrorWithMessage("can't set lambda event out value")
	}

	evtType := reflect.TypeOf(c.lambdaEvent)
	outType := outValue.Elem().Type()
	if evtType != outType {
		return newErrorWithMessage("cannot set event of type %s to output of type %s", evtType.Name(), outType.Name())
	}

	outValue.Elem().Set(reflect.ValueOf(c.lambdaEvent))
	return err
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
