package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/phelix-/psostats/v2/pkg/server"
	"log"
	"time"
)

func main() {
	version := "0.0.1"

	log.Printf("Starting Up PSOStats Server %v", version)

	//awsSession := session.Must(session.NewSessionWithOptions(session.Options{
	//	SharedConfigState: session.SharedConfigEnable,
	//}))
	//
	//dynamoClient := dynamodb.New(awsSession)

	//err := listTables(dynamoClient)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//writeGame(dynamoClient)
	s := server.New(nil /*dynamoClient*/)
	s.Run()
}

//
//func listTables(dynamoClient *dynamodb.DynamoDB) error {
//	listTablesInput := &dynamodb.ListTablesInput{}
//	result, err := dynamoClient.ListTables(listTablesInput)
//	if err != nil {
//		return err
//	}
//	for _, tableName := range result.TableNames {
//		log.Print(*tableName)
//	}
//	return nil
//}

func writeGame(dynamoClient *dynamodb.DynamoDB) error {
	tableName := "games_by_id"
	game := Game{
		Id:        "0",
		Player:    "unseen+42",
		Category:  "2n",
		Quest:     "Sweep-up Operation #4",
		Time:      time.Minute*7 + time.Second*33 + time.Millisecond*300,
		Timestamp: time.Date(2021, 4, 1, 23, 15, 0, 0, time.Local),
	}
	marshalled, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String(tableName),
	}
	_, err = dynamoClient.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

type Game struct {
	Id        string
	Player    string
	Category  string
	Episode   int
	Quest     string
	Time      time.Duration
	Timestamp time.Time
}
