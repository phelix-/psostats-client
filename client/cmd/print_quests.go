package main

import (
	"fmt"
	"github.com/phelix-/psostats/v2/client/internal/pso/quest"
)

func main() {
	for _, aquest := range quest.GetAllQuests() {
		if aquest.Remap == nil {
			fmt.Println(aquest.Name)
		}
	}

}
