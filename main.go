package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vidhanio/dmux"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Loading .env...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	log.Info().Msg(".env loaded.")

	log.Info().Msg("Initializing Discord bot...")

	mux := dmux.NewMux(os.Getenv("DISCORD_TOKEN"), os.Getenv("DISCORD_GUILD_ID"))

	mux.HandleFunc("/gizmo resource:integer", func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		url := fmt.Sprintf("http://localhost:8000/gizmos/%d", dmux.Option(i, "resource").IntValue())

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error().Err(err).Msg("Error creating request")

			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("Error sending request")

			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error().Err(err).Msg("Error getting response")

			return
		}

		gizmoResp := &GizmoResponse{}

		err = json.NewDecoder(resp.Body).Decode(gizmoResp)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding response")

			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					gizmoResp.Gizmo.Embed(),
				},
				Flags: 64,
			},
		},
		)
	})

	log.Info().Msg("Discord bot initialized.")

	log.Info().Msg("Starting Discord bot...")

	go func() {
		err = mux.Serve()
		if err != nil {
			log.Fatal().Err(err).Msg("Error serving")
		}
	}()

	log.Info().Msg("Discord bot started.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc

	log.Info().Msg("Stopping Discord bot...")

	err = mux.Close()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to stop Discord bot")
	}

	log.Info().Msg("Discord bot stopped.")
}
