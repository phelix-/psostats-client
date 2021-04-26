package db_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"log"
	"testing"
	"time"
)

func TestGetPlayerRecentGames(t *testing.T) {

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return
	}
	dynamoClient := dynamodb.New(sess)
	if err = CreateAllTables(dynamoClient); err != nil {
		t.Error(err)
	}

	games, err := db.GetRecentGames(dynamoClient)
	for _, game := range games {
		log.Printf("%v - %v", game.Id, game.Quest)
	}
	questRun := createMockGame()
	id, err := db.WriteGameById(questRun, dynamoClient)
	if err != nil {
		t.Error(err)
	}
	returned, err := db.GetGame(id, dynamoClient)
	if err != nil {
		t.Error(err)
	}
	if returned.QuestName != questRun.QuestName {
		t.Fatalf("Quest names didn't match")
	}
}

func createMockGame() *model.QuestRun {
	startDate := time.Date(2021, time.April, 24, 15, 19, 0, 0, time.Local)
	questRun := model.QuestRun{
		Server:      "unseen",
		PlayerName:  "bvelix",
		PlayerClass: "HUcast",
		GuildCard:   "phelix",
		AllPlayers: []model.BasePlayerInfo{
			{Name: "bvelix", GuildCard: "42", Level: 200, Class: "HUcast"},
			{Name: "player2", GuildCard: "43", Level: 200, Class: "HUnewearl"},
		},
		Id:                  "",
		Difficulty:          "Ultimate",
		Episode:             1,
		QuestName:           "Sweep-up Operation #1",
		QuestComplete:       true,
		QuestStartTime:      startDate,
		QuestEndTime:        startDate.Add(time.Minute*5 + time.Second*20),
		QuestDuration:       (time.Minute*5 + time.Second*20).String(),
		DeathCount:          2,
		HP:                  nil,
		TP:                  nil,
		Meseta:              nil,
		MesetaCharged:       nil,
		Room:                nil,
		IllegalShifta:       false,
		PbCategory:          false,
		ShiftaLvl:           nil,
		DebandLvl:           nil,
		Invincible:          nil,
		Events:              nil,
		Monsters:            nil,
		MonsterCount:        nil,
		MonstersKilledCount: nil,
		MonstersDead:        0,
		WeaponsUsed:         nil,
		FreezeTraps:         nil,
		FTUsed:              12,
		DTUsed:              3,
		CTUsed:              1,
		TPUsed:              0,
	}
	return &questRun
}

func CreateAllTables(dynamoClient *dynamodb.DynamoDB) error {
	result, err := dynamoClient.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return err
	}
	tables := make(map[string]bool)
	for _, tableName := range result.TableNames {
		tables[*tableName] = true
	}
	if _, exists := tables["games_by_id"]; !exists {
		if err = CreateGamesById(dynamoClient); err != nil {
			return err
		}
	}
	if _, exists := tables["games_counter"]; !exists {
		if err = CreateGamesCounter(dynamoClient); err != nil {
			return err
		}
	}
	if _, exists := tables[db.RecentGamesByMonth]; !exists {
		if err = CreateRecentGamesByMonth(dynamoClient); err != nil {
			return err
		}
	}
	return nil
}

func CreateGamesById(dynamoClient *dynamodb.DynamoDB) error {
	attributeDefinition := dynamodb.AttributeDefinition{
		AttributeName: aws.String("Id"),
		AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
	}
	keySchemaElement := dynamodb.KeySchemaElement{
		AttributeName: aws.String("Id"),
		KeyType:       aws.String(dynamodb.KeyTypeHash),
	}
	provisionedThroughput := dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(1),
		WriteCapacityUnits: aws.Int64(1),
	}
	createTableInput := dynamodb.CreateTableInput{
		AttributeDefinitions:  []*dynamodb.AttributeDefinition{&attributeDefinition},
		KeySchema:             []*dynamodb.KeySchemaElement{&keySchemaElement},
		TableName:             aws.String("games_by_id"),
		ProvisionedThroughput: &provisionedThroughput,
	}
	_, err := dynamoClient.CreateTable(&createTableInput)
	return err
}

func CreateGamesCounter(dynamoClient *dynamodb.DynamoDB) error {
	provisionedThroughput := dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(1),
		WriteCapacityUnits: aws.Int64(1),
	}
	attributeDefinition := dynamodb.AttributeDefinition{
		AttributeName: aws.String("key"),
		AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
	}
	keySchemaElement := dynamodb.KeySchemaElement{
		AttributeName: aws.String("key"),
		KeyType:       aws.String(dynamodb.KeyTypeHash),
	}
	createTableInput := dynamodb.CreateTableInput{
		AttributeDefinitions:  []*dynamodb.AttributeDefinition{&attributeDefinition},
		KeySchema:             []*dynamodb.KeySchemaElement{&keySchemaElement},
		TableName:             aws.String("games_counter"),
		ProvisionedThroughput: &provisionedThroughput,
	}
	_, err := dynamoClient.CreateTable(&createTableInput)
	return err
}

func CreateRecentGamesByMonth(dynamoClient *dynamodb.DynamoDB) error {
	provisionedThroughput := dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(1),
		WriteCapacityUnits: aws.Int64(1),
	}
	pk := dynamodb.AttributeDefinition{
		AttributeName: aws.String("Month"),
		AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
	}
	pkSchema := dynamodb.KeySchemaElement{
		AttributeName: aws.String("Month"),
		KeyType:       aws.String(dynamodb.KeyTypeHash),
	}
	sortKey := dynamodb.AttributeDefinition{
		AttributeName: aws.String("Id"),
		AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
	}
	sortKeySchema := dynamodb.KeySchemaElement{
		AttributeName: aws.String("Id"),
		KeyType:       aws.String(dynamodb.KeyTypeRange),
	}
	createTableInput := dynamodb.CreateTableInput{
		AttributeDefinitions:  []*dynamodb.AttributeDefinition{&pk, &sortKey},
		KeySchema:             []*dynamodb.KeySchemaElement{&pkSchema, &sortKeySchema},
		TableName:             aws.String(db.RecentGamesByMonth),
		ProvisionedThroughput: &provisionedThroughput,
	}
	_, err := dynamoClient.CreateTable(&createTableInput)
	return err
}
