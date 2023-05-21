package models

import (
	"encoding/json"
	"time"
)

type GachaItem struct {
	Uid       string `json:"uid"`
	GachaID   string `json:"gacha_id"`
	GachaType string `json:"gacha_type"`
	ItemID    string `json:"item_id"`
	Count     string `json:"count"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	ItemType  string `json:"item_type"`
	RankType  string `json:"rank_type"`
	ID        string `json:"id"`
}

type GachaResponse struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		Page           string      `json:"page"`
		Size           string      `json:"size"`
		List           []GachaItem `json:"list"`
		Region         string      `json:"region"`
		RegionTimeZone int         `json:"region_time_zone"`
	} `json:"data"`
}

func (g *GachaResponse) Marshal() ([]byte, error) {
	return json.Marshal(g)
}

type ItemType string

const (
	Character ItemType = "character"
	LightCone ItemType = "lightcone"
)

var ItemTypes = map[string]ItemType{
	"Character":  Character,
	"Light Cone": LightCone,
}

type BannerType string

const (
	StandardBanner  BannerType = "standard"
	BeginnerBanner  BannerType = "beginner"
	CharacterBanner BannerType = "character"
	LightconeBanner BannerType = "lightcone"
)

var BannerTypes = map[string]BannerType{
	"Standard":   StandardBanner,
	"Beginner":   BeginnerBanner,
	"Character":  CharacterBanner,
	"Light Cone": LightconeBanner,
}

type Warp struct {
	ID         string
	UID        string
	GachaID    string
	GachaType  string
	ItemID     string
	Count      int8
	Name       string
	TimeRaw    string
	Time       time.Time
	ItemType   ItemType
	Rarity     int8
	BannerType BannerType
	Region     string
}
