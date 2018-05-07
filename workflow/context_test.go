package workflow

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContext(t *testing.T) {
	Convey("Context", t, func() {
		Convey("lambdaCtx", func() {
			Convey("GetRequestObject", func() {
				Convey("Should set the out parameter correctly.", func() {
					type test struct{ v string }
					rv := reflect.ValueOf(&test{v: "test"})
					c := lambdaCtx{req: &rv}

					out := new(test)
					err := c.GetRequestObject(out)

					So(err, ShouldBeNil)
					So(out.v, ShouldEqual, "test")
				})
				Convey("Should set the out parameter correctly if there is no req value.", func() {
					type test struct{ v string }
					c := lambdaCtx{}

					out := new(test)
					err := c.GetRequestObject(out)

					So(err, ShouldBeNil)
					So(*out, ShouldResemble, reflect.Zero(reflect.TypeOf(test{})).Interface())
				})
				Convey("Should validate the out parameter.", func() {
					type test struct{ v string }
					rv := reflect.ValueOf(&test{v: "test"})
					c := lambdaCtx{req: &rv}

					out := test{}
					err := c.GetRequestObject(out)

					So(err, ShouldBeError, "the out parameter is not a pointer")
				})
			})
		})
	})
}
