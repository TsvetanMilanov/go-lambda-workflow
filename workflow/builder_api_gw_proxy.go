package workflow

import (
	"net/http"
)

// APIGWProxyWorkflowBuilder AWS Lambda handler workflow builder.
type APIGWProxyWorkflowBuilder struct {
	*BaseWorkflowBuilder
	httpHandlers map[string]*handlerData
}

// AddGetHandler adds the provided handler to the specified path and GET HTTP method.
func (b *APIGWProxyWorkflowBuilder) AddGetHandler(path string, handler interface{}) *APIGWPrePostHandlerActionBuilder {
	return b.AddMethodHandler(http.MethodGet, path, handler)
}

// AddPostHandler adds the provided handler to the specified path and POST HTTP method.
func (b *APIGWProxyWorkflowBuilder) AddPostHandler(path string, handler interface{}) *APIGWPrePostHandlerActionBuilder {
	return b.AddMethodHandler(http.MethodPost, path, handler)
}

// AddPutHandler adds the provided handler to the specified path and PUT HTTP method.
func (b *APIGWProxyWorkflowBuilder) AddPutHandler(path string, handler interface{}) *APIGWPrePostHandlerActionBuilder {
	return b.AddMethodHandler(http.MethodPut, path, handler)
}

// AddDeleteHandler adds the provided handler to the specified path and DELETE HTTP method.
func (b *APIGWProxyWorkflowBuilder) AddDeleteHandler(path string, handler interface{}) *APIGWPrePostHandlerActionBuilder {
	return b.AddMethodHandler(http.MethodDelete, path, handler)
}

// AddMethodHandler adds the provided handler to the specified path with the provided HTTP method.
func (b *APIGWProxyWorkflowBuilder) AddMethodHandler(httpMethod, path string, handler interface{}) *APIGWPrePostHandlerActionBuilder {
	// TODO: Validate handler func.
	hData := &handlerData{handler: handler, preActions: []Action{}, postActions: []Action{}}
	b.httpHandlers[getHandlerKey(httpMethod, path)] = hData
	return newAPIGWPrePostHandlerActionBuilder(b, hData)
}

// SetBootstrap override just to return the correct builder.
func (b *APIGWProxyWorkflowBuilder) SetBootstrap(bootstrap Bootstrap) *APIGWProxyWorkflowBuilder {
	b.BaseWorkflowBuilder.SetBootstrap(bootstrap)
	return b
}

// AddPreAction adds Pre Action to the workflow.
func (b *APIGWProxyWorkflowBuilder) AddPreAction(action Action) *APIGWProxyWorkflowBuilder {
	b.BaseWorkflowBuilder.AddPreAction(action)
	return b
}

// AddPostAction adds Post Action to the workflow.
func (b *APIGWProxyWorkflowBuilder) AddPostAction(action Action) *APIGWProxyWorkflowBuilder {
	b.BaseWorkflowBuilder.AddPostAction(action)
	return b
}

// Build creates the AWS Lambda workflow.
func (b *APIGWProxyWorkflowBuilder) Build() *APIGatewayProxyWorkflow {
	return &APIGatewayProxyWorkflow{
		BaseWorkflow: b.BaseWorkflowBuilder.Build(),
		httpHandlers: b.httpHandlers,
	}
}

// NewAPIGWProxyWorkflowBuilder creates new AWS API Gateway Proxy workflow builder.
func NewAPIGWProxyWorkflowBuilder() *APIGWProxyWorkflowBuilder {
	return &APIGWProxyWorkflowBuilder{
		BaseWorkflowBuilder: NewBaseWorkflowBuilder(),
		httpHandlers:        make(map[string]*handlerData),
	}
}

// APIGWPrePostHandlerActionBuilder is the builder which enables the
// adding pre and post handler actions with fluent API.
type APIGWPrePostHandlerActionBuilder struct {
	*APIGWProxyWorkflowBuilder
	handler *handlerData
}

// WithPreActions adds the pre actions to the previously added handler.
func (b *APIGWPrePostHandlerActionBuilder) WithPreActions(actions ...Action) *APIGWPrePostHandlerActionBuilder {
	b.handler.preActions = append(b.handler.preActions, actions...)
	return b
}

// WithPostActions adds the post actions to the previously added handler.
func (b *APIGWPrePostHandlerActionBuilder) WithPostActions(actions ...Action) *APIGWPrePostHandlerActionBuilder {
	b.handler.postActions = append(b.handler.postActions, actions...)
	return b
}

func newAPIGWPrePostHandlerActionBuilder(b *APIGWProxyWorkflowBuilder, handler *handlerData) *APIGWPrePostHandlerActionBuilder {
	return &APIGWPrePostHandlerActionBuilder{APIGWProxyWorkflowBuilder: b, handler: handler}
}

type handlerData struct {
	handler     interface{}
	preActions  []Action
	postActions []Action
}
