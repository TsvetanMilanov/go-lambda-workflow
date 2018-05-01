package workflow

import (
	"net/http"
)

// BaseWorkflowBuilder is the base workflow builder.
type BaseWorkflowBuilder struct {
	bootstrap   Bootstrap
	preActions  []Action
	postActions []Action
}

// SetBootstrap sets the bootstrap function to the workflow.
func (b *BaseWorkflowBuilder) SetBootstrap(bootstrap Bootstrap) *BaseWorkflowBuilder {
	b.bootstrap = bootstrap
	return b
}

// AddPreAction adds Pre Action to the workflow.
func (b *BaseWorkflowBuilder) AddPreAction(action Action) *BaseWorkflowBuilder {
	b.preActions = append(b.preActions, action)
	return b
}

// AddPostAction adds Post Action to the workflow.
func (b *BaseWorkflowBuilder) AddPostAction(action Action) *BaseWorkflowBuilder {
	b.postActions = append(b.postActions, action)
	return b
}

// Build creates the Base workflow.
func (b *BaseWorkflowBuilder) Build() *BaseWorkflow {
	return &BaseWorkflow{bootstrap: b.bootstrap, preActions: b.preActions, postActions: b.postActions}
}

// NewBaseWorkflowBuilder creates new Base workflow builder.
func NewBaseWorkflowBuilder() *BaseWorkflowBuilder {
	return &BaseWorkflowBuilder{
		preActions:  []Action{},
		postActions: []Action{},
	}
}

// APIGWProxyWorkflowBuilder AWS Lambda handler workflow builder.
type APIGWProxyWorkflowBuilder struct {
	*BaseWorkflowBuilder
	httpHandlers map[string]interface{}
}

// AddGetHandler adds the provided handler to the specified path and GET HTTP method.
func (b *APIGWProxyWorkflowBuilder) AddGetHandler(path string, handler interface{}) *APIGWProxyWorkflowBuilder {
	// TODO: Validate handler func.
	b.httpHandlers[getHandlerKey(http.MethodGet, path)] = handler
	return b
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
		httpHandlers:        make(map[string]interface{}),
	}
}
