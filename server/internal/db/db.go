package db

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/phelix-/psostats/v2/pkg/model"
	"io"
	"log"
	"strconv"
	"time"
)

const (
	gamesTable          = "games_by_id"
	gameCountTable      = "games_counter"
	gameCountPrimaryKey = "game_count"
)

func WriteGame(questRun *model.QuestRun, dynamoClient *dynamodb.DynamoDB) (string, error) {
	gameId, err := incrementAndGetGameId(dynamoClient)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	writer := gzip.NewWriter(buffer)
	jsonQuestBytes, err := json.Marshal(questRun)
	if err != nil {
		return "", err
	}
	_, err = writer.Write(jsonQuestBytes)
	if err != nil {
		return "", err
	}
	err = writer.Flush()
	if err != nil {
		return "", err
	}
	category := fmt.Sprintf("%d", len(questRun.AllPlayers))
	if questRun.PbCategory {
		category += "p"
	} else {
		category += "n"
	}
	duration, err := time.ParseDuration(questRun.QuestDuration)
	if err != nil {
		return "", err
	}

	game := model.Game{
		Id:        fmt.Sprintf("%d", gameId),
		Player:    fmt.Sprintf("%v+%v", questRun.Server, questRun.GuildCard),
		Category:  category,
		Quest:     questRun.QuestName,
		Time:      duration,
		Timestamp: questRun.QuestStartTime,
		Episode:   int(questRun.Episode),
		GameGzip:  buffer.Bytes(),
	}
	marshalled, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		return "", err
	}
	input2 := &dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String(gamesTable),
	}
	_, err = dynamoClient.PutItem(input2)
	if err != nil {
		return "", err
	}
	return game.Id, nil
}

func incrementAndGetGameId(dynamoClient *dynamodb.DynamoDB) (int, error) {
	value := dynamodb.AttributeValue{
		N: aws.String("1"),
	}
	update := dynamodb.AttributeValueUpdate{
		Action: aws.String(dynamodb.AttributeActionAdd),
		Value:  &value,
	}
	primaryKey := dynamodb.AttributeValue{
		S: aws.String(gameCountPrimaryKey),
	}
	updateItemInput := dynamodb.UpdateItemInput{
		TableName:        aws.String(gameCountTable),
		Key:              map[string]*dynamodb.AttributeValue{"key": &primaryKey},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{"count": &update},
		ReturnValues:     aws.String(dynamodb.ReturnValueUpdatedNew),
	}
	item, err := dynamoClient.UpdateItem(&updateItemInput)
	if err != nil {
		return -1, err
	}
	gameIdString := item.Attributes["count"].N
	gameId, err := strconv.Atoi(*gameIdString)
	if err != nil {
		return -1, err
	}
	log.Printf("gameId: %v", gameId)
	return gameId, nil
}

func GetRecentGames(dynamoClient *dynamodb.DynamoDB) ([]model.Game, error) {
	scanInput := dynamodb.ScanInput{
		AttributesToGet: aws.StringSlice([]string{"Id", "Category", "Episode", "Quest", "Time", "Player", "Timestamp"}),
		Limit:           aws.Int64(20),
		TableName:       aws.String(gamesTable),
	}
	scan, err := dynamoClient.Scan(&scanInput)
	if err != nil {
		return nil, err
	}
	games := make([]model.Game, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(scan.Items, &games)
	return games, err
}

func GetGame(gameId string, dynamoClient *dynamodb.DynamoDB) (*model.QuestRun, error) {
	questRun := model.QuestRun{}
	game := model.Game{}
	primaryKey := dynamodb.AttributeValue{
		S: aws.String(gameId),
	}
	getItem := dynamodb.GetItemInput{
		TableName: aws.String(gamesTable),
		Key: map[string]*dynamodb.AttributeValue{"Id": &primaryKey},
	}
	item, err := dynamoClient.GetItem(&getItem)
	if err != nil || item.Item == nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalMap(item.Item, &game)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(game.GameGzip)
	reader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := io.ReadAll(reader)
	if err != io.ErrUnexpectedEOF {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &questRun)
	if err != nil {
		return nil, err
	}

	return &questRun, err
}