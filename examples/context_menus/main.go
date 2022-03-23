package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/ayntgl/astatine"
)

// Bot parameters
var (
	GuildID  = flag.String("guild", "", "Test guild ID")
	BotToken = flag.String("token", "", "Bot access token")
	AppID    = flag.String("app", "", "Application ID")
	Cleanup  = flag.Bool("cleanup", true, "Cleanup of commands")
)

var s *astatine.Session

func init() {
	flag.Parse()
	s = astatine.New("Bot " + *BotToken)
}

func searchLink(message, format, sep string) string {
	return fmt.Sprintf(format, strings.Join(
		strings.Split(
			message,
			" ",
		),
		sep,
	))
}

var (
	commands = []astatine.ApplicationCommand{
		{
			Name: "rickroll-em",
			Type: astatine.UserApplicationCommand,
		},
		{
			Name: "google-it",
			Type: astatine.MessageApplicationCommand,
		},
		{
			Name: "stackoverflow-it",
			Type: astatine.MessageApplicationCommand,
		},
		{
			Name: "godoc-it",
			Type: astatine.MessageApplicationCommand,
		},
		{
			Name: "discordjs-it",
			Type: astatine.MessageApplicationCommand,
		},
		{
			Name: "discordpy-it",
			Type: astatine.MessageApplicationCommand,
		},
	}
	commandsHandlers = map[string]func(s *astatine.Session, i *astatine.InteractionCreate){
		"rickroll-em": func(s *astatine.Session, i *astatine.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &astatine.InteractionResponse{
				Type: astatine.InteractionResponseChannelMessageWithSource,
				Data: &astatine.InteractionResponseData{
					Content: "Operation rickroll has begun",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}

			ch, err := s.UserChannelCreate(
				i.ApplicationCommandData().TargetID,
			)
			if err != nil {
				_, err = s.FollowupMessageCreate(*AppID, i.Interaction, true, &astatine.WebhookParams{
					Content: fmt.Sprintf("Mission failed. Cannot send a message to this user: %q", err.Error()),
					Flags:   1 << 6,
				})
				if err != nil {
					panic(err)
				}
			}
			_, err = s.ChannelMessageSend(
				ch.ID,
				fmt.Sprintf("%s sent you this: https://youtu.be/dQw4w9WgXcQ", i.Member.Mention()),
			)
			if err != nil {
				panic(err)
			}
		},
		"google-it": func(s *astatine.Session, i *astatine.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &astatine.InteractionResponse{
				Type: astatine.InteractionResponseChannelMessageWithSource,
				Data: &astatine.InteractionResponseData{
					Content: searchLink(
						i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content,
						"https://google.com/search?q=%s", "+"),
					Flags: 1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"stackoverflow-it": func(s *astatine.Session, i *astatine.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &astatine.InteractionResponse{
				Type: astatine.InteractionResponseChannelMessageWithSource,
				Data: &astatine.InteractionResponseData{
					Content: searchLink(
						i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content,
						"https://stackoverflow.com/search?q=%s", "+"),
					Flags: 1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"godoc-it": func(s *astatine.Session, i *astatine.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &astatine.InteractionResponse{
				Type: astatine.InteractionResponseChannelMessageWithSource,
				Data: &astatine.InteractionResponseData{
					Content: searchLink(
						i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content,
						"https://pkg.go.dev/search?q=%s", "+"),
					Flags: 1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"discordjs-it": func(s *astatine.Session, i *astatine.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &astatine.InteractionResponse{
				Type: astatine.InteractionResponseChannelMessageWithSource,
				Data: &astatine.InteractionResponseData{
					Content: searchLink(
						i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content,
						"https://discord.js.org/#/docs/main/stable/search?query=%s", "+"),
					Flags: 1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"discordpy-it": func(s *astatine.Session, i *astatine.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &astatine.InteractionResponse{
				Type: astatine.InteractionResponseChannelMessageWithSource,
				Data: &astatine.InteractionResponseData{
					Content: searchLink(
						i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content,
						"https://discordpy.readthedocs.io/en/stable/search.html?q=%s", "+"),
					Flags: 1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
)

func main() {
	s.AddHandler(func(s *astatine.Session, r *astatine.Ready) {
		log.Println("Bot is up!")
	})

	s.AddHandler(func(s *astatine.Session, i *astatine.InteractionCreate) {
		if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	cmdIDs := make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := s.ApplicationCommandCreate(*AppID, *GuildID, &cmd)
		if err != nil {
			log.Fatalf("Cannot create slash command %q: %v", cmd.Name, err)
		}

		cmdIDs[rcmd.ID] = rcmd.Name

	}

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

	if !*Cleanup {
		return
	}

	for id, name := range cmdIDs {
		err := s.ApplicationCommandDelete(*AppID, *GuildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}

}
