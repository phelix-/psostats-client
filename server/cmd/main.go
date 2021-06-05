package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/phelix-/psostats/v2/server/internal/server"
	"log"
	"os"
)

func main() {
	version := "0.6.2"

	log.Printf("Starting Up PSOStats Server %v", version)

	var awsSession *session.Session
	if _, set := os.LookupEnv("AWS_ACCESS_KEY_ID"); set {
		awsSession = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	} else {
		awsSession = session.Must(session.NewSession(&aws.Config{
			Region:   aws.String("us-west-2"),
			Endpoint: aws.String("http://localhost:8000"),
		}))
	}

	dynamoClient := dynamodb.New(awsSession)
	s := server.New(dynamoClient)
	s.Run()
}
