package main

import (
	"database/sql"
	"log"

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

	app.Get("/tz/:zip", func(c *fiber.Ctx) error {
		zip := c.Params("zip")
		var tz string
		row := db.QueryRow("select timezone from tz where zip = ?", zip)
		if err := row.Scan(&tz); err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString(tz)
	})

	log.Fatal(app.Listen(":3000"))
}
