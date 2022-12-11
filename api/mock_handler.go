package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
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
		fmt.Println("ERROR", errEncodingBody)
	}

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
