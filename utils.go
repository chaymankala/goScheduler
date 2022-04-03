package main

import (
	"github.com/gofiber/fiber/v2"
)

func send400Response(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"err": "Invalid data, try again",
	})
}

func send500Response(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"err": "Internal Server Error, try again",
	})
}

func makeParams(message string, jobName string) string {
	output := ""
	output += "message=" + message + "&"
	output += "jobName=" + jobName + "&"
	return output
}
