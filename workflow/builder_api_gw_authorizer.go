package workflow

// APIGWAuthorizerWorkflowBuilder AWS Lambda handler workflow builder.
type APIGWAuthorizerWorkflowBuilder struct {
	*BaseWorkflowBuilder
	handler *handlerData
}

// SetHandler sets the provided handler as the API GW Authorizer.
func (b *APIGWAuthorizerWorkflowBuilder) SetHandler(handler interface{}) *APIGWAuthorizerPrePostHandlerActionBuilder {
	// TODO: Validate handler func.
	hData := &handlerData{handler: handler, preActions: []Action{}, postActions: []Action{}}
	b.handler = hData
	return newAPIGWAuthorizerPrePostHandlerActionBuilder(b)
}

// SetBootstrap override just to return the correct builder.
func (b *APIGWAuthorizerWorkflowBuilder) SetBootstrap(bootstrap Bootstrap) *APIGWAuthorizerWorkflowBuilder {
	b.BaseWorkflowBuilder.SetBootstrap(bootstrap)
	return b
}

// AddPreActions adds Pre Actions to the workflow.
func (b *APIGWAuthorizerWorkflowBuilder) AddPreActions(actions ...Action) *APIGWAuthorizerWorkflowBuilder {
	b.BaseWorkflowBuilder.AddPreActions(actions...)
	return b
}

// AddPostActions adds Post Actions to the workflow.
func (b *APIGWAuthorizerWorkflowBuilder) AddPostActions(actions ...Action) *APIGWAuthorizerWorkflowBuilder {
	b.BaseWorkflowBuilder.AddPostActions(actions...)
	return b
}

// Build creates the AWS Lambda workflow.
func (b *APIGWAuthorizerWorkflowBuilder) Build() *APIGatewayAuthorizerWorkflow {
	return &APIGatewayAuthorizerWorkflow{
		BaseWorkflow: b.BaseWorkflowBuilder.Build(),
		handler:      b.handler,
	}
}

// NewAPIGWAuthorizerWorkflowBuilder creates new AWS API Gateway Authorizer workflow builder.
func NewAPIGWAuthorizerWorkflowBuilder() *APIGWAuthorizerWorkflowBuilder {
	return &APIGWAuthorizerWorkflowBuilder{
		BaseWorkflowBuilder: NewBaseWorkflowBuilder(),
	}
}

// APIGWAuthorizerPrePostHandlerActionBuilder is the builder which enables
// adding pre and post handler actions with fluent API.
type APIGWAuthorizerPrePostHandlerActionBuilder struct {
	*APIGWAuthorizerWorkflowBuilder
}

// WithPreActions adds the pre actions to the previously added handler.
func (b *APIGWAuthorizerPrePostHandlerActionBuilder) WithPreActions(actions ...Action) *APIGWAuthorizerPrePostHandlerActionBuilder {
	b.handler.preActions = append(b.handler.preActions, actions...)
	return b
}

// WithPostActions adds the post actions to the previously added handler.
func (b *APIGWAuthorizerPrePostHandlerActionBuilder) WithPostActions(actions ...Action) *APIGWAuthorizerPrePostHandlerActionBuilder {
	b.handler.postActions = append(b.handler.postActions, actions...)
	return b
}

func newAPIGWAuthorizerPrePostHandlerActionBuilder(b *APIGWAuthorizerWorkflowBuilder) *APIGWAuthorizerPrePostHandlerActionBuilder {
	return &APIGWAuthorizerPrePostHandlerActionBuilder{APIGWAuthorizerWorkflowBuilder: b}
}
