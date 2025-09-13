package server

import (
	"encoding/json"
	"fmt"
	"github.com/phelix-/psostats/v2/pkg/common"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"github.com/phelix-/psostats/v2/server/internal/userdb"
	"github.com/phelix-/psostats/v2/server/internal/weapons"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
)

type Server struct {
	app                     *fiber.App
	dynamoClient            *dynamodb.DynamoDB
	userDb                  userdb.UserDb
	recentGames             []model.QuestRun
	recentGamesCount        int
	recentGamesSize         int
	recentGamesLock         sync.Mutex
	recordsLock             sync.Mutex
	webhookUrl              string
	adminWebhookUrl         string
	indexTemplate           *template.Template
	infoTemplate            *template.Template
	downloadTemplate        *template.Template
	gameTemplate            *template.Template
	gameV3Template          *template.Template
	gameV4Template          *template.Template
	playerTemplate          *template.Template
	gameNotFoundTemplate    *template.Template
	recordsTemplate         *template.Template
	anniversaryTemplate     *template.Template
	anniversary2022Template *template.Template
	comboCalcTemplate       *template.Template
	anniversaryQuests       map[string]struct{}
	anniversaryNamesInOrder []string
}

func New(dynamo *dynamodb.DynamoDB) *Server {
	f := fiber.New(fiber.Config{
		// modify config
	})
	cacheSize := 500
	webhookUrl, _ := os.LookupEnv("WEBHOOK_URL")
	adminWebhookUrl, _ := os.LookupEnv("ADMIN_WEBHOOK_URL")
	return &Server{
		app:              f,
		dynamoClient:     dynamo,
		userDb:           userdb.DynamoInstance(dynamo),
		recentGames:      make([]model.QuestRun, cacheSize),
		recentGamesCount: 0,
		recentGamesSize:  cacheSize,
		webhookUrl:       webhookUrl,
		adminWebhookUrl:  adminWebhookUrl,
		anniversaryQuests: map[string]struct{}{
			"Maximum Attack E: Forest": {},
			"Maximum Attack E: Caves":  {},
			"Maximum Attack E: Mines":  {},
			"Maximum Attack E: Ruins":  {},
			"Maximum Attack E: Temple": {},
			"Maximum Attack E: Space":  {},
			"Maximum Attack E: CCA":    {},
			"Maximum Attack E: Seabed": {},
			"Maximum Attack E: Tower":  {},
			"Maximum Attack E: Crater": {},
			"Maximum Attack E: Desert": {},
			"August Atrocity #1":       {},
			"August Atrocity #2":       {},
		},
		anniversaryNamesInOrder: []string{
			"Maximum Attack E: Forest",
			"Maximum Attack E: Caves",
			"Maximum Attack E: Mines",
			"Maximum Attack E: Ruins",
			"Maximum Attack E: Temple",
			"Maximum Attack E: Space",
			"Maximum Attack E: CCA",
			"Maximum Attack E: Seabed",
			"Maximum Attack E: Tower",
			"Maximum Attack E: Crater",
			"Maximum Attack E: Desert",
		},
	}
}

func (s *Server) Run() {
	s.app.Static("/favicon.ico", "./static/favicon.ico", fiber.Static{})
	s.app.Static("/static/", "./static/", fiber.Static{})
	s.app.Get("/js/game.js", s.GetGameJs)
	s.app.Get("/js/draughts.js", s.DraughtsJs)
	s.app.Get("/js/three.module.js", s.ThreeJs)
	s.app.Get("/js/OrbitControls.js", s.OrbitControlsJs)

	// UI
	s.app.Get("/", s.Index)
	s.app.Get("/game/:gameId/:gem?", s.GamePageV4)
	s.app.Get("/gamev1/:gameId/:gem?", s.GamePage)
	s.app.Get("/gamev3/:gameId/:gem?", s.GamePageV3)
	s.app.Get("/gamev4/:gameId/:gem?", s.GamePageV4)
	s.app.Get("/info", s.InfoPage)
	s.app.Get("/download", s.DownloadPage)
	s.app.Get("/records", s.RecordsV2Page)
	s.app.Get("/anniv2021", s.Anniv2021RecordsPage)
	s.app.Get("/anniv2022", s.Anniv2022RecordsPage)
	s.app.Get("/anniv2023", s.Anniv2023RecordsPage)
	s.app.Get("/anniv2025", s.Anniv2025RecordsPage)
	//s.app.Get("/threejs", s.ThreejsPage)
	//s.app.Get("/geometry", s.GetGeometry)
	s.app.Get("/combo-calculator", s.ComboCalcMultiPage)
	s.app.Get("/tech-calculator", s.TechCalcPage)
	s.app.Get("/combo-calculator/opm", s.ComboCalcOpmPage)
	s.app.Get("/combo-calculator/ultima", s.ComboCalcUltima)
	s.app.Get("/players/:player", s.PlayerV2Page)
	// API
	s.app.Post("/api/game", s.PostGame)
	s.app.Get("/api/game/:gameId/:gem?", s.GetGame)
	s.app.Get("/api/record/:quest", s.GetRecord)
	s.app.Get("/api/record-splits/:quest", s.GetRecordSplits)
	s.app.Get("/api/pb-splits/:quest", s.GetPbSplits)
	s.app.Get("/api/weapons", s.GetWeapons)
	s.app.Post("/api/motd", s.PostMotd)
	s.app.Post("/api/users/register", s.RegisterUser)
	s.indexTemplate = ensureParsed("./server/internal/templates/index.gohtml")
	s.infoTemplate = ensureParsed("./server/internal/templates/info.gohtml")
	s.playerTemplate = ensureParsed("./server/internal/templates/playerV2.gohtml")
	s.gameTemplate = ensureParsed("./server/internal/templates/game.gohtml")
	s.gameV3Template = ensureParsed("./server/internal/templates/gamev3.gohtml")
	s.gameV4Template = ensureParsed("./server/internal/templates/gamev4.gohtml")
	s.downloadTemplate = ensureParsed("./server/internal/templates/download.gohtml")
	s.gameNotFoundTemplate = ensureParsed("./server/internal/templates/gameNotFound.gohtml")
	s.recordsTemplate = ensureParsed("./server/internal/templates/recordsV2.gohtml")
	s.anniversaryTemplate = ensureParsed("./server/internal/templates/anniv2021.gohtml")
	s.anniversary2022Template = ensureParsed("./server/internal/templates/anniv2022.gohtml")
	s.comboCalcTemplate = ensureParsed("./server/internal/templates/comboCalc.gohtml")

	if certLocation, found := os.LookupEnv("CERT"); found {
		keyLocation := os.Getenv("KEY")
		go http.ListenAndServe(":80", http.HandlerFunc(redirectToTls))
		if err := s.app.ListenTLS(":443", certLocation, keyLocation); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := s.app.Listen(":80"); err != nil {
			log.Fatal(err)
		}
	}
}

func ensureParsed(templatePath string) *template.Template {
	t, err := template.ParseFiles("./server/internal/templates/navbar.gohtml", templatePath)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func redirectToTls(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
}

func (s *Server) Index(c *fiber.Ctx) error {
	games, err := db.GetRecentGames(s.dynamoClient)
	if err != nil {
		log.Printf("get recent games %v", err)
		c.Status(500)
		return err
	}
	recentGamesModel := struct {
		Games []model.FormattedGame
	}{
		Games: make([]model.FormattedGame, len(games)),
	}
	for i, game := range games {
		formattedGame := getFormattedGame(game)
		recentGamesModel.Games[i] = formattedGame
	}
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	err = s.indexTemplate.ExecuteTemplate(c.Response().BodyWriter(), "index", recentGamesModel)
	return err
}

func (s *Server) InfoPage(c *fiber.Ctx) error {
	infoModel := struct{}{}
	err := s.infoTemplate.ExecuteTemplate(c.Response().BodyWriter(), "info", infoModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

type QuestAndCount struct {
	Quest string
	Count int
}

func sortGames(games []model.Game) map[int]map[string]map[string]model.FormattedGame {
	recordModel := make(map[int]map[string]map[string]model.FormattedGame)
	for _, game := range games {
		formattedGame := getFormattedGame(game)
		questsForEpisode := recordModel[game.Episode]
		if questsForEpisode == nil {
			questsForEpisode = make(map[string]map[string]model.FormattedGame)
		}
		gamesForQuest := questsForEpisode[game.Quest]
		if gamesForQuest == nil {
			gamesForQuest = make(map[string]model.FormattedGame)
		}
		gamesForQuest[game.Category] = formattedGame
		questsForEpisode[game.Quest] = gamesForQuest
		recordModel[game.Episode] = questsForEpisode
	}
	return recordModel
}

func (s *Server) DownloadPage(c *fiber.Ctx) error {
	downloadModel := struct{}{}
	err := s.downloadTemplate.ExecuteTemplate(c.Response().BodyWriter(), "download", downloadModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func (s *Server) GamePage(c *fiber.Ctx) error {
	gameId := c.Params("gameId")
	gem, err := strconv.Atoi(c.Params("gem"))
	if err != nil {
		gem = -1
	}
	fullGame, err := db.GetFullGame(gameId, s.dynamoClient)
	if err != nil {
		return err
	}
	game, err := db.GetGame(gameId, gem, s.dynamoClient)
	if err != nil {
		return err
	}

	if game == nil {
		err = s.gameNotFoundTemplate.ExecuteTemplate(c.Response().BodyWriter(), "gameNotFound", nil)
	} else {
		duration, err := time.ParseDuration(game.QuestDuration)
		if err != nil {
			return err
		}

		invincibleRanges := make(map[int]int)
		invincibleStart := -1
		for i, invincible := range game.Invincible {
			if invincible {
				if invincibleStart < 0 {
					invincibleStart = i
				}
			} else {
				if invincibleStart > 0 {
					if i-invincibleStart >= 10 {
						invincibleRanges[invincibleStart] = i
					}
					invincibleStart = -1
				}
			}
		}
		for _, split := range game.Splits {
			if split.StartSecond > 1 {
				game.Events = append(game.Events, model.Event{
					Second:      split.StartSecond,
					Description: split.Name,
				})
			}
		}
		playerIndex := -1
		for i, player := range game.AllPlayers {
			if game.GuildCard == player.GuildCard {
				playerIndex = i
			}
		}
		timeMoving := game.TimeByState[2] + game.TimeByState[4]
		timeStanding := game.TimeByState[1]
		timeAttacking := game.TimeByState[5] + game.TimeByState[6] + game.TimeByState[7]
		timeCasting := game.TimeByState[8]
		totalActions := 0
		for _, weapon := range game.Weapons {
			actions := weapon.Attacks + weapon.Techs
			if actions > totalActions {
				totalActions = actions
			}
		}
		model := struct {
			Game                 model.QuestRun
			HasPov               map[int]bool
			FormattedQuestTime   string
			InvincibleRanges     map[int]int
			HpRanges             map[int]uint16
			TpRanges             map[int]uint16
			PbRanges             map[int]int
			MonstersAliveRanges  map[int]int
			MonstersKilledRanges map[int]int
			MesetaChargedRanges  map[int]int
			FreezeTrapRanges     map[int]uint16
			ShiftaRanges         map[int]int16
			DebandRanges         map[int]int16
			HpPoolRanges         map[int]int
			Weapons              []model.Equipment
			Barriers             []model.Equipment
			Frames               []model.Equipment
			Units                []model.Equipment
			Mags                 []model.Equipment
			VideoUrl             string
			TimeMoving           string
			TimeStanding         string
			TimeAttacking        string
			TimeCasting          uint64
			FormattedTimeCasting string
			MapData              []MapData
			PlayerIndex          int
			TechsInOrder         [][]string
			MostActions          int
		}{
			Game: *game,
			HasPov: map[int]bool{
				0: fullGame.P1Gzip != nil,
				1: fullGame.P2Gzip != nil,
				2: fullGame.P3Gzip != nil,
				3: fullGame.P4Gzip != nil,
			},
			FormattedQuestTime:   formatDuration(duration),
			InvincibleRanges:     invincibleRanges,
			HpRanges:             convertU16ToXY(game.HP),
			TpRanges:             convertU16ToXY(game.TP),
			PbRanges:             convertF32ToXY(game.PB),
			MonstersAliveRanges:  convertIntToXY(game.MonsterCount),
			MonstersKilledRanges: convertIntToXY(game.MonstersKilledCount),
			MesetaChargedRanges:  convertIntToXY(game.MesetaCharged),
			FreezeTrapRanges:     convertU16ToXY(game.FreezeTraps),
			ShiftaRanges:         convertToXY(game.ShiftaLvl),
			DebandRanges:         convertToXY(game.DebandLvl),
			HpPoolRanges:         convertIntToXY(game.MonsterHpPool),
			Weapons:              getEquipment(game, model.EquipmentTypeWeapon),
			Barriers:             getEquipment(game, model.EquipmentTypeBarrier),
			Frames:               getEquipment(game, model.EquipmentTypeFrame),
			Units:                getEquipment(game, model.EquipmentTypeUnit),
			Mags:                 getEquipment(game, model.EquipmentTypeMag),
			TimeMoving:           formatDuration(time.Duration(timeMoving) * (time.Second / 30)),
			TimeStanding:         formatDuration(time.Duration(timeStanding) * (time.Second / 30)),
			TimeAttacking:        formatDuration(time.Duration(timeAttacking) * (time.Second / 30)),
			TimeCasting:          timeCasting,
			FormattedTimeCasting: formatDuration(time.Duration(timeCasting) * (time.Second / 30)),
			MapData:              formatMap(game, game.DataFrames),
			PlayerIndex:          playerIndex,
			TechsInOrder: [][]string{
				{"Foie", "Zonde", "Barta"},
				{"Gifoie", "Gizonde", "Gibarta"},
				{"Rafoie", "Razonde", "Rabarta"},
				{"Grants", "Megid"},
				{"Resta", "Anti", "Reverser"},
				{"Shifta", "Deband", "Ryuker"},
				{"Jellen", "Zalure"}},
			MostActions: totalActions,
		}
		err = s.gameTemplate.ExecuteTemplate(c.Response().BodyWriter(), "game", model)
	}
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

type MapData struct {
	MapName      string
	MapNum       uint16
	MapVariation uint16
	Movement     map[string]Movement
}
type Movement struct {
	Title       string
	Coordinates [][]float32
	Time        []int64
	Color       string
}

func formatMap(game *model.QuestRun, data []model.DataFrame) []MapData {
	allMapData := make([]MapData, 0)
	if data == nil || len(data) == 0 {
		return allMapData
	}
	mapNum := uint16(255)
	mapVariation := uint16(255)
	mapData := MapData{}
	playerIndexByGc := make(map[string]int)
	playerByGc := make(map[string]model.BasePlayerInfo)
	for index, player := range game.AllPlayers {
		playerIndexByGc[player.GuildCard] = index
		playerByGc[player.GuildCard] = player
	}
	for _, frame := range data {
		if frame.Map == 0 || frame.Map == 18 || frame.Map == 45 {
			continue
		}
		if frame.Map != mapNum || frame.MapVariation != mapVariation {
			if len(mapData.Movement) > 0 {
				allMapData = append(allMapData, mapData)
			}
			mapNum = frame.Map
			mapVariation = frame.MapVariation
			mapData = MapData{
				MapNum:       frame.Map,
				MapVariation: frame.MapVariation,
				Movement:     make(map[string]Movement),
				MapName:      common.GetFloorName(mapNum),
			}
		}
		// PlayerLocation is deprecated
		for player, location := range frame.PlayerLocation {
			playerId := fmt.Sprintf("%d", player)
			playerData := mapData.Movement[playerId]
			if playerData.Coordinates == nil {
				playerData.Coordinates = make([][]float32, 0)
				playerData.Time = make([]int64, 0)
				playerData.Title = fmt.Sprintf("Player %d: %v", player+1, game.AllPlayers[player].Name)
			}
			playerData.Coordinates = append(playerData.Coordinates, []float32{location.X / 4, -location.Z / 4})
			playerData.Time = append(playerData.Time, frame.Time*1000)
			mapData.Movement[playerId] = playerData
		}
		// New location info
		for gc, location := range frame.PlayerByGcLocation {
			playerIndex := fmt.Sprintf("%d", playerIndexByGc[gc])
			playerData := mapData.Movement[playerIndex]
			if playerData.Coordinates == nil {
				playerData.Coordinates = make([][]float32, 0)
				playerData.Time = make([]int64, 0)
				playerData.Title = fmt.Sprintf("Player %d: %v", playerIndexByGc[gc]+1, playerByGc[gc].Name)
			}
			playerData.Coordinates = append(playerData.Coordinates, []float32{location.X / 4, -location.Z / 4})
			playerData.Time = append(playerData.Time, frame.Time*1000)
			mapData.Movement[playerIndex] = playerData
		}
		for monster, location := range frame.MonsterLocation {
			monsterId := fmt.Sprintf("%d", monster)
			monsterData := mapData.Movement[monsterId]
			if monsterData.Coordinates == nil {
				monsterData.Coordinates = make([][]float32, 0)
				monsterData.Time = make([]int64, 0)
				monsterData.Title = game.Monsters[monster].Name
			}
			monsterData.Coordinates = append(monsterData.Coordinates, []float32{location.X / 4, -location.Z / 4})
			monsterData.Time = append(monsterData.Time, frame.Time*1000)
			mapData.Movement[monsterId] = monsterData
		}
	}
	if len(mapData.Movement) > 0 {
		allMapData = append(allMapData, mapData)
	}
	return allMapData
}

func getEquipment(game *model.QuestRun, equipmentType string) []model.Equipment {
	equipmentOfType := make([]model.Equipment, 0)
	if game.Weapons != nil && len(game.Weapons) > 0 {
		for _, equipment := range game.Weapons {
			if equipment.Type == equipmentType {
				equipment.Attacks = equipment.Attacks + equipment.Techs // ugly
				equipmentOfType = append(equipmentOfType, equipment)
			}
		}
	} else if game.EquipmentUsedTime != nil && len(game.EquipmentUsedTime) > 0 {
		equipmentUsed := game.EquipmentUsedTime[equipmentType]
		for k, v := range equipmentUsed {
			equipmentOfType = append(equipmentOfType, model.Equipment{Display: k, SecondsEquipped: v})
		}
	}
	sort.Slice(equipmentOfType, func(i, j int) bool {
		return equipmentOfType[i].Display < equipmentOfType[j].Display
	})
	return equipmentOfType
}

func formatDurationSeconds(d time.Duration) string {
	d = d.Round(time.Second)
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	d -= seconds * time.Second
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	d -= seconds * time.Second
	milliseconds := d / time.Millisecond
	return fmt.Sprintf("%d:%02d.%03d", minutes, seconds, milliseconds)
}

func formatDurationSecMilli(d time.Duration) string {
	d = d.Round(time.Millisecond)
	seconds := d / time.Second
	d -= seconds * time.Second
	milliseconds := d / time.Millisecond
	return fmt.Sprintf("%d.%03d", seconds, milliseconds)
}

func convertIntToXY(values []int) map[int]int {
	converted := make(map[int]int)
	previousValue := 0
	for i, value := range values {
		if i == 0 || value != previousValue {
			converted[i] = value
			previousValue = value
		}
	}
	converted[len(values)-1] = previousValue
	return converted
}

func convertU16ToXY(values []uint16) map[int]uint16 {
	converted := make(map[int]uint16)
	previousValue := uint16(0)
	for i, value := range values {
		if i == 0 || value != previousValue {
			converted[i] = value
			previousValue = value
		}
	}
	converted[len(values)-1] = previousValue
	return converted
}

func convertF32ToXY(values []float32) map[int]int {
	converted := make(map[int]int)
	previousValue := 0
	for i, value := range values {
		intValue := int(value)
		if i == 0 || intValue != previousValue {
			converted[i] = intValue
			previousValue = intValue
		}
	}
	converted[len(values)-1] = previousValue
	return converted
}

func convertToXY(values []int16) map[int]int16 {
	converted := make(map[int]int16)
	previousValue := int16(0)
	for i, value := range values {
		if i == 0 || value != previousValue {
			converted[i] = value
			previousValue = value
		}
	}
	converted[len(values)-1] = previousValue
	return converted
}

func (s *Server) RecordsV2Page(c *fiber.Ctx) error {
	games, err := db.GetQuestRecords(db.QuestRecordsTable, s.dynamoClient)
	if err != nil {
		return err
	}
	recordModel := sortGames(games)

	err = s.recordsTemplate.ExecuteTemplate(c.Response().BodyWriter(), "index", recordModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func (s *Server) Anniv2021RecordsPage(c *fiber.Ctx) error {
	recordModel := struct{}{}

	err := s.anniversaryTemplate.ExecuteTemplate(c.Response().BodyWriter(), "index", recordModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func getFormattedGame(game model.Game) model.FormattedGame {
	shortCategory := game.Category
	numPlayers, err := strconv.Atoi(string(shortCategory[0]))
	if err != nil {
		log.Fatalf("Couldn't atoi")
	}
	pbRun := string(shortCategory[1])
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatalf("Couldn't find time zone America/Chicago")
	}
	players := make([]model.FormattedPlayerInfo, 4)
	for i := range game.PlayerNames {
		var hasPov bool
		switch i {
		case 0:
			hasPov = game.P1HasStats
		case 1:
			hasPov = game.P2HasStats
		case 2:
			hasPov = game.P3HasStats
		case 3:
			hasPov = game.P4HasStats
		}
		players[i] = model.FormattedPlayerInfo{
			Name:      game.PlayerNames[i],
			GuildCard: game.PlayerGcs[i],
			HasPov:    hasPov,
			Class:     game.PlayerClasses[i],
		}
	}
	var formattedRelativeDate string
	relativeDate := time.Now().Sub(game.Timestamp)
	if relativeDate > time.Hour*24 {
		daysAgo := relativeDate / (time.Hour * 24)
		if daysAgo > 1 {
			formattedRelativeDate = fmt.Sprintf("%d days ago", daysAgo)
		} else {
			formattedRelativeDate = "A day ago"
		}
	} else if relativeDate > time.Hour {
		hoursAgo := relativeDate / time.Hour
		if hoursAgo > 1 {
			formattedRelativeDate = fmt.Sprintf("%d hours ago", hoursAgo)
		} else {
			formattedRelativeDate = "An hour ago"
		}
	} else {
		minutesAgo := relativeDate / time.Minute
		if minutesAgo > 2 {
			formattedRelativeDate = fmt.Sprintf("%d minutes ago", minutesAgo)
		} else {
			formattedRelativeDate = "Just now"
		}
	}
	return model.FormattedGame{
		Id:           game.Id,
		Players:      players,
		PbRun:        pbRun == "p",
		NumPlayers:   numPlayers,
		Episode:      game.Episode,
		Quest:        game.Quest,
		Duration:     game.Time,
		Time:         formatDuration(game.Time),
		RelativeDate: formattedRelativeDate,
		Date:         game.Timestamp.In(location).Format("15:04 01/02/2006"),
	}
}

func (s *Server) GetGame(c *fiber.Ctx) error {
	gameId := c.Params("gameId")
	gem := c.Params("gem")
	gemInt, err := strconv.Atoi(gem)
	if err != nil {
		gemInt = -1
	}
	game, _ := db.GetGame(gameId, gemInt, s.dynamoClient)

	if game == nil {
		c.Status(404)
		return nil
	} else {
		log.Printf("Serving %v", gameId)
		compressed, _ := db.Compress(game)
		c.Response().AppendBody(compressed)
		c.Response().Header.Set("Content-Type", "application/json")
		c.Response().Header.Set("Content-Encoding", "gzip")
		return nil
	}
}

func (s *Server) GetRecord(c *fiber.Ctx) error {
	questName := c.Params("quest")
	questName, err := url.PathUnescape(questName)
	if err != nil {
		return err
	}
	playersString := c.Query("players", "4")
	pbString := c.Query("pb", "false")
	pbCategory := strings.ToLower(pbString) == "true"
	numPlayers, err := strconv.Atoi(playersString)
	if err != nil {
		return err
	}
	questRecord, err := db.GetQuestRecord(questName, numPlayers, pbCategory, s.dynamoClient)
	if questRecord == nil {
		c.Status(404)
		return nil
	} else {
		compressed, _ := db.Compress(questRecord)
		c.Response().AppendBody(compressed)
		c.Response().Header.Set("Content-Type", "application/json")
		c.Response().Header.Set("Content-Encoding", "gzip")
		return nil
	}
}

func (s *Server) GetRecordSplits(c *fiber.Ctx) error {
	questName := c.Params("quest")
	questName, err := url.PathUnescape(questName)
	if err != nil {
		return err
	}
	playersString := c.Query("players", "4")
	pbString := c.Query("pb", "false")
	pbCategory := strings.ToLower(pbString) == "true"
	numPlayers, err := strconv.Atoi(playersString)
	if err != nil {
		return err
	}
	questRecord, err := db.GetQuestRecord(questName, numPlayers, pbCategory, s.dynamoClient)
	if err != nil {
		return err
	}

	if questRecord == nil {
		c.Status(404)
		return nil
	} else {
		game, err := db.GetGame(questRecord.Id, -1, s.dynamoClient)
		if err != nil {
			return err
		}
		if game == nil || game.Splits == nil || len(game.Splits) == 0 {
			c.Status(404)
			return nil
		} else {
			game.Splits[len(game.Splits)-1].End = game.QuestEndTime
			splitBytes, _ := json.Marshal(game.Splits)
			c.Response().AppendBody(splitBytes)
			c.Response().Header.Set("Content-Type", "application/json")
			return nil
		}
	}
}

func (s *Server) GetPbSplits(c *fiber.Ctx) error {
	authorized, user := s.verifyAuth(&c.Request().Header)
	if !authorized {
		c.Status(401)
		return nil
	}
	questName := c.Params("quest")
	questName, err := url.PathUnescape(questName)
	if err != nil {
		return err
	}
	playersString := c.Query("players", "4")
	pbString := c.Query("pb", "false")
	pbCategory := strings.ToLower(pbString) == "true"
	numPlayers, err := strconv.Atoi(playersString)
	if err != nil {
		return err
	}
	questRecord, err := db.GetPlayerPB(questName, user.Id, numPlayers, pbCategory, s.dynamoClient)
	if err != nil {
		return err
	}

	if questRecord == nil {
		c.Status(404)
		return nil
	} else {
		game, err := db.GetGame(questRecord.Id, -1, s.dynamoClient)
		if err != nil {
			return err
		}
		if game == nil || game.Splits == nil || len(game.Splits) == 0 {
			c.Status(404)
			return nil
		} else {
			game.Splits[len(game.Splits)-1].End = game.QuestEndTime
			return respondWithJson(game.Splits, c)
		}
	}
}

func (s *Server) RegisterUser(c *fiber.Ctx) error {
	authorized, user := s.verifyAuth(&c.Request().Header)
	if !authorized {
		c.Status(401)
		return nil
	}
	if !user.Admin {
		c.Status(403)
		return nil
	}
	var newUser userdb.User
	if err := c.BodyParser(&newUser); err != nil {
		c.Status(400)
		return nil
	}
	newUser.Admin = false
	newUser.Password = HashPassword(newUser.Password)

	userForId, err := s.userDb.GetUserByDiscordId(newUser.DiscordId)
	if err != nil {
		return err
	}
	if userForId != nil {
		c.Status(400)
		c.Response().AppendBodyString("User already associated with this discord account")
		c.Response().Header.Set("Content-Type", "plain/text")
		return nil
	}
	existingUser, err := s.userDb.GetUser(newUser.Id)
	if err != nil {
		return err
	}
	if existingUser != nil {
		c.Status(400)
		c.Response().AppendBodyString("User already exists")
		c.Response().Header.Set("Content-Type", "plain/text")
		return nil
	}

	err = s.userDb.CreateUser(newUser)
	if err != nil {
		return err
	}

	s.SendWebhook(Webhook{Embeds: []Embed{{Title: "New User Registered: " + newUser.Id}}}, s.adminWebhookUrl)
	return nil
}

func (s *Server) GetWeapons(c *fiber.Ctx) error {
	return respondWithJson(weapons.GetWeapons(), c)
}

func respondWithJson(data any, c *fiber.Ctx) error {
	jsonBytes, err := db.Compress(data)
	if err != nil {
		return err
	}
	c.Response().AppendBody(jsonBytes)
	c.Response().Header.Set("Content-Type", "application/json")
	c.Response().Header.Set("Content-Encoding", "gzip")
	return nil
}

func (s *Server) PostMotd(c *fiber.Ctx) error {
	authorized, user := s.verifyAuth(&c.Request().Header)
	var clientInfo model.ClientInfo
	if err := c.BodyParser(&clientInfo); err != nil {
		log.Printf("body parser")
		c.Status(400)
		return err
	}
	message := "Not logged in"
	if authorized && user != nil {
		message = fmt.Sprintf("Logged in as %v, up to date", user.Id)
	}
	if getClientVersionInt(clientInfo) < 11200 {
		message = "Update available. https://psostats.com/download"
	}
	motd := model.MessageOfTheDay{
		Authorized: authorized,
		Message:    message,
	}
	jsonBytes, err := json.Marshal(motd)
	if err != nil {
		return err
	}
	c.Response().AppendBody(jsonBytes)
	c.Response().Header.Set("Content-Type", "application/json")
	return nil
}

func GamesMatch(a, b model.QuestRun) bool {
	if a.QuestName != b.QuestName {
		return false
	}
	if a.Difficulty != b.Difficulty {
		return false
	}
	if a.Episode != b.Episode {
		return false
	}
	if a.Server != b.Server {
		return false
	}
	if a.GuildCard == b.GuildCard {
		return false
	}
	if a.UserName == b.UserName {
		return false
	}
	if a.SubmittedTime.Add(time.Second*-30).After(b.SubmittedTime) ||
		a.SubmittedTime.Add(time.Second*30).Before(b.SubmittedTime) {
		return false
	}
	if len(a.AllPlayers) != len(b.AllPlayers) {
		return false
	}
	for i := range a.AllPlayers {
		if a.AllPlayers[i].Name != b.AllPlayers[i].Name {
			return false
		}
		if a.AllPlayers[i].GuildCard != b.AllPlayers[i].GuildCard {
			return false
		}
		if a.AllPlayers[i].Class != b.AllPlayers[i].Class {
			return false
		}
	}
	return true
}

func getClientVersionInt(clientInfo model.ClientInfo) int {
	return clientInfo.VersionMajor*10000 + clientInfo.VersionMinor*100 + clientInfo.VersionPatch
}
