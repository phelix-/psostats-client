package server

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
	"io/ioutil"
)

//func (s *Server) GetGeometry(c *fiber.Ctx) error {
//	gameId := "26873"
//	gem := 1
//	if dataFrames, err := db.GetDataFrames(gameId, gem, s.dynamoClient); err == nil {
//		floorMeshes := GetFloorMeshes(dataFrames[0].Map, dataFrames[0].MapVariation)
//		floorMeshes.DataFrames = dataFrames
//		jsonBytes, err := json.Marshal(floorMeshes)
//		if err != nil {
//			return err
//		}
//		c.Response().AppendBody(jsonBytes)
//		c.Response().Header.Set("Content-Type", "application/json")
//	}
//	return nil
//}

func (s *Server) GetGameJs(c *fiber.Ctx) error {
	content, err := ioutil.ReadFile("./js/game.js")
	if err != nil {
		return err
	}
	_, err = c.Response().BodyWriter().Write(content)
	c.Response().Header.Set("Content-Type", "application/javascript")
	return err
}

func (s *Server) OrbitControlsJs(c *fiber.Ctx) error {
	content, err := ioutil.ReadFile("./js/OrbitControls.js")
	if err != nil {
		return err
	}
	_, err = c.Response().BodyWriter().Write(content)
	c.Response().Header.Set("Content-Type", "application/javascript")
	return err
}

func (s *Server) DraughtsJs(c *fiber.Ctx) error {
	content, err := ioutil.ReadFile("./js/draughts.js")
	if err != nil {
		return err
	}
	_, err = c.Response().BodyWriter().Write(content)
	c.Response().Header.Set("Content-Type", "application/javascript")
	return err
}

func (s *Server) ThreeJs(c *fiber.Ctx) error {
	content, err := ioutil.ReadFile("./js/three.module.js")
	if err != nil {
		return err
	}
	_, err = c.Response().BodyWriter().Write(content)
	c.Response().Header.Set("Content-Type", "application/javascript")
	return err
}

func GetFloorMeshes(mapNum uint16, mapVariation uint16) *FloorMeshes {
	floorName := fmt.Sprintf("Unknown Map %v", mapNum)
	switch mapNum {
	case 0:
		floorName = "city00"
	case 1:
		floorName = "forest01"
	case 2:
		floorName = "forest02"
	case 3:
		floorName = "cave01"
	case 4:
		floorName = "cave02"
	case 5:
		floorName = "cave03"
	case 6:
		floorName = "machine01"
	case 7:
		floorName = "machine02"
	case 8:
		floorName = "ancient01"
	case 9:
		floorName = "ancient02"
	case 10:
		floorName = "ancient03"
	case 11:
		floorName = "boss01"
	case 12:
		floorName = "boss02"
	case 13:
		floorName = "boss03"
	case 14:
		floorName = "darkfalz00"
	case 15:
		floorName = "lobby"
	case 16:
		floorName = "vs01"
	case 17:
		floorName = "vs02"
	case 18:
		floorName = "labo00"
	case 19:
		floorName = "ruins01"
	case 20:
		floorName = "ruins02"
	case 21:
		floorName = "space01"
	case 22:
		floorName = "space02"
	case 23:
		floorName = "jungle01"
	case 24:
		floorName = "jungle02"
	case 25:
		floorName = "jungle03"
	case 26:
		floorName = "jungle04"
	case 27:
		floorName = "jungle05"
	case 28:
		floorName = "seabed01"
	case 29:
		floorName = "seabed02"
	case 30:
		floorName = "boss05"
	case 31:
		floorName = "boss06"
	case 32:
		floorName = "boss07"
	case 33:
		floorName = "boss08"
	case 34:
		floorName = "jungle06"
	case 35:
		floorName = "jungle07"
	case 36:
		floorName = "wilds01"
	case 37:
		floorName = "wilds01"
	case 38:
		floorName = "wilds01"
	case 39:
		floorName = "wilds01"
	case 40:
		floorName = "crater01"
	case 41:
		floorName = "desert01"
	case 42:
		floorName = "desert02"
	case 43:
		floorName = "desert03"
	case 44:
		floorName = "boss09"
	case 45:
		floorName = "city02"
	}

	floorMeshes := FloorMeshes{}
	jsonBytes, _ := ioutil.ReadFile(fmt.Sprintf("js/map_%v_%02dc.json", floorName, mapVariation))
	_ = json.Unmarshal(jsonBytes, &floorMeshes)

	return &floorMeshes
}

type FloorMeshes struct {
	Meshes     [][]Coordinate    `json:"meshes"`
	Normals    [][]Coordinate    `json:"normals"`
	DataFrames []model.DataFrame `json:"dataFrames"`
}

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
