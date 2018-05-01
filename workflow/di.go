package workflow

// Bootstrap is function which registers all dependencies in the
// injector and returns it.
type Bootstrap func() Injector

// Injector describes DI related operations.
type Injector interface {
	Resolve(out interface{}) error
	ResolveByName(name string, out interface{}) error
}
