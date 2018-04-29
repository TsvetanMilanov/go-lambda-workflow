package workflow

// Action describes function which will be executed before or
// after the handler is executed.
type Action func(c Context) error
