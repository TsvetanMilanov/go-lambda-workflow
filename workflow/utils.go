package workflow

import (
	"fmt"
	"reflect"
	"strings"
)

func getHandlerKey(method, path string) string {
	// TODO: sanitize path.
	return fmt.Sprintf("%s-%s", strings.ToLower(method), path)
}

func hasResponse(ctx *lambdaCtx) bool {
	return ctx.rawResponse != nil || ctx.response != nil
}

func validateOutParameterValue(out, valueToSet interface{}, valueName string) (err Error) {
	defer func() {
		if r := recover(); r != nil {
			err = newErrorWithMessage("invalid value for %s, reson is: %s", valueName, r)
		}
	}()

	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr {
		return newErrorWithMessage("the out parameter must be a pointer")
	}

	if !outValue.Elem().CanSet() {
		return newErrorWithMessage("can't set %s out value", valueName)
	}

	valueToSetType := reflect.ValueOf(valueToSet).Elem().Type()
	outType := outValue.Elem().Type()
	if valueToSetType != outType {
		return newErrorWithMessage("cannot set %s of type %s to output of type %s", valueName, valueToSetType.Name(), outType.Name())
	}

	return nil
}

func setOutParameterValue(out, valueToSet interface{}, valueName string) Error {
	if valueToSet == nil {
		out = nil
		return nil
	}

	err := validateOutParameterValue(out, valueToSet, valueName)
	if err != nil {
		return err
	}

	refValue := reflect.ValueOf(valueToSet)
	outElem := reflect.ValueOf(out).Elem()
	if refValue.Kind() == reflect.Ptr {
		outElem.Set(refValue.Elem())
	} else {
		outElem.Set(refValue)
	}

	return nil
}
