package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"github.com/phelix-/psostats/v2/server/internal/userdb"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
)

type Server struct {
	app              *fiber.App
	dynamoClient     *dynamodb.DynamoDB
	userDb           userdb.UserDb
	recentGames      []model.QuestRun
	recentGamesCount int
	recentGamesSize  int
	recentGamesLock sync.Mutex
}

func New(dynamo *dynamodb.DynamoDB) *Server {
	f := fiber.New(fiber.Config{
		// modify config
	})
	cacheSize := 500
	return &Server{
		app:              f,
		dynamoClient:     dynamo,
		userDb:           userdb.DynamoInstance(dynamo),
		recentGames:      make([]model.QuestRun, cacheSize),
		recentGamesCount: 0,
		recentGamesSize:  cacheSize,
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
	s.app.Get("/recordsV2", s.RecordsPage)
	s.app.Get("/players/:player", s.PlayerPage)
	s.app.Get("/playersV2", s.PlayerV2Page)
	s.app.Get("/gc/:gc", s.GcRedirect)
	// API
	s.app.Post("/api/game", s.PostGame)
	s.app.Get("/api/game/:gameId", s.GetGame)
	s.app.Post("/api/motd", s.PostMotd)

	if certLocation, found := os.LookupEnv("CERT"); found {
		keyLocation := os.Getenv("KEY")
		if err := s.app.ListenTLS(":443", certLocation, keyLocation); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := s.app.Listen(":80"); err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Server) Index(c *fiber.Ctx) error {
	t, err := template.ParseFiles("./server/internal/templates/index.gohtml")
	if err != nil {
		c.Status(500)
		return err
	}
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
	err = t.ExecuteTemplate(c.Response().BodyWriter(), "index", recentGamesModel)
	return err
}
func (s *Server) InfoPage(c *fiber.Ctx) error {
	t, err := template.ParseFiles("./server/internal/templates/info.gohtml")
	if err != nil {
		return err
	}
	infoModel := struct{}{}
	err = t.ExecuteTemplate(c.Response().BodyWriter(), "info", infoModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func (s *Server) PlayerV2Page(c *fiber.Ctx) error {
	t, err := template.ParseFiles("./server/internal/templates/playerV2.gohtml")
	if err != nil {
		return err
	}
	infoModel := struct{}{}
	err = t.ExecuteTemplate(c.Response().BodyWriter(), "player", infoModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func (s *Server) DownloadPage(c *fiber.Ctx) error {
	t, err := template.ParseFiles("./server/internal/templates/download.gohtml")
	if err != nil {
		return err
	}
	downloadModel := struct{}{}
	err = t.ExecuteTemplate(c.Response().BodyWriter(), "download", downloadModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func (s *Server) GamePage(c *fiber.Ctx) error {
	gameId := c.Params("gameId")
	gem := c.Params("gem")
	fullGame, err := db.GetFullGame(gameId, s.dynamoClient)
	if err != nil {
		return err
	}
	var gameGzip []byte
	var videoUrl string
	if fullGame != nil {
		if len(gem) > 0 {
			gemNum, err := strconv.Atoi(gem)
			if err == nil {
				switch gemNum + 1 {
				case 1:
					gameGzip = fullGame.P1Gzip
					videoUrl = fullGame.P1Video
				case 2:
					gameGzip = fullGame.P2Gzip
					videoUrl = fullGame.P2Video
				case 3:
					gameGzip = fullGame.P3Gzip
					videoUrl = fullGame.P3Video
				case 4:
					gameGzip = fullGame.P4Gzip
					videoUrl = fullGame.P4Video
				}
			}
		}
		if gameGzip == nil {
			gameGzip = fullGame.GameGzip
		}
	}

	if gameGzip == nil {
		t, err := template.ParseFiles("./server/internal/templates/gameNotFound.gohtml")
		if err != nil {
			return err
		}
		err = t.ExecuteTemplate(c.Response().BodyWriter(), "gameNotFound", nil)
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
		}
		t, err := template.ParseFiles("./server/internal/templates/game.gohtml")
		if err != nil {
			return err
		}
		err = t.ExecuteTemplate(c.Response().BodyWriter(), "game", model)
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
	t, err := template.ParseFiles("./server/internal/templates/recordsV2.gohtml")
	if err != nil {
		return err
	}

	games, err := db.GetQuestRecords(s.dynamoClient)
	if err != nil {
		return err
	}
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

	err = t.ExecuteTemplate(c.Response().BodyWriter(), "index", recordModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func (s *Server) RecordsPage(c *fiber.Ctx) error {
	t, err := template.ParseFiles("./server/internal/templates/records.gohtml")
	if err != nil {
		c.Status(500)
		return err
	}
	games, err := db.GetQuestRecords(s.dynamoClient)
	sort.Slice(games, func(i, j int) bool {
		if games[i].Episode != games[j].Episode {
			return games[i].Episode < games[j].Episode
		}
		if games[i].Quest != games[j].Quest {
			return games[i].Quest < games[j].Quest
		}
		return games[i].Category < games[j].Category
	})

	if err != nil {
		log.Print("get recent games")
		c.Status(500)
		return err
	}
	for i, game := range games {
		addFormattedFields(&game)
		games[i] = game
	}
	model := struct {
		Games []model.Game
	}{
		Games: games,
	}
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	err = t.ExecuteTemplate(c.Response().BodyWriter(), "index", model)
	return err
}

func addFormattedFields(game *model.Game) {
	game.FormattedTime = formatDuration(game.Time)
	shortCategory := game.Category
	numPlayers := string(shortCategory[0])
	pbRun := string(shortCategory[1])
	pbText := ""
	if pbRun == "p" {
		pbText = " PB"
	}
	game.Category = numPlayers + " Player" + pbText
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatalf("Couldn't find time zone America/Chicago")
	}
	game.FormattedDate = game.Timestamp.In(location).Format("15:04 01/02/2006")
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

func (s *Server) PlayerPage(c *fiber.Ctx) error {
	player := c.Params("player")
	t, err := template.ParseFiles("./server/internal/templates/player.gohtml")
	if err != nil {
		c.Status(500)
		return err
	}

	player, err = url.PathUnescape(player)
	if err != nil {
		c.Status(500)
		return err
	}
	pbs, err := db.GetPlayerPbs(player, s.dynamoClient)

	if err != nil {
		log.Print("get pbs")
		c.Status(500)
		return err
	}
	sort.Slice(pbs, func(i, j int) bool { return pbs[i].Quest < pbs[j].Quest })
	for i, game := range pbs {
		addFormattedFields(&game)
		pbs[i] = game
	}

	recent, err := db.GetPlayerRecentGames(player, s.dynamoClient)

	if err != nil {
		log.Print("get recent")
		c.Status(500)
		return err
	}
	for i, game := range recent {
		addFormattedFields(&game)
		recent[i] = game
	}

	model := struct {
		Player      string
		PlayerPbs   []model.Game
		RecentGames []model.Game
	}{
		Player:      player,
		PlayerPbs:   pbs,
		RecentGames: recent,
	}
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	err = t.ExecuteTemplate(c.Response().BodyWriter(), "index", model)
	return err
}

func (s *Server) GcRedirect(c *fiber.Ctx) error {
	gc := c.Params("gc")
	playerName, err := s.userDb.GetUsernameByGc(gc)
	if err != nil {
		log.Printf("loading player by gc %v %v", gc, err)
	}
	return c.Redirect(fmt.Sprintf("/players/%v", playerName))
}

func (s *Server) GetGame(c *fiber.Ctx) error {
	gameId := c.Params("gameId")
	game, err := db.GetGame(gameId, s.dynamoClient)
	if err != nil {
		return err
	}

	if game == nil {
		c.Status(404)
		return nil
	} else {
		jsonBytes, err := json.Marshal(game)
		if err != nil {
			return err
		}
		c.Response().AppendBody(jsonBytes)
		c.Response().Header.Set("Content-Type", "application/json")
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
	if clientInfo.VersionMajor < 0 || clientInfo.VersionMinor < 7 || clientInfo.VersionPatch < 2 {
		message = "Version 0.7.2 available! Head to https://psostats.com/download to update"
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

func IsLeaderboardCandidate(questRun model.QuestRun) bool {
	cmodeRegex := regexp.MustCompile("[12]c\\d")
	allowedDifficulty := questRun.Difficulty == "Ultimate" || cmodeRegex.MatchString(questRun.QuestName)
	return allowedDifficulty && questRun.QuestComplete && !questRun.IllegalShifta
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
	if a.QuestStartTime.Add(time.Second*-30).After(b.QuestStartTime) ||
		a.QuestStartTime.Add(time.Second*30).Before(b.QuestStartTime) {
		return false
	}
	if a.QuestEndTime.Add(time.Second*-30).After(b.QuestEndTime) ||
		a.QuestEndTime.Add(time.Second*30).Before(b.QuestEndTime) {
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
