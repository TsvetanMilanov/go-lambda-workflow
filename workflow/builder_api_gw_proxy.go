package workflow

import (
	"net/http"
	"regexp"
)

var (
	parameterizedPathRegExp = regexp.MustCompile(`.*(\{.*\}).*`)
)

// APIGWProxyWorkflowBuilder AWS Lambda handler workflow builder.
type APIGWProxyWorkflowBuilder struct {
	*BaseWorkflowBuilder
	httpHandlers              map[string]*handlerData
	parameterizedHTTPHandlers []*parameterizedHandlerData
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

	// TODO: Check if path already exist.
	if b.isParameterizedPath(path) {
		phData := &parameterizedHandlerData{
			hData:        hData,
			pathRegExp:   regexp.MustCompile(parameterizedPathRegExp.ReplaceAllString(path, "(.*)")),
			originalPath: path,
			method:       httpMethod,
		}
		b.parameterizedHTTPHandlers = append(b.parameterizedHTTPHandlers, phData)
	} else {
		b.httpHandlers[getHandlerKey(httpMethod, path)] = hData
	}
	return newAPIGWPrePostHandlerActionBuilder(b, hData)
}

// SetBootstrap override just to return the correct builder.
func (b *APIGWProxyWorkflowBuilder) SetBootstrap(bootstrap Bootstrap) *APIGWProxyWorkflowBuilder {
	b.BaseWorkflowBuilder.SetBootstrap(bootstrap)
	return b
}

// AddPreActions adds Pre Actions to the workflow.
func (b *APIGWProxyWorkflowBuilder) AddPreActions(actions ...Action) *APIGWProxyWorkflowBuilder {
	b.BaseWorkflowBuilder.AddPreActions(actions...)
	return b
}

// AddPostActions adds Post Actions to the workflow.
func (b *APIGWProxyWorkflowBuilder) AddPostActions(actions ...Action) *APIGWProxyWorkflowBuilder {
	b.BaseWorkflowBuilder.AddPostActions(actions...)
	return b
}

// Build creates the AWS Lambda workflow.
func (b *APIGWProxyWorkflowBuilder) Build() *APIGatewayProxyWorkflow {
	return &APIGatewayProxyWorkflow{
		BaseWorkflow:              b.BaseWorkflowBuilder.Build(),
		httpHandlers:              b.httpHandlers,
		parameterizedHTTPHandlers: b.parameterizedHTTPHandlers,
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

func (b *APIGWProxyWorkflowBuilder) isParameterizedPath(path string) bool {
	return len(parameterizedPathRegExp.FindStringSubmatchIndex(path)) > 0
}

func newAPIGWPrePostHandlerActionBuilder(b *APIGWProxyWorkflowBuilder, handler *handlerData) *APIGWPrePostHandlerActionBuilder {
	return &APIGWPrePostHandlerActionBuilder{APIGWProxyWorkflowBuilder: b, handler: handler}
}

type handlerData struct {
	handler     interface{}
	preActions  []Action
	postActions []Action
}

type parameterizedHandlerData struct {
	hData        *handlerData
	pathRegExp   *regexp.Regexp
	originalPath string
	method       string
}
