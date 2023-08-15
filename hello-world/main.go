package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("no IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("non 200 response found")
)

type Greeting struct {
	Message string `json:"message"`
	Ip      string `json:"ip"`
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayV2HTTPResponse{}, ErrNon200Response
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayV2HTTPResponse{}, ErrNoIP
	}

	name := request.QueryStringParameters["name"]
	if name == "" {
		name = "World"
	}

	greeting := Greeting{
		Message: fmt.Sprintf("Hello, %s", name),
		Ip:      string(ip),
	}

	str, err := json.Marshal(greeting)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	return events.APIGatewayV2HTTPResponse{
		Body:       string(str),
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
