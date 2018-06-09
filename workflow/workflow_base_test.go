package workflow

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBaseWorkflow(t *testing.T) {
	Convey("Base workflow", t, func() {
		Convey("Should set the handler error in the handler context.", func() {
			handlerError := errors.New("handler error")

			assertErr := func(actual error) {
				So(actual, ShouldNotBeNil)
				So(actual, ShouldBeError, handlerError.Error())
			}

			postAction := func(c Context) error {
				assertErr(c.GetHandlerError())
				return nil
			}

			w := NewBaseWorkflowBuilder().
				AddPostActions(postAction). // Validate that the handler error is set before invoking the workflow post actions.
				Build()

			hData := &handlerData{
				handler: func(Context) error {
					return handlerError
				},
				postActions: []Action{postAction}, // Validate that the handler error is set before invoking the handler post actions.
			}

			_, err := w.InvokeHandler(nil, nil, nil, hData)

			assertErr(err)
		})
	})
}
