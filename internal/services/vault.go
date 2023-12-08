package services

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MadeBaruna/pom-moe/internal/db"
	"github.com/MadeBaruna/pom-moe/internal/models"
	"github.com/MadeBaruna/pom-moe/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var bannerPrefixs = map[byte]models.BannerType{
	'4': models.BeginnerBanner,
	'1': models.StandardBanner,
	'2': models.CharacterBanner,
	'3': models.LightconeBanner,
}

var query = `
insert into warps(id, uid, gacha_id, gacha_type, count, item_id, item_type, name, rarity, time_raw, time, banner_type, region) 
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
on conflict (id) do nothing`

func StoreWarp(data *models.GachaResponse) bool {
	len := len(data.Data.List)

	offset := time.FixedZone(data.Data.Region, data.Data.RegionTimeZone*3600)

	batch := &pgx.Batch{}
	for i := 0; i < len; i++ {
		w := data.Data.List[i]

		count, err := strconv.ParseInt(w.Count, 10, 8)
		if err != nil {
			log.Error().Err(err).Msg("Error convert count")
			return false
		}

		rarity, err := strconv.ParseInt(w.RankType, 10, 8)
		if err != nil {
			log.Error().Err(err).Msg("Error convert rarity")
			return false
		}

		t, err := time.ParseInLocation("2006-01-02 15:04:05", w.Time, offset)
		if err != nil {
			log.Error().Err(err).Msg("Error convert time")
			return false
		}

		batch.Queue(
			query,
			w.ID,
			w.Uid,
			w.GachaID,
			w.GachaType,
			count,
			w.ItemID,
			models.ItemTypes[w.ItemType],
			utils.Slugify(w.Name),
			rarity,
			w.Time,
			t,
			bannerPrefixs[w.GachaID[0]],
			data.Data.Region,
		)
	}

	br := db.Pool().SendBatch(context.Background(), batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		jsonStr, errJson := json.Marshal(data)
		if errJson != nil {
			log.Error().Err(err).Str("data", string(jsonStr)).Msg("Error insert warps")
		} else {
			log.Error().Err(err).Msg("Error insert warps")
		}
		return false
	}

	return true
}
