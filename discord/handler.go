package discord

import (
	"encoding/base64"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"path/filepath"
	"strings"
)

var (
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionAdministrator

	Commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "list-emotes",
			Description:              "List all guild emotes",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
		},
		{
			Name:                     "create-emotes",
			Description:              "Create Porygon emotes",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
		},
		{
			Name:                     "delete-emotes",
			Description:              "Delete all emotes (created by Porygon)",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
		},
		{
			Name:        "report",
			Description: "Post an end-of-day report now",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"list-emotes":   listEmotes,
		"create-emotes": createEmotes,
		"delete-emotes": deleteEmotes,
	}
)

func listEmotes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var emotesList strings.Builder

	guildEmotes, _ := s.GuildEmojis(i.GuildID)

	if len(guildEmotes) > 0 {
		emotesList.WriteString("```")
		for _, emote := range guildEmotes {
			emotesList.WriteString(fmt.Sprintf("<:%s:%s>\n", emote.Name, emote.ID))
		}
		emotesList.WriteString("```")
	} else {
		emotesList.WriteString("No guild emotes.")
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: emotesList.String(),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func deleteEmotes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var output strings.Builder

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	guildEmotes, _ := s.GuildEmojis(i.GuildID)

	if len(guildEmotes) > 0 {
		output.WriteString("```")
		for _, emote := range guildEmotes {
			if emote.User.ID == s.State.User.ID {
				err := s.GuildEmojiDelete(i.GuildID, emote.ID)
				if err != nil {
					output.WriteString(fmt.Sprintf("%s - failed to remove: %s\n", emote.Name, err))
				} else {
					output.WriteString(fmt.Sprintf("%s - removed\n", emote.Name))
				}
			} else {
				output.WriteString(fmt.Sprintf("%s - skipping, other owner %s\n", emote.Name, emote.User.String()))
			}
		}
		output.WriteString("```")
	} else {
		output.WriteString("No guild emotes to delete.")
	}

	_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: output.String(),
	})
}

func createEmotes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var output strings.Builder

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	emotesDir := "emojis/override"
	files, err := filepath.Glob(filepath.Join(emotesDir, "*.png"))

	if len(files) > 0 {
		output.WriteString("Using `emojis/override` directory as emotes source.\n")
	} else {
		output.WriteString("Using `emojis` directory as emotes source.\n")
		emotesDir = "emojis"
		files, err = filepath.Glob(filepath.Join(emotesDir, "*.png"))
	}

	if err != nil {
		output.WriteString(fmt.Sprintf("Error reading emotes directory: %s\n", err))
	} else {
		output.WriteString("```")

		// fetch existing emotes
		guildEmotes, _ := s.GuildEmojis(i.GuildID)
		existingEmotes := make(map[string]bool)
		for _, emote := range guildEmotes {
			existingEmotes[emote.Name] = true
		}

		if len(files) == 0 {
			output.WriteString("no files to upload")
		}

		// check and upload every emote we have under emotesDir
		for _, file := range files {
			emoteName := strings.TrimSuffix(filepath.Base(file), ".png")

			if _, exists := existingEmotes[emoteName]; exists {
				output.WriteString(fmt.Sprintf("%s - already there\n", emoteName))
				continue
			}
			emoteFile, err := os.ReadFile(file)
			encodedImage := base64.StdEncoding.EncodeToString(emoteFile)
			dataURI := fmt.Sprintf("data:image/png;base64,%s", encodedImage)

			_, err = s.GuildEmojiCreate(i.GuildID, &discordgo.EmojiParams{
				Name:  emoteName,
				Image: dataURI,
			})
			if err != nil {
				output.WriteString(fmt.Sprintf("%s - upload error: %s\n", emoteName, err))
				continue
			}
			output.WriteString(fmt.Sprintf("%s - success\n", emoteName))
		}
		output.WriteString("```")
	}

	_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: output.String(),
	})
}
