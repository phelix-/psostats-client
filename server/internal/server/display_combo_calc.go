package server

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/psoclasses"
	"github.com/phelix-/psostats/v2/server/internal/enemies"
	"github.com/phelix-/psostats/v2/server/internal/weapons"
)

func (s *Server) ComboCalcOpmPage(c *fiber.Ctx) error {
	return s.comboCalcPage(
		true,
		psoclasses.GetAll(),
		enemies.GetEnemiesUltOpm(),
		weapons.GetWeapons(),
		weapons.GetFrames(),
		c,
	)
}

func (s *Server) ComboCalcMultiPage(c *fiber.Ctx) error {
	return s.comboCalcPage(
		false,
		psoclasses.GetAll(),
		enemies.GetEnemiesUltMulti(),
		weapons.GetWeapons(),
		weapons.GetFrames(),
		c,
	)
}

func (s *Server) ComboCalcUltima(c *fiber.Ctx) error {
	return s.comboCalcPage(
		false,
		psoclasses.GetAllUltima(),
		enemies.GetEnemiesUltima(),
		weapons.GetWeaponsUltima(),
		weapons.GetFramesUltima(),
		c,
	)
}

func (s *Server) comboCalcPage(
	opm bool,
	allClasses []psoclasses.PsoClass,
	allEnemies []enemies.Enemy,
	allWeapons []weapons.Weapon,
	allFrames []weapons.Frame,
	c *fiber.Ctx,
) error {
	sortedEnemies := make(map[string][]string)
	for _, enemy := range allEnemies {
		enemiesInArea := sortedEnemies[enemy.Location]
		if enemiesInArea == nil {
			enemiesInArea = make([]string, 0)
		}
		enemiesInArea = append(enemiesInArea, enemy.Name)
		sortedEnemies[enemy.Location] = enemiesInArea
	}

	weaponsJson := toJsonMap(allWeapons, func(weapon weapons.Weapon) string { return weapon.Name })
	frameJson := toJsonMap(allFrames, func(frame weapons.Frame) string { return frame.Name })
	psoClassJson := toJsonMap(allClasses, func(class psoclasses.PsoClass) string { return class.Name })
	enemyJson := toJsonMap(allEnemies, func(enemy enemies.Enemy) string { return enemy.Name })
	enemyNameJson, _ := json.Marshal(sortedEnemies)
	infoModel := struct {
		Opm            bool
		Classes        []psoclasses.PsoClass
		ClassStatsJson string
		EnemyNameSort  string
		EnemiesJson    string
		Frames         []weapons.Frame
		FramesJson     string
		Weapons        []weapons.Weapon
		WeaponsJson    string
	}{
		Opm:            opm,
		Classes:        allClasses,
		ClassStatsJson: psoClassJson,
		EnemyNameSort:  string(enemyNameJson),
		EnemiesJson:    enemyJson,
		Frames:         allFrames,
		FramesJson:     frameJson,
		Weapons:        allWeapons,
		WeaponsJson:    weaponsJson,
	}
	//s.comboCalcTemplate = ensureParsed("./server/internal/templates/comboCalc.gohtml")
	err := s.comboCalcTemplate.ExecuteTemplate(c.Response().BodyWriter(), "combo-calc", infoModel)

	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func toJsonMap[V any](items []V, keyMapper func(V) string) string {
	itemMap := make(map[string]V)
	for _, item := range items {
		name := keyMapper(item)
		itemMap[name] = item
	}
	jsonBytes, _ := json.Marshal(itemMap)
	return string(jsonBytes)
}
