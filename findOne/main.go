package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Movie struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func findOne(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//id := request.PathParameters["id"]
	id := request.QueryStringParameters["id"]

	if len(id) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Required input ID not found",
		}, nil
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving AWS credentials",
		}, nil
	}

	svc := dynamodb.New(cfg)
	req := svc.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]dynamodb.AttributeValue{
			"ID": dynamodb.AttributeValue{
				S: aws.String(id),
			},
		},
	})
	res, err := req.Send()
	if err != nil {
		var s string
		s = fmt.Sprintf("Error while fetching movie ID:[%s] Err:[%s]", id, err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       s,
		}, nil
	}

	movie := Movie{
		ID:   *res.Item["ID"].S,
		Name: *res.Item["Name"].S,
	}

	response, err := json.Marshal(movie)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding to string value",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}, nil
}

func main() {
	lambda.Start(findOne)
}
