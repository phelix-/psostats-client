package userdb_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/phelix-/psostats/v2/server/internal/userdb"
	"log"
	"math/rand"
	"reflect"
	"testing"
)

func getFixture() (*userdb.DynamoUserDb, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	dynamoClient := dynamodb.New(sess)
	result, err := dynamoClient.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return nil, err
	}
	tables := make(map[string]bool)
	for _, tableName := range result.TableNames {
		tables[*tableName] = true
	}
	if _, exists := tables[userdb.GcToPlayerTable]; !exists {
		err = CreateGcToPlayerTable(dynamoClient)
		if err != nil {
			return nil, err
		}
	}
	if _, exists := tables[userdb.PlayersTable]; !exists {
		err = CreatePlayersTable(dynamoClient)
		if err != nil {
			return nil, err
		}
	}
	instance := userdb.DynamoInstance(dynamoClient)
	return &instance, nil
}

func TestAwsUserDb_TryEverythingOnce(t *testing.T) {
	fixture, err := getFixture()
	if err != nil {
		t.Error(err)
	}
	user := userdb.User{
		Id:       fmt.Sprintf("test%v", rand.Int()),
		Gcs:      []string{"gc1"},
		Password: "password",
		Admin:    false,
	}
	err = fixture.CreateUser(user)
	if err != nil {
		t.Error(err)
	}
	fromDb, err := fixture.GetUser(user.Id)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(user, *fromDb) {
		t.Error("user returned from db was modified")
	}

	id1 := fmt.Sprintf("test+%v", rand.Int())
	id2 := fmt.Sprintf("test+%v", rand.Int())
	err = fixture.AddGcToUser(user.Id, id1)
	if err != nil {
		t.Error(err)
	}
	err = fixture.AddGcToUser(user.Id, id2)
	if err != nil {
		t.Error(err)
	}
	userName, err := fixture.GetUsernameByGc(id1)
	if err != nil {
		t.Error(err)
	}
	if userName != user.Id {
		t.Error("username by gc")
	}
	fromDb, err = fixture.GetUser(user.Id)
	if err != nil {
		t.Error(err)
	}
	expected := []string{"gc1", id1, id2}
	if !reflect.DeepEqual(expected, fromDb.Gcs) {
		t.Error("failed to map gcs correctly")
	}
}

func CreatePlayersTable(dynamoClient *dynamodb.DynamoDB) error {
	provisionedThroughput := dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(1),
		WriteCapacityUnits: aws.Int64(1),
	}
	attributeDefinition := dynamodb.AttributeDefinition{
		AttributeName: aws.String("Id"),
		AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
	}
	keySchemaElement := dynamodb.KeySchemaElement{
		AttributeName: aws.String("Id"),
		KeyType:       aws.String(dynamodb.KeyTypeHash),
	}
	createTableInput := dynamodb.CreateTableInput{
		AttributeDefinitions:  []*dynamodb.AttributeDefinition{&attributeDefinition},
		KeySchema:             []*dynamodb.KeySchemaElement{&keySchemaElement},
		TableName:             aws.String(userdb.PlayersTable),
		ProvisionedThroughput: &provisionedThroughput,
	}
	_, err := dynamoClient.CreateTable(&createTableInput)
	return err
}

func CreateGcToPlayerTable(dynamoClient *dynamodb.DynamoDB) error {
	provisionedThroughput := dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(1),
		WriteCapacityUnits: aws.Int64(1),
	}
	attributeDefinition := dynamodb.AttributeDefinition{
		AttributeName: aws.String("Gc"),
		AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
	}
	keySchemaElement := dynamodb.KeySchemaElement{
		AttributeName: aws.String("Gc"),
		KeyType:       aws.String(dynamodb.KeyTypeHash),
	}
	createTableInput := dynamodb.CreateTableInput{
		AttributeDefinitions:  []*dynamodb.AttributeDefinition{&attributeDefinition},
		KeySchema:             []*dynamodb.KeySchemaElement{&keySchemaElement},
		TableName:             aws.String(userdb.GcToPlayerTable),
		ProvisionedThroughput: &provisionedThroughput,
	}
	_, err := dynamoClient.CreateTable(&createTableInput)
	return err
}
