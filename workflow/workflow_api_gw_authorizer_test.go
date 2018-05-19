package workflow

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPIGWAuthorizerWorkflow(t *testing.T) {
	Convey("API Gateway Authorizer workflow", t, func() {
		Convey("Should handle raw response.", func() {
			w := NewAPIGWAuthorizerWorkflowBuilder().
				SetHandler(func(ctx Context, evt events.APIGatewayCustomAuthorizerRequest) error {
					ctx.SetRawResponse(getAuthorizerResponse("test"))
					return nil
				}).
				Build()

			res, err := w.GetLambdaHandler()(context.TODO(), events.APIGatewayCustomAuthorizerRequest{})

			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			So(*res, ShouldResemble, getAuthorizerResponse("test"))
		})

		Convey("Should handle response.", func() {
			w := NewAPIGWAuthorizerWorkflowBuilder().
				SetHandler(func(ctx Context, evt events.APIGatewayCustomAuthorizerRequest) error {
					ctx.SetResponse(getAuthorizerResponse("test"))
					return nil
				}).
				Build()

			res, err := w.GetLambdaHandler()(context.TODO(), events.APIGatewayCustomAuthorizerRequest{})

			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			So(*res, ShouldResemble, getAuthorizerResponse("test"))
		})

		Convey("Should handle raw response before response.", func() {
			w := NewAPIGWAuthorizerWorkflowBuilder().
				SetHandler(func(ctx Context, evt events.APIGatewayCustomAuthorizerRequest) error {
					ctx.SetRawResponse(getAuthorizerResponse("raw"))
					ctx.SetResponse(getAuthorizerResponse("response"))
					return nil
				}).
				Build()

			res, err := w.GetLambdaHandler()(context.TODO(), events.APIGatewayCustomAuthorizerRequest{})

			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			So(*res, ShouldResemble, getAuthorizerResponse("raw"))
		})

		Convey("Should handle nil response.", func() {
			w := NewAPIGWAuthorizerWorkflowBuilder().
				SetHandler(func(ctx Context, evt events.APIGatewayCustomAuthorizerRequest) error {
					ctx.SetResponse(nil)
					return nil
				}).
				Build()

			res, err := w.GetLambdaHandler()(context.TODO(), events.APIGatewayCustomAuthorizerRequest{})

			So(err, ShouldBeNil)
			So(res, ShouldBeNil)
		})

		Convey("Should handle invalid response.", func() {
			w := NewAPIGWAuthorizerWorkflowBuilder().
				SetHandler(func(ctx Context, evt events.APIGatewayCustomAuthorizerRequest) error {
					ctx.SetResponse(make(map[string]string))
					return nil
				}).
				Build()

			_, err := w.GetLambdaHandler()(context.TODO(), events.APIGatewayCustomAuthorizerRequest{})

			So(err, ShouldBeError, "invalid response")
		})
	})
}

func getAuthorizerResponse(principal string) events.APIGatewayCustomAuthorizerResponse {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: principal,
	}
}
