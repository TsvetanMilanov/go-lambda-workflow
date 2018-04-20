package workflow

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// Context is the AWS Lambda Workflow context.
type Context interface {
	GetLambdaContext() context.Context
	GetLambdaEvent(out interface{}) error
	SetResponse(interface{}) Context
	SetResponseStatusCode(int) Context
}

type lambdaCtx struct {
	lambdaContext      context.Context
	lambdaEvent        interface{}
	response           interface{}
	responseStatusCode int
}

func (c *lambdaCtx) SetResponse(res interface{}) Context {
	c.response = res
	return c
}

func (c *lambdaCtx) SetResponseStatusCode(code int) Context {
	c.responseStatusCode = code
	return c
}

func (c *lambdaCtx) GetLambdaContext() context.Context {
	return c.lambdaContext
}

func (c *lambdaCtx) GetLambdaEvent(out interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("cannot fet lambda event, reson is: %s", r)
		}
	}()

	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr {
		return errors.New("the out parameter must be a pointer")
	}

	if !outValue.Elem().CanSet() {
		return errors.New("can't set lambda event out value")
	}

	evtType := reflect.TypeOf(c.lambdaEvent)
	outType := outValue.Elem().Type()
	if evtType != outType {
		return fmt.Errorf("cannot set event of type %s to output of type %s", evtType.Name(), outType.Name())
	}

	outValue.Elem().Set(reflect.ValueOf(c.lambdaEvent))
	return err
}
