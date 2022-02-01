package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Gizmo struct {
	Title       string   `json:"title"`
	Materials   string   `json:"materials"`
	Description string   `json:"description"`
	Resource    int      `json:"resource"`
	Answers     []string `json:"answers"`
}

type GizmoResponse struct {
	Message string `json:"message"`
	Gizmo   *Gizmo `json:"gizmo,omitempty"`
}

type GizmosResponse struct {
	Message string   `json:"message"`
	Gizmos  []*Gizmo `json:"gizmos,omitempty"`
}

func (g Gizmo) Embed() *discordgo.MessageEmbed {
	fields := make([]*discordgo.MessageEmbedField, len(g.Answers))
	for i, answer := range g.Answers {
		fields[i] = &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("Question %d", i+1),
			Value: answer,
		}
	}

	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: "https://i.imgur.com/eB8Z1NZ.png",
			Name:    "ExploreLearning Gizmos",
		},
		Title:       g.Title,
		URL:         fmt.Sprintf("https://gizmos.explorelearning.com/index.cfm?method=cResource.dspDetail&resourceID=%d", g.Resource),
		Description: g.Description,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://el-gizmos.s3.amazonaws.com/img/GizmoSnap/%dDET.jpg", g.Resource),
		},
		Fields: fields,
	}
}
