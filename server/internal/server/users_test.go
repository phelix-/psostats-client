package server_test

import (
	"github.com/phelix-/psostats/v2/server/internal/server"
	"log"
	"testing"
)

func Test_hashPassword(t *testing.T) {
	passwordIn := "test"
	log.Printf("%v - %v", passwordIn, server.HashPassword(passwordIn))
}