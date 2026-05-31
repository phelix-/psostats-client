# Quest Pages Design

**Date:** 2026-05-31
**Status:** Approved

## Problem

The player profile page (`/players/{name}`) caps the game history at 15 recent games and only shows PBs. There is no way for a player to find their complete run history for a specific quest. The records page lists top times per quest but provides no drill-down or player-specific view.

## Goal

Add per-quest dedicated pages reachable from the records page. Each quest page allows a visitor to search by player name and see that player's complete recorded run history for the quest. The player profile improvement is deferred — the quest pages are the correct home for this capability.

## Out of Scope

- Browsing all runs for a quest across all players (no `games_by_quest` table exists; adding one is a separate decision)
- Pagination UI on the quest page (full result set rendered server-side; add later if needed)
- Changes to the player profile page (deferred)

## Changes

### 1. Records page — quest name links

**File:** `server/internal/templates/recordsV2.gohtml`

Replace the plain quest heading with a link:

```html
<h5><a href="/quest/{{ $quest | urlquery }}">{{ $quest }}</a></h5>
```

`urlquery` percent-encodes the quest name for safe URL embedding. No handler changes.

### 2. New route

**File:** `server/internal/server/server.go`

Register alongside existing routes:

```
GET /quest/:quest  →  s.QuestPage
```

Wire a `questTemplate` field on the `Server` struct and parse `questPage.gohtml` at startup using the existing `ensureParsed` pattern.

### 3. New DB functions

**File:** `server/internal/db/db.go`

#### `GetQuestRecordsForQuest(quest string, client *dynamodb.DynamoDB) ([]model.Game, error)`

Queries the `quest_records` table with `KeyCondition: Quest = quest`. Returns all category records (1P/2P/3P/4P) for the quest. Avoids the existing full-table scan used on the records page.

#### `GetAllGamesForPlayerQuest(player, quest string, client *dynamodb.DynamoDB) ([]model.Game, error)`

Queries `recent_games_by_player_2` with:
- `KeyCondition: Player = player`
- `FilterExpression: Quest = quest`
- `ScanIndexForward: false`
- No `Limit` — paginates through all DynamoDB pages until `LastEvaluatedKey` is nil

Returns the complete run history for that player+quest combination, sorted newest-first. Reads all of the player's game records on each call (DynamoDB reads before filtering); acceptable for current data volumes with no new tables required.

### 4. New handler

**File:** `server/internal/server/display_quest.go`

```
func (s *Server) QuestPage(c *fiber.Ctx) error
```

1. Decode `c.Params("quest")` via `url.PathUnescape`
2. Call `db.GetQuestRecordsForQuest` — always
3. Read optional `?player=` query param
4. If player present:
   - Call `db.GetAllGamesForPlayerQuest(player, quest, ...)`
   - Identify player's PB run (shortest time in result set) and flag it
   - Mark any run whose ID matches the site-wide record
5. Build template model and render `questPage.gohtml`

Template model fields:
- `QuestName string`
- `Records map[string]model.FormattedGame` — keyed by category string (e.g. `"1p"`, `"2n"`). Handler calls `getFormattedGame` on each item from `GetQuestRecordsForQuest` and keys by `game.Category`, consistent with how `sortGames` works elsewhere.
- `Player string` — empty if no search
- `PlayerGames []model.FormattedGame` — empty if no search
- `PlayerNotFound bool` — true if player param given but zero games returned

### 5. New template

**File:** `server/internal/templates/questPage.gohtml`

Layout (top to bottom):
1. Navbar
2. `<h1>` quest name
3. **Records section** — compact 1P/2P/3P/4P top times, same row style as existing pages. Secondary context.
4. **Player search form** — `<form method="GET">`, text input `name="player"`, submit button. Pre-filled with current player value if one was searched. Primary UI element.
5. **Player results** (rendered only when `Player` is non-empty):
   - Subheading with player name linking to `/players/{name}`
   - If `PlayerNotFound`: a "No runs found" message
   - Otherwise: run table — columns: date (relative), party size, time, PB/Record badge, party members with POV links. Matches the style of the recent games section on the player profile page.

## DynamoDB Access Pattern Summary

| Query | Table | Key Used | Notes |
|---|---|---|---|
| Quest top records | `quest_records` | PK=Quest | Targeted query, no scan |
| Player run history for quest | `recent_games_by_player_2` | PK=Player + FilterExpr on Quest | Reads all player games server-side |

## Deferred

- **Player profile quest history** — shelved in favour of quest pages. Revisit post-launch.
- **`games_by_quest` table** — if "browse all runs for a quest without specifying a player" becomes a requirement, a new table with PK=Quest is the right solution. Requires write-side changes to `PostGame` and a historical backfill.
