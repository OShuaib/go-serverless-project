package main

import (
	"os"
	"github.com/OShuaib/go-serverless/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

)

var (
	dbClient dynamodbiface.DynamoDBAPI
)

func main(){
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return
	}
	dbClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}


const tableName = "go-serverless"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error){
	switch req.HTTPMethod{
	case "GET":
		return handlers.GetUser(req, tableName, dbClient)
	case "POST":
		return handlers.CreateUser(req, tableName, dbClient)
	case "PUT":
		return handlers.UpdateUser(req, tableName, dbClient)
	case "DELETE":
		return handlers.DeleteUser(req, tableName, dbClient)
	default:
		return handlers.UnhandleMethod()
	}
}