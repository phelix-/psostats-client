package quest

import (
	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/client/internal/numbers"
	"log"
)

func GetQuestPointer(handle w32.HANDLE) uintptr {
	return uintptr(numbers.ReadU32Unchecked(handle, uintptr(0x00A95AA8)))
}

func GetQuestDataPointer(handle w32.HANDLE, questPtr uintptr) uintptr {
	return uintptr(numbers.ReadU32Unchecked(handle, questPtr+0x19C))
}

func GetQuestRegisterPointer(handle w32.HANDLE, questPtr uintptr) uintptr {
	return uintptr(numbers.ReadU32Unchecked(handle, questPtr+0x2C))
}

func IsRegisterSet(handle w32.HANDLE, registerId uint16, questRegisterAddress uintptr) bool {
	value := GetRegisterValue(handle, registerId, questRegisterAddress)
	return value > 0
}

func GetRegisterValue(handle w32.HANDLE, registerId uint16, questRegisterAddress uintptr) uint16 {
	value := uint16(0)
	if questRegisterAddress != 0 {
		return numbers.ReadU16(handle, questRegisterAddress+(4*uintptr(registerId)))
	}
	return value
}

func warpIn() Trigger {
	return Trigger{WarpIn: true}
}

func register(register int) Trigger {
	registerU16 := uint16(register)
	return Trigger{Register: &registerU16}
}

func floorSwitch(floor int, switchId int) Trigger {
	return Trigger{Floor: uint16(floor), Switch: uint16(switchId)}
}

func remap(questName string) *string {
	return &questName
}

type Trigger struct {
	Register *uint16
	Floor    uint16
	Switch   uint16
	WarpIn   bool
}

type Quest struct {
	Episode       int
	Name          string
	Number        uint16
	Ignore        bool
	ForceTerminal bool
	Remap         *string
	CmodeStage    int
	Start         Trigger
	End           Trigger
	Splits        []Split
}

type Split struct {
	Name    string
	Trigger Trigger
}

type Quests struct {
	questsById   map[uint16]Quest
	allQuests    map[int]map[string]Quest
	warnedQuests map[string]bool
}

func NewQuests() Quests {
	questsById := make(map[uint16]Quest)
	allQuests := make(map[int]map[string]Quest)
	unsortedQuests := getAllQuests()

	for _, quest := range unsortedQuests {
		questsForEpisode := allQuests[quest.Episode]
		if questsForEpisode == nil {
			questsForEpisode = make(map[string]Quest)
		}
		questsForEpisode[quest.Name] = quest
		allQuests[quest.Episode] = questsForEpisode
		questsById[quest.Number] = quest
	}

	return Quests{
		questsById:   questsById,
		allQuests:    allQuests,
		warnedQuests: make(map[string]bool),
	}
}

func (q *Quests) GetQuestConfig(questNumber uint16, episode int, questName string) (Quest, bool) {
	quest, questFound := q.questsById[questNumber]
	if !questFound {
		questsForEpisode, exists := q.allQuests[episode]
		if !exists {
			log.Printf("Episode %v not found?", episode)
			return Quest{}, false
		}
		quest, questFound = questsForEpisode[questName]
	}
	if !questFound {
		if warned := q.warnedQuests[questName]; !warned {
			q.warnedQuests[questName] = true
			log.Printf("Episode %v Quest '%v' number %d not configured", episode, questName, questNumber)
		}
	}
	if questFound && quest.Remap != nil {
		return q.GetQuestConfig(0, episode, *quest.Remap)
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
	return q.ForceTerminal || (!q.StartsAtWarpIn() && !q.StartsOnRegister())
}

func (q *Quest) EndsOnRegister() bool {
	return q.End.Register != nil
}

func (q *Quest) GetCmodeStage() int {
	if q.CmodeStage > 0 {
		return q.CmodeStage
	} else {
		return -1
	}
}
