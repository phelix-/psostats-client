package server

import (
	"log"
	"testing"
)

func Test_hashPassword(t *testing.T) {
	passwordIn := "test"
	log.Printf("%v - %v", passwordIn, hashPassword(passwordIn))
}