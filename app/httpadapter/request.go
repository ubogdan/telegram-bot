package httpadapter

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// EventToRequest converts an API Gateway proxy event into an http.Request object.
// Returns the populated request maintaining headers
func EventToRequest(req events.LambdaFunctionURLRequest) (*http.Request, error) {
	decodedBody := []byte(req.Body)
	if req.IsBase64Encoded {
		base64Body, err := base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return nil, err
		}
		decodedBody = base64Body
	}

	path := req.RequestContext.HTTP.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	serverAddress := "https://" + req.RequestContext.DomainName

	path = serverAddress + path

	if len(req.QueryStringParameters) > 0 {
		// Support `QueryStringParameters` for backward compatibility.
		// https://github.com/awslabs/aws-lambda-go-api-proxy/issues/37
		queryString := ""
		for q := range req.QueryStringParameters {
			if queryString != "" {
				queryString += "&"
			}
			queryString += url.QueryEscape(q) + "=" + url.QueryEscape(req.QueryStringParameters[q])
		}
		path += "?" + queryString
	}

	httpRequest, err := http.NewRequest(
		strings.ToUpper(req.RequestContext.HTTP.Method),
		path,
		bytes.NewReader(decodedBody),
	)

	if err != nil {
		fmt.Printf("Could not convert request %s:%s to http.Request\n", req.RequestContext.HTTP.Method, req.RequestContext.HTTP.Path)
		log.Println(err)
		return nil, err
	}

	httpRequest.RemoteAddr = req.RequestContext.HTTP.SourceIP

	for h := range req.Headers {
		httpRequest.Header.Add(h, req.Headers[h])
	}

	httpRequest.RequestURI = httpRequest.URL.RequestURI()

	return httpRequest, nil
}
