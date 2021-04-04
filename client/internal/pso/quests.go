package pso

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Trigger struct {
	Register *int   `yaml:"register"`
	Floor    uint16 `yaml:"floor"`
	Switch   int    `yaml:"switch"`
	WarpIn   bool   `yaml:"warpIn"`
}

type Quest struct {
	Ignore       bool    `yaml:"ignore"`
	StartTrigger Trigger `yaml:"start"`
	EndTrigger   Trigger `yaml:"end"`
}

type Quests struct {
	allQuests    map[string]map[string]Quest
	warnedQuests map[string]bool
}

func NewQuests() Quests {
	allQuests := make(map[string]map[string]Quest)
	data, err := ioutil.ReadFile("./quests.yaml")
	if err != nil {
		log.Fatalf("Error reading quests file %v", err)
	}

	err = yaml.Unmarshal(data, allQuests)
	if err != nil {
		log.Fatalf("Error parsing quests file %v", err)
	}

	return Quests{
		allQuests:    allQuests,
		warnedQuests: make(map[string]bool),
	}
}

func (q *Quests) GetQuest(episode int, questName string) (Quest, bool) {
	questsForEpisode, exists := q.allQuests[fmt.Sprintf("Episode %v", episode)]
	if !exists {
		log.Printf("Episode %v not found?", episode)
		return Quest{}, false
	}
	quest, questFound := questsForEpisode[questName]
	if !questFound {
		if warned := q.warnedQuests[questName]; !warned {
			q.warnedQuests[questName] = true
			log.Printf("Episode %v Quest '%v' not configured", episode, questName)
		}
	}
	return quest, questFound
}

func (q *Quest) StartsOnRegister() bool {
	return q.StartTrigger.Register != nil
}

func (q *Quest) StartsAtWarpIn() bool {
	return q.StartTrigger.WarpIn
}

func (q *Quest) TerminalQuest() bool {
	return !q.StartsAtWarpIn() && !q.StartsOnRegister()
}

func (q *Quest) EndsOnRegister() bool {
	return q.EndTrigger.Register != nil
}
