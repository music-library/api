package api

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/global"
	version "gitlab.com/music-library/music-api/version"
)

type MockResponse struct {
	Body        map[string]interface{}
	BodyString  string
	BodyLength  int
	Host        string
	Method      string
	Params      string
	Path        string
	QueryString string
	Timestamp   int64
	Version     string
}

func MockHandler(c *fiber.Ctx) error {
	// Encode unknown body data
	var parsedBody map[string]interface{}
	errEncodingBody := c.BodyParser(parsedBody)

	if errEncodingBody != nil {
		log.Error("ERROR", errEncodingBody)
	}

	log.Info("MockHandler", global.DATA_DIR, c.Method(), c.Path(), c.Params("*"), c.Body(), parsedBody)

	res := MockResponse{
		Body:        parsedBody,
		BodyString:  string(c.Body()),
		BodyLength:  len(c.Body()),
		Host:        c.Hostname(),
		Method:      c.Method(),
		Params:      c.Params("*"),
		Path:        c.Path(),
		QueryString: string(c.Request().URI().QueryString()),
		Timestamp:   time.Now().UnixNano() / int64(time.Millisecond),
		Version:     version.GetVersion().FullVersionNumber(false),
	}

	// time.Sleep(20 * time.Second)

	// Send response as JSON
	c.Response().Header.Add("Content-Type", "application/json")
	return c.JSON(res)
}

func Example(c *fiber.Ctx) error {
	// output, _ := sonic.Marshal(res.Response().Body)
	// fmt.Println(string(res.Bytes()[:]))

	msg := fmt.Sprintf("✋ %s", c.Params("*"))
	return c.SendString(msg) // => ✋ register
}
