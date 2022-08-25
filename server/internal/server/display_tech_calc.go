package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) TechCalcPage(c *fiber.Ctx) error {
	model := TechCalcModel{}

	techCalcTemplate := ensureParsed("./server/internal/templates/techCalc.gohtml")
	err := techCalcTemplate.ExecuteTemplate(c.Response().BodyWriter(), "tech-calc", model)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

type TechCalcModel struct {
}
