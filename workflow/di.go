package workflow

// Injector describes DI related operations.
type Injector interface {
	Resolve(depType interface{}, out interface{}) error
	ResolveByName(name string, out interface{}) error
}
