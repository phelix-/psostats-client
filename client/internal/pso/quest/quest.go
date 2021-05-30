package quest

import (
	"errors"
	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/client/internal/numbers"
	"log"
	"strconv"
	"strings"
)

func GetQuestPointer(handle w32.HANDLE) uintptr {
	return uintptr(numbers.ReadU32Unchecked(handle, uintptr(0x00A95AA8)))
}

func GetQuestDataPointer(handle w32.HANDLE, questPtr uintptr) uintptr {
	return uintptr(numbers.ReadU32Unchecked(handle, questPtr+0x19C))
}

func GetQuestRegisterPointer(handle w32.HANDLE) uintptr {
	return uintptr(numbers.ReadU32Unchecked(handle, uintptr(0x00A954B0)))
}

func IsRegisterSet(handle w32.HANDLE, registerId uint16, questRegisterAddress uintptr) (bool, error) {
	registerSet := false
	if questRegisterAddress != 0 {
		buf, _, ok := w32.ReadProcessMemory(handle, questRegisterAddress+(4*uintptr(registerId)), 2)
		if !ok {
			return false, errors.New("unable to isQuestRegisterSet")
		}
		registerSet = buf[0] > 0
	}
	return registerSet, nil
}

func warpIn() Trigger {
	return Trigger{ WarpIn: true }
}

func register(register int) Trigger {
	registerU16 := uint16(register)
	return Trigger{ Register: &registerU16 }
}

func floorSwitch(floor int, switchId int) Trigger {
	return Trigger{Floor: uint16(floor), Switch: uint16(switchId)}
}

func remap(questName string) *string {
	return &questName
}

type Trigger struct {
	Register *uint16 `yaml:"register"`
	Floor    uint16  `yaml:"floor"`
	Switch   uint16  `yaml:"switch"`
	WarpIn   bool    `yaml:"warpIn"`
}

type Quest struct {
	Episode int
	Name    string
	Ignore  bool    `yaml:"ignore"`
	Remap   *string `yaml:"remap"`
	Start   Trigger `yaml:"start"`
	End     Trigger `yaml:"end"`
}

type Quests struct {
	allQuests    map[int]map[string]Quest
	warnedQuests map[string]bool
}

func NewQuests() Quests {
	allQuests := make(map[int]map[string]Quest)
	unsortedQuests := getAllQuests()

	for _, quest := range unsortedQuests {
		questsForEpisode := allQuests[quest.Episode]
		if questsForEpisode == nil {
			questsForEpisode = make(map[string]Quest)
		}
		questsForEpisode[quest.Name] = quest
		allQuests[quest.Episode] = questsForEpisode
	}

	return Quests{
		allQuests:    allQuests,
		warnedQuests: make(map[string]bool),
	}
}

func (q *Quests) GetQuestConfig(episode int, questName string) (Quest, bool) {
	questsForEpisode, exists := q.allQuests[episode]
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
	if questFound && quest.Remap != nil {
		return q.GetQuestConfig(episode, *quest.Remap)
	}
	return quest, questFound
}

func (q *Quest) StartsOnRegister() bool {
	return q.Start.Register != nil
}

func (q *Quest) StartsAtWarpIn() bool {
	return q.Start.WarpIn
}

func (q *Quest) TerminalQuest() bool {
	return !q.StartsAtWarpIn() && !q.StartsOnRegister()
}

func (q *Quest) EndsOnRegister() bool {
	return q.End.Register != nil
}

func (q *Quest) GetCmodeStage() int {
	if strings.HasPrefix(q.Name, "Stage") {
		stageNumber := strings.TrimPrefix(q.Name, "Stage")
		num, _ := strconv.Atoi(stageNumber)
		return num
	} else {
		return -1
	}
}