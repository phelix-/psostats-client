package userdb

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	PlayersTable    = "players"
	GcToPlayerTable = "gc_to_player"
)

type User struct {
	Id       string
	Gcs      []string
	Password string
	Admin    bool
}

type UserDb interface {
	GetUser(userName string) (*User, error)
	CreateUser(user User) error
	AddGcToUser(userName, guildCard string) error
	GetUsernameByGc(gc string) (string, error)
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
	if len(user.Gcs) < 1 {
		return errors.New("user must have at least one assigned guild card")
	}
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

	for _, gc := range user.Gcs {
		userNameAttributeValue := dynamodb.AttributeValue{S: aws.String(user.Id)}
		guildCardAttributeValue := dynamodb.AttributeValue{S: aws.String(gc)}

		gcToPlayerItem := map[string]*dynamodb.AttributeValue{
			"Gc":     &guildCardAttributeValue,
			"Player": &userNameAttributeValue,
		}
		gcToPlayerUpdate := &dynamodb.PutItemInput{
			Item:      gcToPlayerItem,
			TableName: aws.String(GcToPlayerTable),
		}
		_, err := d.dynamoClient.PutItem(gcToPlayerUpdate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d DynamoUserDb) AddGcToUser(userName, guildCard string) error {
	userNameAttributeValue := dynamodb.AttributeValue{S: aws.String(userName)}
	guildCardAttributeValue := dynamodb.AttributeValue{S: aws.String(guildCard)}

	gcToPlayerItem := map[string]*dynamodb.AttributeValue{
		"Gc":     &guildCardAttributeValue,
		"Player": &userNameAttributeValue,
	}
	gcToPlayerUpdate := &dynamodb.PutItemInput{
		Item:      gcToPlayerItem,
		TableName: aws.String(GcToPlayerTable),
	}
	_, err := d.dynamoClient.PutItem(gcToPlayerUpdate)
	if err != nil {
		return err
	}
	updateItemInput := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":0": {L: []*dynamodb.AttributeValue{&guildCardAttributeValue}},
		},
		ExpressionAttributeNames: map[string]*string{"#0": aws.String("Gcs")},
		Key:                      map[string]*dynamodb.AttributeValue{"Id": &userNameAttributeValue},
		TableName:                aws.String(PlayersTable),

		UpdateExpression: aws.String("SET #0 = list_append(#0, :0)"),
	}
	_, err = d.dynamoClient.UpdateItem(updateItemInput)
	return err
}

func (d DynamoUserDb) GetUsernameByGc(gc string) (string, error) {
	gcAttributeValue := dynamodb.AttributeValue{S: aws.String(gc)}
	getItem := dynamodb.GetItemInput{
		TableName:       aws.String(GcToPlayerTable),
		AttributesToGet: aws.StringSlice([]string{"Player"}),
		Key:             map[string]*dynamodb.AttributeValue{"Gc": &gcAttributeValue},
	}
	item, err := d.dynamoClient.GetItem(&getItem)
	if err != nil {
		return "", err
	}
	playerName := ""
	if len(item.Item) > 0 {
		playerName = *item.Item["Player"].S
	}
	return playerName, nil
}
