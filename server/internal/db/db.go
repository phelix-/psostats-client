package db

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/phelix-/psostats/v2/pkg/model"
	"io"
	"log"
	"sort"
	"strconv"
	"time"
)

const (
	GamesByIdTable           = "games_by_id"
	QuestRecordsTable        = "quest_records"
	Anniv2021RecordsTable    = "records_anniv_2021"
	RecentGamesByPlayerTable = "recent_games_by_player_2"
	RecentGamesByMonth       = "recent_games_by_month_2"
	GameCountTable           = "games_counter"
	gameCountPrimaryKey      = "game_count"
	PlayerPbTable            = "player_pb"
	PlayerClassCount         = "player_class_count"
	PlayerQuestCount         = "player_quest_count"
	OverallQuestCount        = "overall_quest_count"
	QuestSeriesPbTable       = "quest_series_pb"
	QuestDataFramesTable     = "quest_data_frames"
)

type PsoStatsDb struct {
	dynamoClient *dynamodb.DynamoDB
}

func getCategoryFromQuest(questRun model.QuestRun) string {
	return getCategoryString(len(questRun.AllPlayers), questRun.PbCategory)
}

func getCategoryString(numPlayers int, pbCategory bool) string {
	category := fmt.Sprintf("%d", numPlayers)
	if pbCategory {
		category += "p"
	} else {
		category += "n"
	}
	return category
}

func WriteGameById(questRun *model.QuestRun, dynamoClient *dynamodb.DynamoDB) (string, error) {
	gameId, err := incrementAndGetGameId(dynamoClient)
	if err != nil {
		return "", err
	}
	questRun.Id = fmt.Sprintf("%d", gameId)
	_ = writeDataFrames(questRun, dynamoClient)
	questRun.DataFrames = make([]model.DataFrame, 0)
	gameGzip, err := CompressGame(questRun)
	if err != nil {
		return "", err
	}

	playerIndex := -1
	for i, player := range questRun.AllPlayers {
		if player.GuildCard == questRun.GuildCard {
			playerIndex = i + 1
		}
	}

	game := summaryFromQuestRun(*questRun)
	game.GameGzip = gameGzip
	switch playerIndex {
	case 1:
		game.P1Gzip = gameGzip
		game.P1HasStats = true
	case 2:
		game.P2Gzip = gameGzip
		game.P2HasStats = true
	case 3:
		game.P3Gzip = gameGzip
		game.P3HasStats = true
	case 4:
		game.P4Gzip = gameGzip
		game.P4HasStats = true
	}

	marshalled, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		return "", err
	}
	delete(marshalled, "FormattedDate")
	delete(marshalled, "FormattedTime")
	delete(marshalled, "QuestAndCategory")
	input2 := &dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String(GamesByIdTable),
	}
	_, err = dynamoClient.PutItem(input2)
	if err != nil {
		return "", err
	}
	err = incrementOverallQuestCount(game.Episode, game.Quest, dynamoClient)
	if err != nil {
		return "", err
	}
	err = writeRecentGame(game, dynamoClient)
	return game.Id, err
}

func getPlayerIndex(questRun model.QuestRun) (int, error) {
	playerIndex := -1
	for i, player := range questRun.AllPlayers {
		if player.GuildCard == questRun.GuildCard {
			playerIndex = i + 1
		}
	}
	if playerIndex <= 0 || playerIndex > 4 {
		return -1, errors.New(fmt.Sprintf("player index was out of bounds: %d", playerIndex))
	}
	return playerIndex, nil
}

type QuestDataFrameItem struct {
	QuestAndPlayerId     string
	CompressedDataFrames []byte
}

func writeDataFrames(questRun *model.QuestRun, db *dynamodb.DynamoDB) error {
	buffer := new(bytes.Buffer)
	writer := gzip.NewWriter(buffer)
	jsonQuestBytes, err := json.Marshal(questRun.DataFrames)
	if err != nil {
		return err
	}
	_, err = writer.Write(jsonQuestBytes)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	index, _ := getPlayerIndex(*questRun)
	item := QuestDataFrameItem{
		QuestAndPlayerId:     fmt.Sprintf("%s_%d", questRun.Id, index),
		CompressedDataFrames: buffer.Bytes(),
	}
	marshalled, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	putItem := &dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String(QuestDataFramesTable),
	}
	_, err = db.PutItem(putItem)
	return err
}

func GetDataFrames(gameId string, gem int, db *dynamodb.DynamoDB) ([]model.DataFrame, error) {
	questDataFrameItem := QuestDataFrameItem{}
	questAndCategory := fmt.Sprintf("%s_%d", gameId, gem)
	primaryKey := dynamodb.AttributeValue{
		S: aws.String(questAndCategory),
	}
	getItem := dynamodb.GetItemInput{
		TableName: aws.String(QuestDataFramesTable),
		Key:       map[string]*dynamodb.AttributeValue{"QuestAndPlayerId": &primaryKey},
	}
	item, err := db.GetItem(&getItem)
	if err != nil || item.Item == nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalMap(item.Item, &questDataFrameItem)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(questDataFrameItem.CompressedDataFrames)
	reader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := io.ReadAll(reader)
	if err != io.ErrUnexpectedEOF {
		return nil, err
	}
	dataFrames := make([]model.DataFrame, 0)
	err = json.Unmarshal(jsonBytes, &dataFrames)
	return dataFrames, err
}

func AttachGameToId(questRun model.QuestRun, id string, dynamoClient *dynamodb.DynamoDB) error {
	_ = writeDataFrames(&questRun, dynamoClient)
	questRun.DataFrames = make([]model.DataFrame, 0)
	gameGzip, err := CompressGame(&questRun)
	if err != nil {
		return err
	}
	playerIndex, err := getPlayerIndex(questRun)
	if err != nil {
		return err
	}
	key := make(map[string]*dynamodb.AttributeValue)
	idAttribute := dynamodb.AttributeValue{S: aws.String(id)}
	key["Id"] = &idAttribute
	gzipAttribute := dynamodb.AttributeValue{B: gameGzip}
	values := make(map[string]*dynamodb.AttributeValue)
	values[":g"] = &gzipAttribute
	trueAttribute := dynamodb.AttributeValue{BOOL: aws.Bool(true)}
	values[":h"] = &trueAttribute

	putItemInput := dynamodb.UpdateItemInput{
		Key:                       key,
		UpdateExpression:          aws.String(fmt.Sprintf("SET P%dGzip = :g, P%dHasStats = :h", playerIndex, playerIndex)),
		ExpressionAttributeValues: values,
		TableName:                 aws.String(GamesByIdTable),
	}
	_, err = dynamoClient.UpdateItem(&putItemInput)
	if err != nil {
		return err
	}

	month := questRun.QuestStartTime.UTC().Format("01/2006")
	monthAttribute := dynamodb.AttributeValue{S: aws.String(month)}
	byMonthKey := map[string]*dynamodb.AttributeValue{
		"IdInt": {N: aws.String(id)},
		"Month": &monthAttribute,
	}
	values = map[string]*dynamodb.AttributeValue{
		":h": &trueAttribute,
	}

	putItemInput = dynamodb.UpdateItemInput{
		Key:                       byMonthKey,
		UpdateExpression:          aws.String(fmt.Sprintf("SET P%dHasStats = :h", playerIndex)),
		ExpressionAttributeValues: values,
		TableName:                 aws.String(RecentGamesByMonth),
	}
	_, err = dynamoClient.UpdateItem(&putItemInput)
	return err
}

func AddPovToRecord(tableName string, questRun model.QuestRun, dynamoClient *dynamodb.DynamoDB) error {
	gameSummary := summaryFromQuestRun(questRun)
	playerIndex, err := getPlayerIndex(questRun)
	if err != nil {
		return err
	}
	questName := dynamodb.AttributeValue{S: aws.String(gameSummary.Quest)}
	category := dynamodb.AttributeValue{S: aws.String(gameSummary.Category)}
	key := map[string]*dynamodb.AttributeValue{
		"Quest":    &questName,
		"Category": &category,
	}
	trueAttribute := dynamodb.AttributeValue{BOOL: aws.Bool(true)}
	values := map[string]*dynamodb.AttributeValue{
		":h": &trueAttribute,
	}

	updateItemInput := dynamodb.UpdateItemInput{
		Key:                       key,
		UpdateExpression:          aws.String(fmt.Sprintf("SET P%dHasStats = :h", playerIndex)),
		ExpressionAttributeValues: values,
		TableName:                 aws.String(tableName),
	}
	_, err = dynamoClient.UpdateItem(&updateItemInput)
	return err
}

func writeRecentGame(game model.Game, dynamoClient *dynamodb.DynamoDB) error {
	month := game.Timestamp.UTC().Format("01/2006")
	marshalled, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		return err
	}
	delete(marshalled, "FormattedDate")
	delete(marshalled, "FormattedTime")
	delete(marshalled, "GameGzip")
	delete(marshalled, "P1Gzip")
	delete(marshalled, "P2Gzip")
	delete(marshalled, "P3Gzip")
	delete(marshalled, "P4Gzip")
	delete(marshalled, "QuestAndCategory")
	monthAttribute := dynamodb.AttributeValue{
		S: aws.String(month),
	}
	marshalled["Month"] = &monthAttribute
	input := &dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String(RecentGamesByMonth),
	}
	_, err = dynamoClient.PutItem(input)
	return err
}

func CompressGame(questRun *model.QuestRun) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer := gzip.NewWriter(buffer)
	jsonQuestBytes, err := json.Marshal(questRun)
	if err != nil {
		return nil, err
	}
	_, err = writer.Write(jsonQuestBytes)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func writeRecord(tableName string, questRun *model.QuestRun, dynamoClient *dynamodb.DynamoDB) error {
	gameSummary := summaryFromQuestRun(*questRun)
	marshalledSummary, err := dynamodbattribute.MarshalMap(gameSummary)
	delete(marshalledSummary, "GameGzip")
	delete(marshalledSummary, "P1Gzip")
	delete(marshalledSummary, "P2Gzip")
	delete(marshalledSummary, "P3Gzip")
	delete(marshalledSummary, "P4Gzip")
	delete(marshalledSummary, "FormattedDate")
	delete(marshalledSummary, "FormattedTime")
	if err != nil {
		return err
	}
	gamesByQuestInput := &dynamodb.PutItemInput{
		Item:      marshalledSummary,
		TableName: aws.String(tableName),
	}
	_, err = dynamoClient.PutItem(gamesByQuestInput)
	return err
}

func WriteGameByQuestRecord(questRun *model.QuestRun, dynamoClient *dynamodb.DynamoDB) error {
	return writeRecord(QuestRecordsTable, questRun, dynamoClient)
}

func WriteAnniv2021Record(questRun *model.QuestRun, dynamoClient *dynamodb.DynamoDB) error {
	return writeRecord(Anniv2021RecordsTable, questRun, dynamoClient)
}

func getRecord(
	tableName string,
	quest string,
	numPlayers int,
	pbCategory bool,
	dynamoClient *dynamodb.DynamoDB,
) (*model.Game, error) {
	category := getCategoryString(numPlayers, pbCategory)
	requestExpression, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("Quest"), expression.Value(quest)).
			And(expression.KeyEqual(expression.Key("Category"), expression.Value(category)))).
		Build()
	if err != nil {
		return nil, err
	}
	result, err := dynamoClient.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  requestExpression.Names(),
		ExpressionAttributeValues: requestExpression.Values(),
		KeyConditionExpression:    requestExpression.KeyCondition(),
		TableName:                 aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}
	games := make([]model.Game, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &games)
	if len(games) > 0 {
		return &games[0], nil
	} else {
		return nil, nil
	}
}

func GetQuestRecord(quest string, numPlayers int, pbCategory bool, dynamoClient *dynamodb.DynamoDB) (*model.Game, error) {
	return getRecord(QuestRecordsTable, quest, numPlayers, pbCategory, dynamoClient)
}

func GetAnniv2021Record(quest string, numPlayers int, pbCategory bool, dynamoClient *dynamodb.DynamoDB) (*model.Game, error) {
	return getRecord(Anniv2021RecordsTable, quest, numPlayers, pbCategory, dynamoClient)
}

func GetQuestRecords(tableName string, dynamoClient *dynamodb.DynamoDB) ([]model.Game, error) {
	scanInput := dynamodb.ScanInput{
		AttributesToGet: aws.StringSlice([]string{"Id", "Category", "Episode", "Quest",
			"Time", "Player", "Timestamp", "PlayerNames", "PlayerClasses", "PlayerGcs",
			"P1HasStats", "P2HasStats", "P3HasStats", "P4HasStats", "Points"}),
		Limit:     aws.Int64(1000),
		TableName: aws.String(tableName),
	}
	scan, err := dynamoClient.Scan(&scanInput)
	if err != nil {
		return nil, err
	}
	games := make([]model.Game, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(scan.Items, &games)
	return games, err
}

func GetPlayerPbs(player string, dynamoClient *dynamodb.DynamoDB) ([]model.Game, error) {
	requestExpression, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("Player"), expression.Value(player))).
		Build()
	if err != nil {
		return nil, err
	}
	result, err := dynamoClient.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  requestExpression.Names(),
		ExpressionAttributeValues: requestExpression.Values(),
		KeyConditionExpression:    requestExpression.KeyCondition(),
		TableName:                 aws.String(PlayerPbTable),
	})
	if err != nil {
		return nil, err
	}

	games := make([]model.Game, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &games)
	return games, err
}

func GetPlayerRecentGames(player string, dynamoClient *dynamodb.DynamoDB, limit int64) ([]model.Game, error) {
	requestExpression, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("Player"), expression.Value(player))).
		Build()
	if err != nil {
		return nil, err
	}
	result, err := dynamoClient.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  requestExpression.Names(),
		ExpressionAttributeValues: requestExpression.Values(),
		KeyConditionExpression:    requestExpression.KeyCondition(),
		ScanIndexForward:          aws.Bool(false),
		Limit:                     aws.Int64(limit),
		TableName:                 aws.String(RecentGamesByPlayerTable),
	})
	if err != nil {
		return nil, err
	}
	games := make([]model.Game, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &games)
	if len(result.Items) < int(limit) {
		oldGames, err := GetPlayerRecentGamesOld(player, dynamoClient)
		if err != nil {
			return nil, err
		}
		games = append(games, oldGames...)
	}
	sort.Slice(games, func(i, j int) bool { return games[i].Timestamp.After(games[j].Timestamp) })
	if len(games) > int(limit) {
		games = games[0:limit]
	}
	return games, err
}

func GetPlayerRecentGamesOld(player string, dynamoClient *dynamodb.DynamoDB) ([]model.Game, error) {
	requestExpression, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("Player"), expression.Value(player))).
		Build()
	if err != nil {
		return nil, err
	}
	result, err := dynamoClient.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  requestExpression.Names(),
		ExpressionAttributeValues: requestExpression.Values(),
		KeyConditionExpression:    requestExpression.KeyCondition(),
		ScanIndexForward:          aws.Bool(false),
		TableName:                 aws.String("recent_games_by_player"),
	})
	if err != nil {
		return nil, err
	}
	games := make([]model.Game, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &games)
	sort.Slice(games, func(i, j int) bool { return games[i].Timestamp.After(games[j].Timestamp) })
	return games, err
}

type QuestSeriesPb struct {
	Series       string
	UserAndQuest string
	User         string
	QuestName    string
	Id           int
	Time         time.Duration
	Timestamp    time.Time
}

func questSeriesPbFromQuestRun(series string, questRun *model.QuestRun) *QuestSeriesPb {
	idInt, _ := strconv.Atoi(questRun.Id)
	duration, _ := time.ParseDuration(questRun.QuestDuration)
	pb := QuestSeriesPb{
		Series:       series,
		UserAndQuest: fmt.Sprintf("%s+%s", questRun.UserName, questRun.QuestName),
		User:         questRun.UserName,
		QuestName:    questRun.QuestName,
		Id:           idInt,
		Time:         duration,
		Timestamp:    questRun.QuestStartTime,
	}
	return &pb
}

func WriteQuestSeriesPb(series string, questRun *model.QuestRun, db *dynamodb.DynamoDB) error {
	pb := questSeriesPbFromQuestRun(series, questRun)
	marshalled, err := dynamodbattribute.MarshalMap(pb)
	if err != nil {
		return err
	}
	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String(QuestSeriesPbTable),
	})
	return err
}

func GetQuestSeriesPb(series, user, quest string, db *dynamodb.DynamoDB) (*QuestSeriesPb, error) {
	userAndQuest := fmt.Sprintf("%s+%s", user, quest)
	requestExpression, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("UserAndQuest"), expression.Value(userAndQuest)).
			And(expression.KeyEqual(expression.Key("Series"), expression.Value(series)))).
		Build()
	if err != nil {
		return nil, err
	}

	result, err := db.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  requestExpression.Names(),
		ExpressionAttributeValues: requestExpression.Values(),
		KeyConditionExpression:    requestExpression.KeyCondition(),
		TableName:                 aws.String(QuestSeriesPbTable),
	})
	if err != nil {
		return nil, err
	}

	pb := QuestSeriesPb{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &pb)
	return &pb, err
}

func WriteAnniversaryStats(questRun *model.QuestRun, db *dynamodb.DynamoDB) {
	mesetaCharged := int64(questRun.MesetaCharged[len(questRun.MesetaCharged)-1])
	addAnniversaryCounter(questRun.QuestName, "MesetaCharged", mesetaCharged)
	addAnniversaryCounter(questRun.QuestName, fmt.Sprintf("MesetaCharged.%s", questRun.PlayerClass), mesetaCharged)
	addAnniversaryCounter(questRun.QuestName, "Runs", 1)
	addAnniversaryCounter(questRun.QuestName, "Moving", int64(questRun.TimeByState[2]+questRun.TimeByState[4]))
	addAnniversaryCounter(questRun.QuestName, "Standing", int64(questRun.TimeByState[1]))
	addAnniversaryCounter(questRun.QuestName, "Attacking", int64(questRun.TimeByState[5]+questRun.TimeByState[6]+questRun.TimeByState[7]))
	addAnniversaryCounter(questRun.QuestName, "Casting", int64(questRun.TimeByState[8]))
	addAnniversaryCounter(questRun.QuestName, "Deaths", int64(questRun.DeathCount))
	addAnniversaryCounter(questRun.QuestName, fmt.Sprintf("ClassUse.%s", questRun.PlayerClass), 1)
}

func addAnniversaryCounter(questName, counterName string, value int64) {

}

func GetPlayerPB(quest, player string, numPlayers int, pbCategory bool, dynamoClient *dynamodb.DynamoDB) (*model.Game, error) {
	category := getCategoryString(numPlayers, pbCategory)
	questAndCategory := fmt.Sprintf("%v+%v", quest, category)
	requestExpression, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("Player"), expression.Value(player)).
			And(expression.KeyEqual(expression.Key("QuestAndCategory"), expression.Value(questAndCategory)))).
		Build()
	if err != nil {
		return nil, err
	}
	result, err := dynamoClient.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  requestExpression.Names(),
		ExpressionAttributeValues: requestExpression.Values(),
		KeyConditionExpression:    requestExpression.KeyCondition(),
		TableName:                 aws.String(PlayerPbTable),
	})
	if err != nil {
		return nil, err
	}
	games := make([]model.Game, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &games)
	if len(games) > 0 {
		return &games[0], nil
	} else {
		return nil, nil
	}
}

func WritePlayerPb(questRun *model.QuestRun, dynamoClient *dynamodb.DynamoDB) error {
	gameSummary := summaryFromQuestRun(*questRun)
	marshalledSummary, err := dynamodbattribute.MarshalMap(gameSummary)
	if err != nil {
		return err
	}
	delete(marshalledSummary, "GameGzip")
	delete(marshalledSummary, "P1Gzip")
	delete(marshalledSummary, "P2Gzip")
	delete(marshalledSummary, "P3Gzip")
	delete(marshalledSummary, "P4Gzip")
	delete(marshalledSummary, "FormattedDate")
	delete(marshalledSummary, "FormattedTime")
	gamesByQuestInput := &dynamodb.PutItemInput{
		Item:      marshalledSummary,
		TableName: aws.String(PlayerPbTable),
	}
	_, err = dynamoClient.PutItem(gamesByQuestInput)
	return err
}

func WriteGameByPlayer(questRun *model.QuestRun, dynamoClient *dynamodb.DynamoDB) error {
	gameSummary := summaryFromQuestRun(*questRun)
	marshalledSummary, err := dynamodbattribute.MarshalMap(gameSummary)
	delete(marshalledSummary, "GameGzip")
	delete(marshalledSummary, "P1Gzip")
	delete(marshalledSummary, "P2Gzip")
	delete(marshalledSummary, "P3Gzip")
	delete(marshalledSummary, "P4Gzip")
	delete(marshalledSummary, "FormattedDate")
	delete(marshalledSummary, "FormattedTime")
	if err != nil {
		return err
	}
	gamesByQuestInput := &dynamodb.PutItemInput{
		Item:      marshalledSummary,
		TableName: aws.String(RecentGamesByPlayerTable),
	}
	_, err = dynamoClient.PutItem(gamesByQuestInput)
	if err != nil {
		return err
	}
	err = incrementPlayerQuestCount(questRun.UserName, int(questRun.Episode), questRun.QuestName, dynamoClient)
	if err != nil {
		return err
	}
	err = incrementPlayerClassCount(questRun.UserName, questRun.PlayerClass, dynamoClient)
	return err
}

func GetPlayerClassCounts(playerName string, dynamoClient *dynamodb.DynamoDB) (map[string]int, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("Player"), expression.Value(playerName))).
		Build()
	if err != nil {
		return nil, err
	}
	result, err := dynamoClient.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(PlayerClassCount),
	})
	classUsage := map[string]int{}
	for _, item := range result.Items {
		class := item["Class"].S
		count := item["count"].N
		atoi, err := strconv.Atoi(*count)
		if err != nil {
			return nil, err
		}
		classUsage[*class] = atoi
	}
	return classUsage, nil
}

func GetPlayerQuestCounts(playerName string, dynamoClient *dynamodb.DynamoDB) (map[string]int, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("Player"), expression.Value(playerName))).
		Build()
	if err != nil {
		return nil, err
	}
	result, err := dynamoClient.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(PlayerQuestCount),
	})
	questsPlayed := map[string]int{}
	for _, item := range result.Items {
		class := item["Quest"].S
		count := item["count"].N
		atoi, err := strconv.Atoi(*count)
		if err != nil {
			return nil, err
		}
		questsPlayed[*class] = atoi
	}
	return questsPlayed, nil
}

func incrementPlayerQuestCount(playerName string, episode int, questName string, dynamoClient *dynamodb.DynamoDB) error {
	oneValue := dynamodb.AttributeValue{N: aws.String("1")}
	update := dynamodb.AttributeValueUpdate{
		Action: aws.String(dynamodb.AttributeActionAdd),
		Value:  &oneValue,
	}
	playerNameAttr := dynamodb.AttributeValue{S: aws.String(playerName)}
	questNameAttr := dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("%d_%v", episode, questName))}

	updateItemInput := dynamodb.UpdateItemInput{
		TableName:        aws.String(PlayerQuestCount),
		Key:              map[string]*dynamodb.AttributeValue{"Player": &playerNameAttr, "Quest": &questNameAttr},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{"count": &update},
	}
	_, err := dynamoClient.UpdateItem(&updateItemInput)
	return err
}

func incrementPlayerClassCount(playerName string, className string, dynamoClient *dynamodb.DynamoDB) error {
	oneValue := dynamodb.AttributeValue{N: aws.String("1")}
	update := dynamodb.AttributeValueUpdate{
		Action: aws.String(dynamodb.AttributeActionAdd),
		Value:  &oneValue,
	}
	playerNameAttr := dynamodb.AttributeValue{S: aws.String(playerName)}
	questNameAttr := dynamodb.AttributeValue{S: aws.String(className)}

	updateItemInput := dynamodb.UpdateItemInput{
		TableName:        aws.String(PlayerClassCount),
		Key:              map[string]*dynamodb.AttributeValue{"Player": &playerNameAttr, "Class": &questNameAttr},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{"count": &update},
	}
	_, err := dynamoClient.UpdateItem(&updateItemInput)
	return err
}

func incrementOverallQuestCount(episode int, questName string, dynamoClient *dynamodb.DynamoDB) error {
	oneValue := dynamodb.AttributeValue{N: aws.String("1")}
	update := dynamodb.AttributeValueUpdate{
		Action: aws.String(dynamodb.AttributeActionAdd),
		Value:  &oneValue,
	}
	questNameAttr := dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("%d_%v", episode, questName))}

	updateItemInput := dynamodb.UpdateItemInput{
		TableName:        aws.String(OverallQuestCount),
		Key:              map[string]*dynamodb.AttributeValue{"Quest": &questNameAttr},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{"count": &update},
	}
	_, err := dynamoClient.UpdateItem(&updateItemInput)
	return err
}

func summaryFromQuestRun(questRun model.QuestRun) model.Game {
	category := getCategoryFromQuest(questRun)
	duration, err := time.ParseDuration(questRun.QuestDuration)
	if err != nil {
		log.Printf("Failed parsing duration gameId %v", questRun.Id)
	}
	questAndCategory := fmt.Sprintf("%v+%v", questRun.QuestName, category)
	playerNames := make([]string, 4)
	playerClasses := make([]string, 4)
	playerGcs := make([]string, 4)
	for i, basePlayerInfo := range questRun.AllPlayers {
		playerNames[i] = basePlayerInfo.Name
		playerClasses[i] = basePlayerInfo.Class
		playerGcs[i] = basePlayerInfo.GuildCard
	}
	playerIndex, _ := getPlayerIndex(questRun)
	idInt, _ := strconv.Atoi(questRun.Id)
	return model.Game{
		Id:               questRun.Id,
		IdInt:            idInt,
		Player:           questRun.UserName,
		PlayerNames:      playerNames,
		PlayerClasses:    playerClasses,
		PlayerGcs:        playerGcs,
		P1HasStats:       playerIndex == 1,
		P2HasStats:       playerIndex == 2,
		P3HasStats:       playerIndex == 3,
		P4HasStats:       playerIndex == 4,
		Category:         category,
		Quest:            questRun.QuestName,
		QuestAndCategory: questAndCategory,
		Time:             duration,
		Timestamp:        questRun.QuestStartTime,
		Episode:          int(questRun.Episode),
		Points:           int(questRun.Points),
	}
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
		TableName:        aws.String(GameCountTable),
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
	thisMonth := time.Now().UTC().Format("01/2006")
	lastMonth := time.Now().UTC().AddDate(0, -1, 0).Format("01/2006")
	games, err := GetGamesForMonth(thisMonth, 30, dynamoClient)
	if err != nil {
		return nil, err
	}
	if len(games) < 30 {
		lastMonthGames, err := GetGamesForMonth(lastMonth, 30, dynamoClient)
		if err != nil {
			return nil, err
		}
		games = append(games, lastMonthGames...)
	}

	sort.Slice(games, func(i, j int) bool { return games[i].Timestamp.After(games[j].Timestamp) })
	if len(games) > 30 {
		games = games[0:30]
	}
	return games, err
}

func GetGamesForMonth(month string, limit int, dynamoClient *dynamodb.DynamoDB) ([]model.Game, error) {
	requestExpression, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("Month"), expression.Value(month))).
		Build()
	if err != nil {
		return nil, err
	}
	result, err := dynamoClient.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  requestExpression.Names(),
		ExpressionAttributeValues: requestExpression.Values(),
		KeyConditionExpression:    requestExpression.KeyCondition(),
		ScanIndexForward:          aws.Bool(false),
		Limit:                     aws.Int64(int64(limit)),
		TableName:                 aws.String(RecentGamesByMonth),
	})

	if err != nil {
		return nil, err
	}
	games := make([]model.Game, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &games)
	return games, err
}

func GetFullGame(gameId string, dynamoClient *dynamodb.DynamoDB) (*model.Game, error) {
	game := model.Game{}
	primaryKey := dynamodb.AttributeValue{
		S: aws.String(gameId),
	}
	getItem := dynamodb.GetItemInput{
		TableName: aws.String(GamesByIdTable),
		Key:       map[string]*dynamodb.AttributeValue{"Id": &primaryKey},
	}
	item, err := dynamoClient.GetItem(&getItem)
	if err != nil || item.Item == nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalMap(item.Item, &game)
	return &game, err
}

func GetGame(gameId string, gem int, dynamoClient *dynamodb.DynamoDB) (*model.QuestRun, error) {
	questRun := model.QuestRun{}
	game := model.Game{}
	primaryKey := dynamodb.AttributeValue{
		S: aws.String(gameId),
	}
	getItem := dynamodb.GetItemInput{
		TableName: aws.String(GamesByIdTable),
		Key:       map[string]*dynamodb.AttributeValue{"Id": &primaryKey},
	}
	item, err := dynamoClient.GetItem(&getItem)
	if err != nil || item.Item == nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalMap(item.Item, &game)
	if err != nil {
		return nil, err
	}

	gameGzip := getGameGzip(&game, gem)

	buffer := bytes.NewBuffer(gameGzip)
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

	if len(questRun.DataFrames) == 0 {
		playerIndex, _ := getPlayerIndex(questRun)
		dataFrames, err := GetDataFrames(questRun.Id, playerIndex, dynamoClient)
		if err != nil {
			return nil, err
		}
		questRun.DataFrames = dataFrames
	}

	return &questRun, err
}

func getGameGzip(game *model.Game, gem int) []byte {
	var gameGzip []byte
	if game != nil {
		switch gem + 1 {
		case 1:
			gameGzip = game.P1Gzip
		case 2:
			gameGzip = game.P2Gzip
		case 3:
			gameGzip = game.P3Gzip
		case 4:
			gameGzip = game.P4Gzip
		}
	}
	if gameGzip == nil {
		gameGzip = game.GameGzip
	}

	return gameGzip
}
