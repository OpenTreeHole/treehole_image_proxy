package main

import "github.com/gofiber/fiber/v2"

func GetAuthToken(c *fiber.Ctx) error {
	return c.SendString("PF.obj.config.auth_token = \"123456789\"")
}

func UploadImage(c *fiber.Ctx) error {
	var response CheveretoUploadResponse
	file, err := c.FormFile("source")
	if err != nil {
		return err
	}
	err = ProxyUploadImage(file, &response)
	if err != nil {
		return err
	}
	return c.JSON(&response)
}
