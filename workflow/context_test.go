package workflow

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
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

			Convey("GetRawResult", func() {
				Convey("Should validate if the out parameter", func() {
					Convey("Is pointer.", func() {
						c := new(lambdaCtx)
						c.rawResponse = &events.APIGatewayProxyResponse{}
						err := c.GetRawResponse(events.APIGatewayProxyResponse{})
						So(err, ShouldBeError, "the out parameter must be a pointer")
					})
					Convey("Can be set.", func() {
						c := new(lambdaCtx)
						c.rawResponse = &events.APIGatewayProxyResponse{}
						out := struct {
							private *events.APIGatewayProxyResponse
						}{}
						err := c.GetRawResponse(out.private)
						So(err, ShouldBeError, "can't set raw response out value")
					})
					Convey("Is the same type as the raw response.", func() {
						c := new(lambdaCtx)
						c.rawResponse = &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
						res := new(events.APIGatewayProxyRequest)
						err := c.GetRawResponse(res)
						So(err, ShouldBeError, "cannot set raw response of type APIGatewayProxyResponse to output of type APIGatewayProxyRequest")
					})
				})

				Convey("Should set the out parameter.", func() {
					c := new(lambdaCtx)
					c.rawResponse = &events.APIGatewayProxyResponse{StatusCode: http.StatusTeapot}
					res := new(events.APIGatewayProxyResponse)
					err := c.GetRawResponse(res)

					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)
					So(res.StatusCode, ShouldEqual, http.StatusTeapot)
				})
			})
		})
	})
}
