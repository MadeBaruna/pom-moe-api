package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/MadeBaruna/pom-moe/internal/queue"
	"github.com/MadeBaruna/pom-moe/internal/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
	"golang.org/x/net/publicsuffix"
)

type GetWarpBody struct {
	Url string `json:"url"`
}

var allowedHosts = []string{
	"mihoyo.com",
	"hoyoverse.com",
}

func GetWarp(c *fiber.Ctx) error {
	_, proxy := utils.GetProxy(c.IP())

	body := new(GetWarpBody)
	if err := c.BodyParser(body); err != nil {
		return c.SendStatus(400)
	}

	u, err := url.Parse(body.Url)
	if err != nil {
		return c.SendStatus(400)
	}

	domain, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname())
	if err != nil {
		return c.SendStatus(400)
	}
	if !slices.Contains(allowedHosts, domain) {
		return c.SendStatus(400)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	res, err := client.Get(u.String())
	if err != nil {
		log.Println(err)
		if res != nil {
			return c.SendStatus(res.StatusCode)
		}

		return c.SendStatus(400)
	}

	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)
	queue.Publish("store", data)

	c.Response().Header.SetContentType(res.Header.Get("Content-Type"))
	return c.Status(res.StatusCode).Send(data)
}
