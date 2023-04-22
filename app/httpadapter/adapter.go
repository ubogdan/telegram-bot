package httpadapter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func New(handler http.Handler) func(context.Context, events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	return func(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
		req, err := EventToRequest(event)
		if err != nil {
			return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError}, fmt.Errorf("error while converting request.go to http.Request: %v", err)
		}

		w := NewURLResponseWriter()
		handler.ServeHTTP(http.ResponseWriter(w), req)

		resp, err := w.GetProxyResponse()
		if err != nil {
			return events.LambdaFunctionURLResponse{StatusCode: http.StatusGatewayTimeout}, fmt.Errorf("error while generating proxy response: %v", err)
		}

		return resp, nil
	}
}
