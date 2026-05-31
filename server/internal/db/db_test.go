package db_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestGetQuestRecordsForQuest(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		t.Fatal(err)
	}
	dynamoClient := dynamodb.New(sess)
	if err = CreateAllTables(dynamoClient); err != nil {
		t.Fatal(err)
	}

	const quest = "Sweep-up Operation #1"
	writeTestQuestRecord(t, dynamoClient, quest, "1n", "game-1")
	writeTestQuestRecord(t, dynamoClient, quest, "2n", "game-2")
	writeTestQuestRecord(t, dynamoClient, "Some Other Quest", "1n", "game-3")

	records, err := db.GetQuestRecordsForQuest(quest, dynamoClient)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records for quest, got %d", len(records))
	}
	for _, r := range records {
		if r.Quest != quest {
			t.Errorf("expected quest %q, got %q", quest, r.Quest)
		}
	}
}

func TestGetAllGamesForPlayerQuest(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		t.Fatal(err)
	}
	dynamoClient := dynamodb.New(sess)
	if err = CreateAllTables(dynamoClient); err != nil {
		t.Fatal(err)
	}

	const player = "testplayer"
	const targetQuest = "Sweep-up Operation #1"

	older := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

	writeTestPlayerGame(t, dynamoClient, player, targetQuest, "pg-1", 1, older)
	writeTestPlayerGame(t, dynamoClient, player, targetQuest, "pg-2", 2, newer)
	writeTestPlayerGame(t, dynamoClient, player, "Lost WORKS Machine", "pg-3", 3, older)

	games, err := db.GetAllGamesForPlayerQuest(player, targetQuest, dynamoClient)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 2 {
		t.Fatalf("expected 2 games for player+quest, got %d", len(games))
	}
	for _, g := range games {
		if g.Quest != targetQuest {
			t.Errorf("expected quest %q, got %q", targetQuest, g.Quest)
		}
	}
	if !games[0].Timestamp.After(games[1].Timestamp) {
		t.Errorf("expected games sorted descending by Timestamp: games[0]=%v games[1]=%v",
			games[0].Timestamp, games[1].Timestamp)
	}
}

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
	returned, err := db.GetGame(id, 0, dynamoClient)
	if err != nil {
		t.Error(err)
	}
	if returned == nil || returned.QuestName != questRun.QuestName {
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
	if _, exists := tables[db.QuestRecordsTable]; !exists {
		if err = CreateQuestRecordsTable(dynamoClient); err != nil {
			return err
		}
	}
	if _, exists := tables[db.RecentGamesByPlayerTable]; !exists {
		if err = CreateRecentGamesByPlayerTable(dynamoClient); err != nil {
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

func CreateQuestRecordsTable(dynamoClient *dynamodb.DynamoDB) error {
	_, err := dynamoClient.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(db.QuestRecordsTable),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("Quest"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("Category"), AttributeType: aws.String("S")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("Quest"), KeyType: aws.String(dynamodb.KeyTypeHash)},
			{AttributeName: aws.String("Category"), KeyType: aws.String(dynamodb.KeyTypeRange)},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	return err
}

func CreateRecentGamesByPlayerTable(dynamoClient *dynamodb.DynamoDB) error {
	_, err := dynamoClient.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(db.RecentGamesByPlayerTable),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("Player"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("IdInt"), AttributeType: aws.String("N")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("Player"), KeyType: aws.String(dynamodb.KeyTypeHash)},
			{AttributeName: aws.String("IdInt"), KeyType: aws.String(dynamodb.KeyTypeRange)},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	return err
}

func writeTestPlayerGame(t *testing.T, dynamoClient *dynamodb.DynamoDB, player, questName, gameId string, idInt int, ts time.Time) {
	t.Helper()
	_, err := dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(db.RecentGamesByPlayerTable),
		Item: map[string]*dynamodb.AttributeValue{
			"Player":    {S: aws.String(player)},
			"IdInt":     {N: aws.String(strconv.Itoa(idInt))},
			"Id":        {S: aws.String(gameId)},
			"Quest":     {S: aws.String(questName)},
			"Category":  {S: aws.String("2n")},
			"Episode":   {N: aws.String("1")},
			"Timestamp": {S: aws.String(ts.Format(time.RFC3339))},
		},
	})
	if err != nil {
		t.Fatalf("failed to write test player game: %v", err)
	}
}

func writeTestQuestRecord(t *testing.T, dynamoClient *dynamodb.DynamoDB, questName, category, gameId string) {
	t.Helper()
	_, err := dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(db.QuestRecordsTable),
		Item: map[string]*dynamodb.AttributeValue{
			"Quest":    {S: aws.String(questName)},
			"Category": {S: aws.String(category)},
			"Id":       {S: aws.String(gameId)},
		},
	})
	if err != nil {
		t.Fatalf("failed to write test quest record: %v", err)
	}
}
