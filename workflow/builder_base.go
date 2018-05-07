package workflow

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
