package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var idRegexp = regexp.MustCompile(`[0-9]`)
var errorLogger = log.New(os.Stderr, "Error:", log.Llongfile)

type user struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Hobbies string `json:"hobbies"`
}

//router switch on the HTTP request method to determin whcih action to takes
func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return show(req)
	case "POST":
		return create(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

//create adds a new user to the db
func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	u := new(user)

	err := json.Unmarshal([]byte(req.Body), u)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}
	if !idRegexp.MatchString(u.ID) {
		return clientError(http.StatusBadRequest)
	}
	if u.Name == "" || u.Hobbies == "" {
		return clientError(http.StatusBadRequest)
	}

	err = putItem(u)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": fmt.Sprintf("/users?id=%s", u.ID)},
	}, nil

}

//show retrieves from the db the user
func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	id := req.QueryStringParameters["id"]
	if !idRegexp.MatchString(id) {
		return clientError(http.StatusBadRequest)
	}
	//Fetch the user recored from the db based on the ID value
	u, err := getItem(id)
	if err != nil {
		return serverError(err)
	}
	//the APIGatewayProxyResponse.Body field needs to be a string, so
	//we marshal the users record into a JSON
	js, err := json.Marshal(u)
	if err != nil {
		return serverError(err)
	}
	fmt.Printf("u: %v", u)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

//serverError is a helper for handling errors. This logs any error to os.Stderr and returns a 500 Internal Server Error
//response that the AWS API Gateway understands
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(show)
}
