package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/TsvetanMilanov/go-lambda-workflow/workflow"
	. "github.com/smartystreets/goconvey/convey"
)

type JSONReq struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func TestAPIGWProxyWorkflow(t *testing.T) {
	Convey("API Gateway proxy workflow", t, func() {
		input := JSONReq{
			Message: "Hello World!",
			Code:    123,
		}

		apigwReq := getAPIGWProxyRequest(http.MethodGet, "/", input)

		Convey("Should return API Gateway proxy response with JSON request body.", func() {
			handler := func(c workflow.Context, req JSONReq) error {
				So(req, ShouldResemble, input)
				c.SetResponse(req).SetResponseStatusCode(req.Code)
				return nil
			}

			w := workflow.NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, input.Code)
			So(res.Body, ShouldEqual, getStringBody(input))
		})

		Convey("Should return API Gateway proxy response with JSON array request body.", func() {
			input := []JSONReq{
				{
					Message: "Hello World!",
					Code:    123,
				},
				{
					Message: "Second",
					Code:    456,
				},
			}

			apigwReq := getAPIGWProxyRequest(http.MethodGet, "/", input)

			handler := func(c workflow.Context, req []JSONReq) error {
				So(req, ShouldResemble, input)
				c.SetResponse(req).SetResponseStatusCode(http.StatusOK)
				return nil
			}

			w := workflow.NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusOK)
			So(res.Body, ShouldEqual, getStringBody(input))
		})

		Convey("Should return API Gateway proxy response with JSON request body and pointer input.", func() {
			ptrHandler := func(c workflow.Context, req *JSONReq) error {
				So(*req, ShouldResemble, input)
				c.SetResponse(req).SetResponseStatusCode(req.Code)
				return nil
			}

			w := workflow.NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", ptrHandler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, input.Code)
			So(res.Body, ShouldEqual, getStringBody(input))
		})

		Convey("Should return API Gateway proxy response with JSON request body and no input parameter in the handler.", func() {
			handler := func(c workflow.Context) error {
				c.SetResponseStatusCode(http.StatusOK)
				return nil
			}

			w := workflow.NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusOK)
			So(res.Body, ShouldEqual, "")
		})

		Convey("Should return not found API Gateway proxy response when no handler was found.", func() {
			handler := func(c workflow.Context) error {
				c.SetResponseStatusCode(http.StatusOK)
				return nil
			}

			w := workflow.NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/no-such-path", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusNotFound)
			So(res.Body, ShouldEqual, "")
		})

		Convey("Should return API Gateway proxy response with string request body and string handler input.", func() {
			input := "test"

			apigwReq := getAPIGWProxyRequest(http.MethodGet, "/", input)

			handler := func(c workflow.Context, input string) error {
				c.SetResponse(input).SetResponseStatusCode(http.StatusOK)
				return nil
			}

			w := workflow.NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusOK)
			So(res.Body, ShouldEqual, fmt.Sprintf(`"%s"`, input))
		})

		Convey("Should return API Gateway proxy response with string request body and []byte handler input.", func() {
			input := "test"

			apigwReq := getAPIGWProxyRequest(http.MethodGet, "/", input)

			handler := func(c workflow.Context, input []byte) error {
				c.SetResponse(input).SetResponseStatusCode(http.StatusOK)
				return nil
			}

			w := workflow.NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusOK)
			So(res.Body, ShouldEqual, fmt.Sprintf(`"%s"`, input))
		})

		Convey("Should return API Gateway proxy response with raw response set in the handler context.", func() {
			rawRes := events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "test"}
			handler := func(c workflow.Context) error {
				c.SetRawResponse(rawRes)
				return nil
			}

			w := workflow.NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(*res, ShouldResemble, rawRes)

			handler = func(c workflow.Context) error {
				c.SetRawResponse(&rawRes)
				return nil
			}

			w = workflow.NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err = w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(*res, ShouldResemble, rawRes)
		})
	})
}

func getStringBody(input interface{}) string {
	bodyBytes, _ := json.Marshal(input)
	return string(bodyBytes)
}

func getAPIGWProxyRequest(method, path string, input interface{}) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{
		Path:       path,
		Body:       getStringBody(input),
		HTTPMethod: method,
	}
}
