package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/phelix-/psostats/v2/server/internal/server"
	"log"
)

func main() {
	version := "0.1.0"

	log.Printf("Starting Up PSOStats Server %v", version)

	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamoClient := dynamodb.New(awsSession)
	s := server.New(dynamoClient)
	s.Run()
}
