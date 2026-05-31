# Quest Pages Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add per-quest pages at `/quest/{questName}` reachable via links on the records page, showing top records and a player search that returns the player's complete run history for that quest.

**Architecture:** Two new DB functions query existing DynamoDB tables (`quest_records` and `recent_games_by_player_2`). A new Fiber handler renders a new Go template. The records page template gets a one-line change to link quest names. No new DynamoDB tables required.

**Tech Stack:** Go, Fiber v2, AWS SDK v1 (`github.com/aws/aws-sdk-go`), `text/template`, Bootstrap 5, DynamoDB Local (`amazon/dynamodb-local`) for integration tests.

---

## File Map

| Action | File | Purpose |
|--------|------|---------|
| Modify | `server/internal/db/db.go` | Add `GetQuestRecordsForQuest` and `GetAllGamesForPlayerQuest` |
| Modify | `server/internal/db/db_test.go` | Add table helpers + integration tests for both new DB functions |
| Modify | `server/internal/server/server.go` | Add FuncMap to `ensureParsed`, add `questTemplate` field, register route |
| Create | `server/internal/server/display_quest.go` | `QuestPage` handler |
| Create | `server/internal/templates/questPage.gohtml` | Quest page template |
| Modify | `server/internal/templates/recordsV2.gohtml` | Wrap quest names in links |

---

## Task 1: Start DynamoDB Local

**Files:** none (environment setup)

- [ ] **Step 1: Start DynamoDB Local**

```bash
cd /Users/lukebrutvan/Desktop/psostats-client/db_tests
docker compose up -d
```

- [ ] **Step 2: Verify it's running**

```bash
aws dynamodb list-tables --endpoint-url http://localhost:8000 --region us-west-2
```

Expected output:
```json
{
    "TableNames": []
}
```

If `aws` CLI is not available, verify with curl:
```bash
curl -s http://localhost:8000
```
Expected: any non-connection-refused response.

---

## Task 2: DB — `GetQuestRecordsForQuest`

**Files:**
- Modify: `server/internal/db/db_test.go`
- Modify: `server/internal/db/db.go`

- [ ] **Step 1: Add table creation helpers to `db_test.go`**

Add these two functions at the bottom of `server/internal/db/db_test.go`:

```go
func CreateQuestRecordsTable(dynamoClient *dynamodb.DynamoDB) error {
	_, err := dynamoClient.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(db.QuestRecordsTable),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("Quest"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("Category"), AttributeType: aws.String("S")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("Quest"), KeyType: aws.String(dynamodb.KeyTypeHash)},
			{AttributeName: aws.String("Category"), KeyType: aws.String(dynamodb.KeyTypeRange)},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	return err
}

func writeTestQuestRecord(t *testing.T, dynamoClient *dynamodb.DynamoDB, questName, category, gameId string) {
	t.Helper()
	_, err := dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(db.QuestRecordsTable),
		Item: map[string]*dynamodb.AttributeValue{
			"Quest":    {S: aws.String(questName)},
			"Category": {S: aws.String(category)},
			"Id":       {S: aws.String(gameId)},
		},
	})
	if err != nil {
		t.Fatalf("failed to write test quest record: %v", err)
	}
}
```

Also add to `CreateAllTables` after the existing `if _, exists := tables[db.RecentGamesByMonth]` block:

```go
	if _, exists := tables[db.QuestRecordsTable]; !exists {
		if err = CreateQuestRecordsTable(dynamoClient); err != nil {
			return err
		}
	}
```

- [ ] **Step 2: Write the failing test in `db_test.go`**

```go
func TestGetQuestRecordsForQuest(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		t.Fatal(err)
	}
	dynamoClient := dynamodb.New(sess)
	if err = CreateAllTables(dynamoClient); err != nil {
		t.Fatal(err)
	}

	const quest = "Sweep-up Operation #1"
	writeTestQuestRecord(t, dynamoClient, quest, "1n", "game-1")
	writeTestQuestRecord(t, dynamoClient, quest, "2n", "game-2")
	writeTestQuestRecord(t, dynamoClient, "Some Other Quest", "1n", "game-3")

	records, err := db.GetQuestRecordsForQuest(quest, dynamoClient)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records for quest, got %d", len(records))
	}
	for _, r := range records {
		if r.Quest != quest {
			t.Errorf("expected quest %q, got %q", quest, r.Quest)
		}
	}
}
```

- [ ] **Step 3: Run the test to verify it fails**

```bash
cd /Users/lukebrutvan/Desktop/psostats-client
go test ./server/internal/db/... -run TestGetQuestRecordsForQuest -v
```

Expected: compile error — `db.GetQuestRecordsForQuest undefined`.

- [ ] **Step 4: Implement `GetQuestRecordsForQuest` in `db.go`**

Add after the existing `GetQuestRecords` function in `server/internal/db/db.go`:

```go
func GetQuestRecordsForQuest(quest string, dynamoClient *dynamodb.DynamoDB) ([]model.Game, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("Quest"), expression.Value(quest))).
		Build()
	if err != nil {
		return nil, err
	}
	result, err := dynamoClient.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(QuestRecordsTable),
	})
	if err != nil {
		return nil, err
	}
	games := make([]model.Game, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &games)
	return games, err
}
```

- [ ] **Step 5: Run the test to verify it passes**

```bash
go test ./server/internal/db/... -run TestGetQuestRecordsForQuest -v
```

Expected: `PASS`.

- [ ] **Step 6: Commit**

```bash
git add server/internal/db/db.go server/internal/db/db_test.go
git commit -m "feat: add GetQuestRecordsForQuest DB function"
```

---

## Task 3: DB — `GetAllGamesForPlayerQuest`

**Files:**
- Modify: `server/internal/db/db_test.go`
- Modify: `server/internal/db/db.go`

- [ ] **Step 1: Add table helper and write helper to `db_test.go`**

Add after the `CreateQuestRecordsTable` function:

```go
func CreateRecentGamesByPlayerTable(dynamoClient *dynamodb.DynamoDB) error {
	_, err := dynamoClient.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(db.RecentGamesByPlayerTable),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("Player"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("IdInt"), AttributeType: aws.String("N")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("Player"), KeyType: aws.String(dynamodb.KeyTypeHash)},
			{AttributeName: aws.String("IdInt"), KeyType: aws.String(dynamodb.KeyTypeRange)},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	return err
}

func writeTestPlayerGame(t *testing.T, dynamoClient *dynamodb.DynamoDB, player, questName, gameId string, idInt int) {
	t.Helper()
	_, err := dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(db.RecentGamesByPlayerTable),
		Item: map[string]*dynamodb.AttributeValue{
			"Player":   {S: aws.String(player)},
			"IdInt":    {N: aws.String(strconv.Itoa(idInt))},
			"Id":       {S: aws.String(gameId)},
			"Quest":    {S: aws.String(questName)},
			"Category": {S: aws.String("2n")},
			"Episode":  {N: aws.String("1")},
		},
	})
	if err != nil {
		t.Fatalf("failed to write test player game: %v", err)
	}
}
```

Add `"strconv"` to the import block in `db_test.go` if it is not already present.

Also add to `CreateAllTables` after the `QuestRecordsTable` block just added in Task 2:

```go
	if _, exists := tables[db.RecentGamesByPlayerTable]; !exists {
		if err = CreateRecentGamesByPlayerTable(dynamoClient); err != nil {
			return err
		}
	}
```

- [ ] **Step 2: Write the failing test**

```go
func TestGetAllGamesForPlayerQuest(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		t.Fatal(err)
	}
	dynamoClient := dynamodb.New(sess)
	if err = CreateAllTables(dynamoClient); err != nil {
		t.Fatal(err)
	}

	const player = "testplayer"
	const targetQuest = "Sweep-up Operation #1"

	writeTestPlayerGame(t, dynamoClient, player, targetQuest, "pg-1", 1)
	writeTestPlayerGame(t, dynamoClient, player, targetQuest, "pg-2", 2)
	writeTestPlayerGame(t, dynamoClient, player, "Lost WORKS Machine", "pg-3", 3)

	games, err := db.GetAllGamesForPlayerQuest(player, targetQuest, dynamoClient)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 2 {
		t.Fatalf("expected 2 games for player+quest, got %d", len(games))
	}
	for _, g := range games {
		if g.Quest != targetQuest {
			t.Errorf("expected quest %q, got %q", targetQuest, g.Quest)
		}
	}
}
```

- [ ] **Step 3: Run the test to verify it fails**

```bash
go test ./server/internal/db/... -run TestGetAllGamesForPlayerQuest -v
```

Expected: compile error — `db.GetAllGamesForPlayerQuest undefined`.

- [ ] **Step 4: Implement `GetAllGamesForPlayerQuest` in `db.go`**

Add after `GetQuestRecordsForQuest`:

```go
func GetAllGamesForPlayerQuest(player, quest string, dynamoClient *dynamodb.DynamoDB) ([]model.Game, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key("Player"), expression.Value(player))).
		WithFilter(expression.Equal(expression.Name("Quest"), expression.Value(quest))).
		Build()
	if err != nil {
		return nil, err
	}

	games := make([]model.Game, 0)
	var lastKey map[string]*dynamodb.AttributeValue

	for {
		input := &dynamodb.QueryInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
			FilterExpression:          expr.Filter(),
			ScanIndexForward:          aws.Bool(false),
			TableName:                 aws.String(RecentGamesByPlayerTable),
		}
		if lastKey != nil {
			input.ExclusiveStartKey = lastKey
		}

		result, err := dynamoClient.Query(input)
		if err != nil {
			return nil, err
		}

		page := make([]model.Game, 0)
		if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &page); err != nil {
			return nil, err
		}
		games = append(games, page...)

		lastKey = result.LastEvaluatedKey
		if lastKey == nil {
			break
		}
	}

	sort.Slice(games, func(i, j int) bool {
		return games[i].Timestamp.After(games[j].Timestamp)
	})

	return games, nil
}
```

- [ ] **Step 5: Run the test to verify it passes**

```bash
go test ./server/internal/db/... -run TestGetAllGamesForPlayerQuest -v
```

Expected: `PASS`.

- [ ] **Step 6: Run all DB tests to check nothing regressed**

```bash
go test ./server/internal/db/... -v
```

Expected: all tests `PASS`.

- [ ] **Step 7: Commit**

```bash
git add server/internal/db/db.go server/internal/db/db_test.go
git commit -m "feat: add GetAllGamesForPlayerQuest DB function"
```

---

## Task 4: Server Wiring — FuncMap, Route, Template Field, Stub Handler

**Files:**
- Modify: `server/internal/server/server.go`
- Create: `server/internal/server/display_quest.go`
- Create: `server/internal/templates/questPage.gohtml`

- [ ] **Step 1: Add `sharedTemplateFuncs` and update `ensureParsed` in `server.go`**

`net/url` is already imported. Add the `sharedTemplateFuncs` var immediately before the `ensureParsed` function (around line 163):

```go
var sharedTemplateFuncs = template.FuncMap{
	"urlPathEscape": url.PathEscape,
}

func ensureParsed(templatePath string) *template.Template {
	t, err := template.New("").Funcs(sharedTemplateFuncs).ParseFiles(
		"./server/internal/templates/navbar.gohtml", templatePath)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
```

This replaces the existing `ensureParsed` body. The change from `template.ParseFiles(...)` to `template.New("").Funcs(...).ParseFiles(...)` is safe: all `{{ define "name" }}` blocks are still accessible by name in `ExecuteTemplate` calls — only the template set's root name changes from the first file's basename to `""`.

- [ ] **Step 2: Add `questTemplate` field to the `Server` struct**

In the `Server` struct (around line 26), add after `comboCalcTemplate`:

```go
	questTemplate           *template.Template
```

- [ ] **Step 3: Wire the template and route in `Run()`**

In the `Run()` function, after the existing `s.comboCalcTemplate = ensureParsed(...)` line, add:

```go
	s.questTemplate = ensureParsed("./server/internal/templates/questPage.gohtml")
```

After `s.app.Get("/players/:player", s.PlayerV2Page)`, add:

```go
	s.app.Get("/quest/:quest", s.QuestPage)
```

- [ ] **Step 4: Create stub template `server/internal/templates/questPage.gohtml`**

```html
{{ define "quest" }}
<html lang="en">
<head><title>Quest - PSOStats</title></head>
<body>{{ template "navbar" }}<p>Quest page coming soon</p></body>
</html>
{{ end }}
```

- [ ] **Step 5: Create stub handler `server/internal/server/display_quest.go`**

```go
package server

import (
	"github.com/gofiber/fiber/v2"
	"net/url"
)

func (s *Server) QuestPage(c *fiber.Ctx) error {
	quest, err := url.PathUnescape(c.Params("quest"))
	if err != nil {
		c.Status(500)
		return err
	}
	_ = quest
	err = s.questTemplate.ExecuteTemplate(c.Response().BodyWriter(), "quest", nil)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}
```

- [ ] **Step 6: Build and run the server**

```bash
go build ./server/... && AWS_ACCESS_KEY_ID="" go run ./server/cmd/main.go
```

- [ ] **Step 7: Verify the stub renders**

Open `http://localhost:80/quest/test` in a browser.
Expected: page loads with navbar and "Quest page coming soon".

- [ ] **Step 8: Commit**

```bash
git add server/internal/server/server.go server/internal/server/display_quest.go server/internal/templates/questPage.gohtml
git commit -m "feat: wire QuestPage route and stub template"
```

---

## Task 5: Implement the Full `QuestPage` Handler

**Files:**
- Modify: `server/internal/server/display_quest.go`

- [ ] **Step 1: Replace the stub with the full handler**

Replace the entire contents of `server/internal/server/display_quest.go`:

```go
package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"net/url"
	"time"
)

func (s *Server) QuestPage(c *fiber.Ctx) error {
	quest, err := url.PathUnescape(c.Params("quest"))
	if err != nil {
		c.Status(500)
		return err
	}
	player := c.Query("player")

	questRecords, err := db.GetQuestRecordsForQuest(quest, s.dynamoClient)
	if err != nil {
		return err
	}

	recordsByCategory := make(map[string]model.FormattedGame)
	for _, game := range questRecords {
		fg := getFormattedGame(game)
		fg.Record = true
		recordsByCategory[game.Category] = fg
	}

	var playerGames []model.FormattedGame
	playerNotFound := false

	if player != "" {
		rawGames, err := db.GetAllGamesForPlayerQuest(player, quest, s.dynamoClient)
		if err != nil {
			return err
		}
		if len(rawGames) == 0 {
			playerNotFound = true
		}

		pbIds := make(map[string]string)
		bestTimes := make(map[string]time.Duration)
		for _, g := range rawGames {
			if best, ok := bestTimes[g.Category]; !ok || g.Time < best {
				bestTimes[g.Category] = g.Time
				pbIds[g.Category] = g.Id
			}
		}

		for _, game := range rawGames {
			fg := getFormattedGame(game)
			if pbIds[game.Category] == game.Id {
				fg.Pb = true
			}
			if rec, ok := recordsByCategory[game.Category]; ok && rec.Id == game.Id {
				fg.Record = true
				fg.Pb = false
			}
			playerGames = append(playerGames, fg)
		}
	}

	pageModel := struct {
		QuestName      string
		Records        map[string]model.FormattedGame
		Player         string
		PlayerGames    []model.FormattedGame
		PlayerNotFound bool
	}{
		QuestName:      quest,
		Records:        recordsByCategory,
		Player:         player,
		PlayerGames:    playerGames,
		PlayerNotFound: playerNotFound,
	}

	err = s.questTemplate.ExecuteTemplate(c.Response().BodyWriter(), "quest", pageModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}
```

- [ ] **Step 2: Build to verify it compiles**

```bash
go build ./server/...
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add server/internal/server/display_quest.go
git commit -m "feat: implement QuestPage handler"
```

---

## Task 6: Build the Full Quest Page Template

**Files:**
- Modify: `server/internal/templates/questPage.gohtml`

- [ ] **Step 1: Replace the stub with the full template**

Replace the entire contents of `server/internal/templates/questPage.gohtml`:

```html
{{ define "quest" }}
<html lang="en">
<head>
    <meta name="viewport" content="width=device-width">
    <title>{{ .QuestName }} - PSOStats</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-+0n0xVW2eSR5OomGNYDnhzAbDsOXxcvSN1TPprVMTNDbiYZCxYbOOl7+AMvyTG2x" crossorigin="anonymous">
    <link href="/static/main2.css" rel="stylesheet" type="text/css">
</head>
<body>
<div class="container">
    {{ template "navbar" }}
    <div class="row">
        <div class="col">
            <h1>{{ .QuestName }}</h1>
        </div>
    </div>

    {{ if .Records }}
    <div class="row">
        <div class="col">
            <h2>Records</h2>
        </div>
    </div>
    {{ range $category, $game := .Records }}
    <div class="row quest-row">
        <div class="col-3 col-xl-2">
            <span class="quest-category">{{ $game.NumPlayers }}P
                {{ if $game.PbRun }}
                    <img src="/static/twins_cropped.png" height="24px" width="24px" alt="PB" title="PB" style="margin-bottom: 4px"/>
                {{ else }}
                    <img src="/static/shifta_cropped.png" height="24px" width="24px" alt="No-PB" title="No-PB" style="margin-bottom: 4px"/>
                {{ end }}
            </span>
        </div>
        <div class="col-3">
            <a href="/game/{{ $game.Id }}" class="quest-time">{{ $game.Time }} 🥇</a>
        </div>
        <div class="col-6">
            {{ range $index, $player := $game.Players }}
                {{ if gt (len .Name) 0 }}
                    <div>
                        {{ if $player.HasPov }}
                            <a href="/game/{{ $game.Id }}/{{ $index }}"><span style="width:85px; display: inline-block">{{ $player.Class }}</span>{{ $player.Name }}</a>
                        {{ else }}
                            <span style="width:85px; display: inline-block">{{ $player.Class }}</span>{{ $player.Name }}
                        {{ end }}
                    </div>
                {{ end }}
            {{ end }}
        </div>
    </div>
    {{ end }}
    {{ end }}

    <div class="row" style="margin-top: 24px; margin-bottom: 8px;">
        <div class="col-12 col-md-6">
            <h2>Player History</h2>
            <form method="GET">
                <div class="input-group">
                    <input type="text" class="form-control" name="player" value="{{ .Player }}" placeholder="Player name">
                    <button class="btn btn-primary" type="submit">Search</button>
                </div>
            </form>
        </div>
    </div>

    {{ if .Player }}
        {{ if .PlayerNotFound }}
        <div class="row">
            <div class="col">
                <p class="text-muted">No runs found for <strong>{{ .Player }}</strong> on this quest.</p>
            </div>
        </div>
        {{ else }}
        <div class="row" style="margin-top: 8px;">
            <div class="col">
                <h3><a href="/players/{{ .Player }}">{{ .Player }}</a></h3>
            </div>
        </div>
        {{ range .PlayerGames }}
            {{ $game := . }}
            <div class="row quest-row">
                <div class="col-8 col-md-3">
                    <h6 class="text-muted" title="{{ .Date }}">{{ .RelativeDate }}</h6>
                </div>
                <div class="col-4 col-md-1">
                    <span class="quest-category">{{ .NumPlayers }}P</span>
                </div>
                <div class="col-4 col-md-2">
                    <span class="quest-time">{{ .Time }}<small>{{ if $game.Record }} 🥇 {{ else if $game.Pb }} PB {{ end }}</small></span>
                </div>
                <div class="col-8 col-md-6">
                    {{ range $index, $player := .Players }}
                        {{ if gt (len .Name) 0 }}
                            <div>
                                {{ if $player.HasPov }}
                                    <a href="/game/{{ $game.Id }}/{{ $index }}"><span style="width:85px; display: inline-block">{{ $player.Class }}</span>{{ $player.Name }}</a>
                                {{ else }}
                                    <span style="width:85px; display: inline-block">{{ $player.Class }}</span>{{ $player.Name }}
                                {{ end }}
                            </div>
                        {{ end }}
                    {{ end }}
                </div>
            </div>
        {{ end }}
        {{ end }}
    {{ end }}
</div>
</body>
</html>
{{ end }}
```

- [ ] **Step 2: Run the server and test the records section**

```bash
AWS_ACCESS_KEY_ID="" go run ./server/cmd/main.go
```

Visit `http://localhost:80/quest/Sweep-up%20Operation%20%231` (or any quest that has records in the live DB). The server defaults to `http://localhost:8000` (DynamoDB Local) when `AWS_ACCESS_KEY_ID` is unset — to test against the real DB, set `AWS_ACCESS_KEY_ID` to your key.

Expected: quest name heading, records section if records exist, player search form.

- [ ] **Step 3: Test the player search**

Append `?player=<a known player name>` to the URL.
Expected: player name subheading linking to `/players/{name}`, list of runs with time, date, party members, PB/record badges.

- [ ] **Step 4: Test the not-found case**

Append `?player=zzzznonexistent` to the URL.
Expected: "No runs found for zzzznonexistent on this quest."

- [ ] **Step 5: Commit**

```bash
git add server/internal/templates/questPage.gohtml
git commit -m "feat: add questPage template"
```

---

## Task 7: Records Page — Quest Name Links

**Files:**
- Modify: `server/internal/templates/recordsV2.gohtml`

- [ ] **Step 1: Wrap the quest name heading in a link**

In `server/internal/templates/recordsV2.gohtml`, find line 45:

```html
                    <h5>{{ $quest }}</h5>
```

Replace with:

```html
                    <h5><a href="/quest/{{ $quest | urlPathEscape }}">{{ $quest }}</a></h5>
```

`urlPathEscape` is the custom function registered in `sharedTemplateFuncs` in Task 4. It calls `url.PathEscape`, encoding spaces as `%20` and colons as `%3A` so quest names like `"Maximum Attack E: Desert"` become valid URL path segments.

- [ ] **Step 2: Run the server and verify records page links**

Visit `http://localhost:80/records`.
Expected: every quest name is now a clickable link.

- [ ] **Step 3: Click a quest link and verify the quest page loads**

Expected: quest page renders with the quest name heading and records section.

- [ ] **Step 4: Commit**

```bash
git add server/internal/templates/recordsV2.gohtml
git commit -m "feat: link quest names on records page to quest pages"
```

---

## Self-Review Checklist (completed inline)

- **Spec coverage:**
  - ✅ Records page quest links → Task 7
  - ✅ New `/quest/:quest` route → Task 4
  - ✅ `GetQuestRecordsForQuest` DB function → Task 2
  - ✅ `GetAllGamesForPlayerQuest` DB function → Task 3
  - ✅ Quest page: records section + player search + player results → Tasks 5, 6
  - ✅ `PlayerNotFound` case → Task 5 (handler), Task 6 (template)
  - ✅ FuncMap for `urlPathEscape` → Task 4

- **Placeholder scan:** None found. All steps include exact code.

- **Type consistency:**
  - `GetQuestRecordsForQuest` → returns `[]model.Game` → handler calls `getFormattedGame` on each → `map[string]model.FormattedGame` keyed by `game.Category` → template uses `{{ range $category, $game := .Records }}` ✅
  - `GetAllGamesForPlayerQuest` → returns `[]model.Game` → handler builds `[]model.FormattedGame` → template uses `{{ range .PlayerGames }}` ✅
  - Template fields match handler's anonymous struct exactly ✅
  - `urlPathEscape` registered in `sharedTemplateFuncs`, used in `recordsV2.gohtml` ✅
