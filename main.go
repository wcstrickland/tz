package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./tz.db")
	if err != nil {
		log.Fatalln("error connecting to db", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln("error pinging db", err)
	}

	app := fiber.New()
	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("/tz/zip:zip\treturns a timezone from a given zipcode\nEXAMPLE: /tz/zip/35210\n/tz/area/:area\treturns a timezone from a given area code\n/tz/state/:state\treturns a timezone from a given state\nEXAMPLE: /tz/state/AL\n\n/download\tdownloads the underlying sqlite database file")
	})

	app.Get("/tz/zip/:zip", func(c *fiber.Ctx) error {
		zip := c.Params("zip")
		var tz string
		row := db.QueryRow("select timezone from tz where zip = ?", zip)
		if err := row.Scan(&tz); err != nil {
			return c.Status(404).SendString("No Dice")
		}
		return c.SendString(tz)
	})

	app.Get("/tz/area/:area", func(c *fiber.Ctx) error {
		area := c.Params("area")
		var tz string
		row := db.QueryRow("select timezone from tz where area_codes = ?", area)
		if err := row.Scan(&tz); err != nil {
			return c.Status(404).SendString("No Dice")
		}
		return c.SendString(tz)
	})

	app.Get("/tz/state/:state", func(c *fiber.Ctx) error {
		state := c.Params("state")
		var tz string
		row := db.QueryRow("select timezone from tz where state = ?", state)
		if err := row.Scan(&tz); err != nil {
			return c.Status(404).SendString("No Dice")
		}
		return c.SendString(tz)
	})

	app.Get("/download", func(c *fiber.Ctx) error {
		return c.Download("./tz.db")
	})
	port := os.Getenv("PORT")
	if os.Getenv("PORT") == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
