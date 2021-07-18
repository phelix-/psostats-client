package userdb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	PlayersTable    = "players"
)

type User struct {
	Id       string
	Password string
	Admin    bool
}

type UserDb interface {
	GetUser(userName string) (*User, error)
	CreateUser(user User) error
}

type DynamoUserDb struct {
	dynamoClient *dynamodb.DynamoDB
}

func DynamoInstance(dynamoClient *dynamodb.DynamoDB) DynamoUserDb {
	return DynamoUserDb{dynamoClient: dynamoClient}
}

func (d DynamoUserDb) GetUser(userName string) (*User, error) {
	user := User{}
	primaryKey := dynamodb.AttributeValue{
		S: aws.String(userName),
	}
	getItem := dynamodb.GetItemInput{
		TableName: aws.String(PlayersTable),
		Key:       map[string]*dynamodb.AttributeValue{"Id": &primaryKey},
	}
	item, err := d.dynamoClient.GetItem(&getItem)
	if err != nil || item.Item == nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalMap(item.Item, &user)
	return &user, err
}

func (d DynamoUserDb) CreateUser(user User) error {
	marshalled, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}
	userInput := &dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String(PlayersTable),
	}
	_, err = d.dynamoClient.PutItem(userInput)
	if err != nil {
		return err
	}

	return nil
}
