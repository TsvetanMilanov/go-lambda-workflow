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

// AddPreActions adds Pre Actions to the workflow.
func (b *BaseWorkflowBuilder) AddPreActions(actions ...Action) *BaseWorkflowBuilder {
	b.preActions = append(b.preActions, actions...)
	return b
}

// AddPostActions adds Post Actions to the workflow.
func (b *BaseWorkflowBuilder) AddPostActions(actions ...Action) *BaseWorkflowBuilder {
	b.postActions = append(b.postActions, actions...)
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
