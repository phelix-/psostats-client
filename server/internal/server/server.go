package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/phelix-/psostats/v2/pkg/common"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"github.com/phelix-/psostats/v2/server/internal/enemies"
	"github.com/phelix-/psostats/v2/server/internal/userdb"
	"github.com/phelix-/psostats/v2/server/internal/weapons"
	"io"
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
	app                  *fiber.App
	dynamoClient         *dynamodb.DynamoDB
	userDb               userdb.UserDb
	recentGames          []model.QuestRun
	recentGamesCount     int
	recentGamesSize      int
	recentGamesLock      sync.Mutex
	recordsLock          sync.Mutex
	webhookUrl           string
	indexTemplate        *template.Template
	infoTemplate         *template.Template
	downloadTemplate     *template.Template
	gameTemplate         *template.Template
	playerTemplate       *template.Template
	gameNotFoundTemplate *template.Template
	recordsTemplate      *template.Template
	anniversaryTemplate *template.Template
	comboCalcTemplate    *template.Template
	anniversaryQuests    map[string]struct{}
}

func New(dynamo *dynamodb.DynamoDB) *Server {
	f := fiber.New(fiber.Config{
		// modify config
	})
	cacheSize := 500
	webhookUrl, _ := os.LookupEnv("WEBHOOK_URL")
	return &Server{
		app:              f,
		dynamoClient:     dynamo,
		userDb:           userdb.DynamoInstance(dynamo),
		recentGames:      make([]model.QuestRun, cacheSize),
		recentGamesCount: 0,
		recentGamesSize:  cacheSize,
		webhookUrl:       webhookUrl,
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
		},
	}
}

func (s *Server) Run() {
	s.app.Static("/favicon.ico", "./static/favicon.ico", fiber.Static{})
	s.app.Static("/static/", "./static/", fiber.Static{})
	// UI
	s.app.Get("/", s.Index)
	s.app.Get("/game/:gameId/:gem?", s.GamePage)
	s.app.Get("/info", s.InfoPage)
	s.app.Get("/download", s.DownloadPage)
	s.app.Get("/records", s.RecordsV2Page)
	s.app.Get("/anniv2021", s.Anniv2021RecordsPage)
	s.app.Get("/combo-calculator", s.ComboCalcMultiPage)
	s.app.Get("/combo-calculator/opm", s.ComboCalcOpmPage)
	s.app.Get("/players/:player", s.PlayerV2Page)
	// API
	s.app.Post("/api/game", s.PostGame)
	s.app.Get("/api/game/:gameId/:gem?", s.GetGame)
	s.app.Post("/api/motd", s.PostMotd)
	s.indexTemplate = ensureParsed("./server/internal/templates/index.gohtml")
	s.infoTemplate = ensureParsed("./server/internal/templates/info.gohtml")
	s.playerTemplate = ensureParsed("./server/internal/templates/playerV2.gohtml")
	s.gameTemplate = ensureParsed("./server/internal/templates/game.gohtml")
	s.downloadTemplate = ensureParsed("./server/internal/templates/download.gohtml")
	s.gameNotFoundTemplate = ensureParsed("./server/internal/templates/gameNotFound.gohtml")
	s.recordsTemplate = ensureParsed("./server/internal/templates/recordsV2.gohtml")
	s.anniversaryTemplate = ensureParsed("./server/internal/templates/anniv2021.gohtml")
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

func (s *Server) ComboCalcOpmPage(c *fiber.Ctx) error {
	return s.comboCalcPage(true, c)
}

func (s *Server) ComboCalcMultiPage(c *fiber.Ctx) error {
	return s.comboCalcPage(false, c)
}

func (s *Server) comboCalcPage(opm bool, c *fiber.Ctx) error {
	sortedEnemies := make(map[string][]enemies.Enemy)

	var allEnemies []enemies.Enemy
	if opm {
		allEnemies = enemies.GetEnemiesUltOpm()
	} else {
		allEnemies = enemies.GetEnemiesUltMulti()
	}
	for _, enemy := range allEnemies {
		enemiesInArea := sortedEnemies[enemy.Location]
		if enemiesInArea == nil {
			enemiesInArea = make([]enemies.Enemy, 0)
		}
		enemiesInArea = append(enemiesInArea, enemy)
		sortedEnemies[enemy.Location] = enemiesInArea
	}
	infoModel := struct {
		Opm     bool
		Classes []weapons.PsoClass
		Enemies map[string][]enemies.Enemy
		Weapons []weapons.Weapon
	}{
		Opm:     opm,
		Classes: weapons.GetClasses(),
		Enemies: sortedEnemies,
		Weapons: weapons.GetWeapons(),
	}
	err := s.comboCalcTemplate.ExecuteTemplate(c.Response().BodyWriter(), "combo-calc", infoModel)

	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func (s *Server) PlayerV2Page(c *fiber.Ctx) error {
	player := c.Params("player")
	player, err := url.PathUnescape(player)
	if err != nil {
		c.Status(500)
		return err
	}
	recentGames, err := db.GetPlayerRecentGames(player, s.dynamoClient, 15)
	if err != nil {
		return err
	}
	games, err := db.GetQuestRecords(db.QuestRecordsTable, s.dynamoClient)
	if err != nil {
		return err
	}
	sortedRecords := sortGames(games)
	playerPbs, err := db.GetPlayerPbs(player, s.dynamoClient)
	if err != nil {
		return err
	}
	sortedPbs := sortGames(playerPbs)

	for _, pb := range playerPbs {
		if sortedRecords[pb.Episode][pb.Quest][pb.Category].Id == pb.Id {
			pbsForEp := sortedPbs[pb.Episode]
			pbsForQuest := pbsForEp[pb.Quest]
			pbForCategory := pbsForQuest[pb.Category]
			pbForCategory.Record = true
			pbsForQuest[pb.Category] = pbForCategory
			pbsForEp[pb.Quest] = pbsForQuest
			sortedPbs[pb.Episode] = pbsForEp
		}
	}
	classUsage, err := db.GetPlayerClassCounts(player, s.dynamoClient)
	if err != nil {
		return err
	}
	for _, class := range common.GetAllClasses() {
		if _, exists := classUsage[class.Name]; !exists {
			classUsage[class.Name] = 0
		}
	}
	questsPlayed, err := db.GetPlayerQuestCounts(player, s.dynamoClient)
	if err != nil {
		return err
	}
	gamesByEpisode := map[int]int{1: 0, 2: 0, 4: 0}
	questCounts := make([]QuestAndCount, 0)
	totalGames := 0
	for quest, count := range questsPlayed {
		split := strings.SplitN(quest, "_", 2)
		if len(split) == 2 {
			episode, err := strconv.Atoi(split[0])
			questName := split[1]
			if err == nil {
				totalGames += count
				gamesByEpisode[episode] += count
				questCounts = append(questCounts, QuestAndCount{questName, count})
			}
		}
	}
	favoriteQuest := QuestAndCount{"None", 0}
	if len(questCounts) > 0 {
		sort.Slice(questCounts, func(i, j int) bool { return questCounts[i].Count > questCounts[j].Count })
		favoriteQuest = questCounts[0]
	}

	infoModel := struct {
		PlayerName     string
		Classes        map[string]int
		TotalGames     int
		GamesByEpisode map[int]int
		FavoriteQuest  QuestAndCount
		RecentGames    []model.FormattedGame
		PbGames        map[int]map[string]map[string]model.FormattedGame
	}{
		PlayerName:     player,
		Classes:        classUsage,
		TotalGames:     totalGames,
		GamesByEpisode: gamesByEpisode,
		FavoriteQuest:  favoriteQuest,
		RecentGames:    make([]model.FormattedGame, 0),
		PbGames:        sortedPbs,
	}
	for _, game := range recentGames {
		formattedGame := getFormattedGame(game)
		pbForQuestAndCategory := sortedPbs[game.Episode][game.Quest][game.Category]
		if pbForQuestAndCategory.Id == game.Id {
			formattedGame.Pb = true
			if pbForQuestAndCategory.Record {
				formattedGame.Record = true
			}
		}
		infoModel.RecentGames = append(infoModel.RecentGames, formattedGame)
	}
	err = s.playerTemplate.ExecuteTemplate(c.Response().BodyWriter(), "player", infoModel)
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

func getGameGzip(game *model.Game, gem string) ([]byte, string) {
	var gameGzip []byte
	var videoUrl string
	if game != nil {
		if len(gem) > 0 {
			if gemNum, err := strconv.Atoi(gem); err == nil {
				switch gemNum + 1 {
				case 1:
					gameGzip = game.P1Gzip
					videoUrl = game.P1Video
				case 2:
					gameGzip = game.P2Gzip
					videoUrl = game.P2Video
				case 3:
					gameGzip = game.P3Gzip
					videoUrl = game.P3Video
				case 4:
					gameGzip = game.P4Gzip
					videoUrl = game.P4Video
				}
			}
		}
		if gameGzip == nil {
			gameGzip = game.GameGzip
		}
	}
	return gameGzip, videoUrl
}

func (s *Server) GamePage(c *fiber.Ctx) error {
	gameId := c.Params("gameId")
	gem := c.Params("gem")
	fullGame, err := db.GetFullGame(gameId, s.dynamoClient)
	if err != nil {
		return err
	}
	gameGzip, videoUrl := getGameGzip(fullGame, gem)

	if gameGzip == nil {
		err = s.gameNotFoundTemplate.ExecuteTemplate(c.Response().BodyWriter(), "gameNotFound", nil)
	} else {
		game, err := parseGameGzip(gameGzip)
		if err != nil {
			return err
		}
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
		timeMoving := game.TimeByState[2] + game.TimeByState[4]
		timeStanding := game.TimeByState[1]
		timeAttacking := game.TimeByState[5] + game.TimeByState[6] + game.TimeByState[7]
		timeCasting := game.TimeByState[8]
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
			VideoUrl:             videoUrl,
			TimeMoving:           formatDuration(time.Duration(timeMoving) * (time.Second / 30)),
			TimeStanding:         formatDuration(time.Duration(timeStanding) * (time.Second / 30)),
			TimeAttacking:        formatDuration(time.Duration(timeAttacking) * (time.Second / 30)),
			TimeCasting:          timeCasting,
			FormattedTimeCasting: formatDuration(time.Duration(timeCasting) * (time.Second / 30)),
		}
		err = s.gameTemplate.ExecuteTemplate(c.Response().BodyWriter(), "game", model)
	}
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func getEquipment(game *model.QuestRun, equipmentType string) []model.Equipment {
	equipmentOfType := make([]model.Equipment, 0)
	if game.Weapons != nil && len(game.Weapons) > 0 {
		for _, equipment := range game.Weapons {
			if equipment.Type == equipmentType {
				equipmentOfType = append(equipmentOfType, equipment)
			}
		}
	} else if game.EquipmentUsedTime != nil && len(game.EquipmentUsedTime) > 0 {
		equipmentUsed := game.EquipmentUsedTime[equipmentType]
		for k, v := range equipmentUsed {
			equipmentOfType = append(equipmentOfType, model.Equipment{Display: k, SecondsEquipped: v})
		}
	}
	return equipmentOfType
}

func parseGameGzip(gameBytes []byte) (*model.QuestRun, error) {
	questRun := model.QuestRun{}
	buffer := bytes.NewBuffer(gameBytes)
	reader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := io.ReadAll(reader)
	if err != io.ErrUnexpectedEOF {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &questRun)
	if err != nil {
		return nil, err
	}

	return &questRun, err
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
		Time:         formatDuration(game.Time),
		RelativeDate: formattedRelativeDate,
		Date:         game.Timestamp.In(location).Format("15:04 01/02/2006"),
	}
}

func (s *Server) GetGame(c *fiber.Ctx) error {
	gameId := c.Params("gameId")
	gem := c.Params("gem")
	fullGame, err := db.GetFullGame(gameId, s.dynamoClient)
	if err != nil {
		return err
	}
	gameGzip, _ := getGameGzip(fullGame, gem)

	if gameGzip == nil {
		c.Status(404)
		return nil
	} else {
		c.Response().AppendBody(gameGzip)
		c.Response().Header.Set("Content-Type", "application/json")
		c.Response().Header.Set("Content-Encoding", "gzip")
		return nil
	}
}

func (s *Server) PostMotd(c *fiber.Ctx) error {
	authorized, user := s.verifyAuth(&c.Request().Header)
	var clientInfo model.ClientInfo
	if err := c.BodyParser(&clientInfo); err != nil {
		log.Printf("body parser")
		c.Status(400)
		return err
	}
	message := fmt.Sprintf("Logged in as %v, up to date", user)
	if clientInfo.VersionMajor < 1 || clientInfo.VersionMinor < 1 || clientInfo.VersionPatch < 0 {
		message = "Update available: fixes MAE Cave and Mine. https://psostats.com/download"
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
		if a.AllPlayers[i] != b.AllPlayers[i] {
			return false
		}
	}
	return true
}
