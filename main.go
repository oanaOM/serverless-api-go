package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type user struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Hobbies string `json:"hobbies"`
}

func show() (*user, error) {
	//Step1
	// u := &user{
	//     ID:      "1",
	//     Name:    "Jamie Oliver",
	//     Hobbies: "Cooking",
	// }
	//Step2
	u, err := getItem("3")
	if err != nil {
		return nil, err
	}
	fmt.Printf("u: %v", u)
	return u, nil
}

func main() {
	lambda.Start(show)
}
