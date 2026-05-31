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
