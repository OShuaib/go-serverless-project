package handlers

import (
	"net/http"

	"github.com/OShuaib/go-serverless/pkg/user"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var ErrMethodNotAllowed = "method not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error, omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, tableName string, dbClient dynamodbiface.DynamoDBAPI)(*events.APIGatewayProxyResponse, error){
	
	email := req.QueryStringParameters["email"]

	if len(email) > 0 {
		result, err := user.FetchUser(email, tableName, dbClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}
	result, err := user.FetchUsers(tableName, dbClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)
}
