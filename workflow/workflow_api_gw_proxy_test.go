package workflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
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
			handler := func(c Context, req JSONReq) error {
				So(req, ShouldResemble, input)
				c.SetResponse(req).SetResponseStatusCode(req.Code)
				return nil
			}

			w := NewAPIGWProxyWorkflowBuilder().
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

			handler := func(c Context, req []JSONReq) error {
				So(req, ShouldResemble, input)
				c.SetResponse(req).SetResponseStatusCode(http.StatusOK)
				return nil
			}

			w := NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusOK)
			So(res.Body, ShouldEqual, getStringBody(input))
		})

		Convey("Should return API Gateway proxy response with JSON request body and pointer input.", func() {
			ptrHandler := func(c Context, req *JSONReq) error {
				So(*req, ShouldResemble, input)
				c.SetResponse(req).SetResponseStatusCode(req.Code)
				return nil
			}

			w := NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", ptrHandler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, input.Code)
			So(res.Body, ShouldEqual, getStringBody(input))
		})

		Convey("Should return API Gateway proxy response with JSON request body and no input parameter in the handler.", func() {
			handler := func(c Context) error {
				c.SetResponseStatusCode(http.StatusOK)
				return nil
			}

			w := NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusOK)
			So(res.Body, ShouldEqual, "")
		})

		Convey("Should return not found API Gateway proxy response when no handler was found.", func() {
			handler := func(c Context) error {
				c.SetResponseStatusCode(http.StatusOK)
				return nil
			}

			w := NewAPIGWProxyWorkflowBuilder().
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

			handler := func(c Context, input string) error {
				c.SetResponse(input).SetResponseStatusCode(http.StatusOK)
				return nil
			}

			w := NewAPIGWProxyWorkflowBuilder().
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

			handler := func(c Context, input []byte) error {
				c.SetResponse(input).SetResponseStatusCode(http.StatusOK)
				return nil
			}

			w := NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusOK)
			So(res.Body, ShouldEqual, fmt.Sprintf(`"%s"`, input))
		})

		Convey("Should return API Gateway proxy response with raw response set in the handler context.", func() {
			rawRes := events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "test"}
			handler := func(c Context) error {
				c.SetRawResponse(rawRes)
				return nil
			}

			w := NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(*res, ShouldResemble, rawRes)

			handler = func(c Context) error {
				c.SetRawResponse(&rawRes)
				return nil
			}

			w = NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).
				Build()

			res, err = w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(*res, ShouldResemble, rawRes)
		})

		Convey("Should execute pre actions before the handler.", func() {
			flow := ""
			handler := func(c Context) error {
				flow += "handler"
				return nil
			}
			a1 := func(c Context) error {
				flow += "pre1"
				return nil
			}
			a2 := func(c Context) error {
				flow += "pre2"
				return nil
			}

			w := NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).WithPreActions(a1, a2).
				Build()

			_, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(flow, ShouldEqual, "pre1pre2handler")
		})

		Convey("Should execute post actions before the handler.", func() {
			flow := ""
			handler := func(c Context) error {
				flow += "handler"
				return nil
			}
			a1 := func(c Context) error {
				flow += "post1"
				return nil
			}
			a2 := func(c Context) error {
				flow += "post2"
				return nil
			}

			w := NewAPIGWProxyWorkflowBuilder().
				AddGetHandler("/", handler).WithPostActions(a1, a2).
				Build()

			_, err := w.GetLambdaHandler()(nil, apigwReq)

			So(err, ShouldBeNil)
			So(flow, ShouldEqual, "handlerpost1post2")
		})

		Convey("Should handle paths correctly", func() {
			type testCase struct {
				testName        string
				cleanPath       string
				specialPath     string
				testCleanPath   string
				testSpecialPath string
			}
			testCases := []testCase{
				{
					testName:        "Should make difference between paths containing / and paths which does not on root level.",
					cleanPath:       "",
					specialPath:     "/",
					testCleanPath:   "",
					testSpecialPath: "/",
				},
				{
					testName:        "Should make difference between paths containing / and paths which does not on non root level.",
					cleanPath:       "/test",
					specialPath:     "/test/",
					testCleanPath:   "/test",
					testSpecialPath: "/test/",
				},
				{
					testName:        "Should handle parameterized paths.",
					cleanPath:       "/",
					specialPath:     "/{id}",
					testCleanPath:   "/",
					testSpecialPath: "/5",
				},
				{
					testName:        "Should handle nested parameterized paths.",
					cleanPath:       "/car/",
					specialPath:     "/car/{id}",
					testCleanPath:   "/car/",
					testSpecialPath: "/car/5",
				},
				{
					testName:        "Should handle nested parameterized paths with multiple parameters.",
					cleanPath:       "/car/6/model/",
					specialPath:     "/car/{id}/model/{number}",
					testCleanPath:   "/car/6/model/",
					testSpecialPath: "/car/5/model/6",
				},
				{
					testName:        "Should choose clean path over special path if there is registered clean path with hard-coded value which is like parameterized one.",
					cleanPath:       "/car/6",
					specialPath:     "/car/{id}",
					testCleanPath:   "/car/6",
					testSpecialPath: "/car/7",
				},
			}

			for _, tc := range testCases {
				Convey(tc.testName, func() {
					specialHandlerCalled := false
					specialHandler := func(Context) error {
						specialHandlerCalled = true
						return nil
					}
					cleanHandlerCalled := false
					cleanHandler := func(Context) error {
						cleanHandlerCalled = true
						return nil
					}
					h := NewAPIGWProxyWorkflowBuilder().
						AddGetHandler(tc.specialPath, specialHandler).
						AddGetHandler(tc.cleanPath, cleanHandler).
						Build().
						GetLambdaHandler()

					_, err := h(nil, getAPIGWProxyRequest(http.MethodGet, tc.testSpecialPath, nil))
					So(err, ShouldBeNil)
					So(specialHandlerCalled, ShouldBeTrue)
					So(cleanHandlerCalled, ShouldBeFalse)

					specialHandlerCalled = false
					cleanHandlerCalled = false

					_, err = h(nil, getAPIGWProxyRequest(http.MethodGet, tc.testCleanPath, nil))
					So(err, ShouldBeNil)
					So(specialHandlerCalled, ShouldBeFalse)
					So(cleanHandlerCalled, ShouldBeTrue)
				})
			}
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
