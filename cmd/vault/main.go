package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/MadeBaruna/pom-moe/internal/db"
	"github.com/MadeBaruna/pom-moe/internal/models"
	"github.com/MadeBaruna/pom-moe/internal/queue"
	"github.com/MadeBaruna/pom-moe/internal/services"
	"github.com/MadeBaruna/pom-moe/internal/utils"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

func main() {
	utils.InitLogger()
	log.Info().Msg("warp vault started")

	utils.LoadEnv()
	queue.LoadNats()
	defer queue.Drain()

	db.Connect()
	defer db.Close()

	sig, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	sub, err := queue.GetStream().PullSubscribe("store", "warp-vault", nats.BindStream("WARP"), nats.MaxDeliver(2))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to subscribe to stream")
	}

l:
	for {
		select {
		case <-sig.Done():
			break l
		default:
		}

		ctx, cancel := context.WithTimeout(sig, 30*time.Second)
		defer cancel()

		msgs, err := sub.Fetch(1, nats.Context(ctx))
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Error().Err(err).Msg("failed to fetch message")
		}

		for _, msg := range msgs {
			var data models.GachaResponse
			err := json.Unmarshal(msg.Data, &data)
			if err != nil {
				log.Error().Err(err).Msg("failed to parse message")
				msg.Ack()
				continue
			}

			if data.Retcode != 0 || len(data.Data.List) == 0 {
				msg.Ack()
				continue
			}

			pending, _ := sub.QueuedMsgs()

			log.Info().Str("uid", data.Data.List[0].Uid).Int("count", len(data.Data.List)).Int("pending", pending).Msg("storing warp data")
			success := services.StoreWarp(&data)
			if success {
				msg.Ack()
				continue
			}

			msg.NakWithDelay(60 * time.Second)
		}
	}
}
