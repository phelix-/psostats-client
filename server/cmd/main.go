package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/phelix-/psostats/v2/server/internal/server"
	"log"
)

func main() {
	version := "0.0.1"

	log.Printf("Starting Up PSOStats Server %v", version)

	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamoClient := dynamodb.New(awsSession)

	// err := listTables(dynamoClient)
	// if err != nil {
	//	log.Fatal(err)
	// }
	// writeGame(dynamoClient)
	s := server.New(dynamoClient)
	s.Run()
}

//
// func listTables(dynamoClient *dynamodb.DynamoDB) error {
//	listTablesInput := &dynamodb.ListTablesInput{}
//	result, err := dynamoClient.ListTables(listTablesInput)
//	if err != nil {
//		return err
//	}
//	for _, tableName := range result.TableNames {
//		log.Print(*tableName)
//	}
//	return nil
// }


