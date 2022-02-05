package user

import (
	"encoding/json"
	"errors"

	"github.com/OShuaib/go-serverless/pkg/validators"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorFailedToFetchRecord ="failed to fetch record"
	ErrorFailedUnmarshalRecord ="failed to unmarshal record"
	ErrorInvalidUserData = "invalid user data"
	ErrorInvalidEmail = "invalid email"
	ErrorCouldNotMarshalItem = "could not marshal item"
	ErrorCouldNotDeleteItem = "could not delete item"
	ErrorCouldNotDynamoPutItem = "could not dynamo put item"
	ErrorUserAlreadyExists = "user already exists"
	ErrorUserDoesNotExist = "user does not exist"
)

type User struct {
	Email 		string	`json:"email"`
	FirstName	string	`json:"first_nmae"`
	LastName	string	`json:"last_name"`
}

func FetchUser(email, tableName string, dbClient dynamodbiface.DynamoDBAPI)(*User, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := dbClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	 item :=new(User)
	 err = dynamodbattribute.UnmarshalMap(result.Item, item)
	 if err != nil {
		 return nil, errors.New(ErrorFailedUnmarshalRecord)
	 }

	 return item, nil
}

func FetchUsers(tableName string, dbClient dynamodbiface.DynamoDBAPI)(*[]User,error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dbClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]User)
	err = dynamodbattribute.UnmarshalMap(result.Items, item)
	return item, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dbClient dynamodbiface.DynamoDBAPI)(*User, error) {
	var u User

	if err := json.Unmarshal([]byte(req.Body),&u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(u.Email){
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ :=  FetchUser(u.Email, tableName, dbClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	user, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: user,
		TableName: aws.String(tableName),
	}
	_, err = dbClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &u,nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dbClient dynamodbiface.DynamoDBAPI)( *User, error) {

}

func DeleteUser() error{

}