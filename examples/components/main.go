package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ayntgl/astatine"
)

// Bot parameters
var (
	GuildID  = flag.String("guild", "", "Test guild ID")
	BotToken = flag.String("token", "", "Bot access token")
	AppID    = flag.String("app", "", "Application ID")
)

var s *astatine.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = astatine.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

// Important note: call every command in order it's placed in the example.

var (
	componentsHandlers = map[string]func(s *astatine.Session, i *astatine.InteractionCreate){
		"fd_no": func(s *astatine.Session, i *astatine.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &astatine.InteractionResponse{
				Type: astatine.InteractionResponseChannelMessageWithSource,
				Data: &astatine.InteractionResponseData{
					Content: "Huh. I see, maybe some of these resources might help you?",
					Flags:   1 << 6,
					Components: []astatine.MessageComponent{
						astatine.ActionsRow{
							Components: []astatine.MessageComponent{
								astatine.Button{
									Emoji: astatine.ComponentEmoji{
										Name: "📜",
									},
									Label: "Documentation",
									Style: astatine.LinkButton,
									URL:   "https://discord.com/developers/docs/interactions/message-components#buttons",
								},
								astatine.Button{
									Emoji: astatine.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: astatine.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								astatine.Button{
									Emoji: astatine.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: astatine.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"fd_yes": func(s *astatine.Session, i *astatine.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &astatine.InteractionResponse{
				Type: astatine.InteractionResponseChannelMessageWithSource,
				Data: &astatine.InteractionResponseData{
					Content: "Great! If you wanna know more or just have questions, feel free to visit Discord Devs and Discord Gophers server. " +
						"But now, when you know how buttons work, let's move onto select menus (execute `/selects single`)",
					Flags: 1 << 6,
					Components: []astatine.MessageComponent{
						astatine.ActionsRow{
							Components: []astatine.MessageComponent{
								astatine.Button{
									Emoji: astatine.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: astatine.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								astatine.Button{
									Emoji: astatine.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: astatine.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"select": func(s *astatine.Session, i *astatine.InteractionCreate) {
			var response *astatine.InteractionResponse

			data := i.MessageComponentData()
			switch data.Values[0] {
			case "go":
				response = &astatine.InteractionResponse{
					Type: astatine.InteractionResponseChannelMessageWithSource,
					Data: &astatine.InteractionResponseData{
						Content: "This is the way.",
						Flags:   1 << 6,
					},
				}
			default:
				response = &astatine.InteractionResponse{
					Type: astatine.InteractionResponseChannelMessageWithSource,
					Data: &astatine.InteractionResponseData{
						Content: "It is not the way to go.",
						Flags:   1 << 6,
					},
				}
			}
			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second) // Doing that so user won't see instant response.
			_, err = s.FollowupMessageCreate(*AppID, i.Interaction, true, &astatine.WebhookParams{
				Content: "Anyways, now when you know how to use single select menus, let's see how multi select menus work. " +
					"Try calling `/selects multi` command.",
				Flags: 1 << 6,
			})
			if err != nil {
				panic(err)
			}
		},
		"stackoverflow_tags": func(s *astatine.Session, i *astatine.InteractionCreate) {
			data := i.MessageComponentData()

			const stackoverflowFormat = `https://stackoverflow.com/questions/tagged/%s`

			err := s.InteractionRespond(i.Interaction, &astatine.InteractionResponse{
				Type: astatine.InteractionResponseChannelMessageWithSource,
				Data: &astatine.InteractionResponseData{
					Content: "Here is your stackoverflow URL: " + fmt.Sprintf(stackoverflowFormat, strings.Join(data.Values, "+")),
					Flags:   1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second) // Doing that so user won't see instant response.
			_, err = s.FollowupMessageCreate(*AppID, i.Interaction, true, &astatine.WebhookParams{
				Content: "Now you know everything about select component. If you want to know more or ask a question - feel free to.",
				Components: []astatine.MessageComponent{
					astatine.ActionsRow{
						Components: []astatine.MessageComponent{
							astatine.Button{
								Emoji: astatine.ComponentEmoji{
									Name: "📜",
								},
								Label: "Documentation",
								Style: astatine.LinkButton,
								URL:   "https://discord.com/developers/docs/interactions/message-components#select-menus",
							},
							astatine.Button{
								Emoji: astatine.ComponentEmoji{
									Name: "🔧",
								},
								Label: "Discord developers",
								Style: astatine.LinkButton,
								URL:   "https://discord.gg/discord-developers",
							},
							astatine.Button{
								Emoji: astatine.ComponentEmoji{
									Name: "🦫",
								},
								Label: "Discord Gophers",
								Style: astatine.LinkButton,
								URL:   "https://discord.gg/7RuRrVHyXF",
							},
						},
					},
				},
				Flags: 1 << 6,
			})
			if err != nil {
				panic(err)
			}
		},
	}
	commandsHandlers = map[string]func(s *astatine.Session, i *astatine.InteractionCreate){
		"buttons": func(s *astatine.Session, i *astatine.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &astatine.InteractionResponse{
				Type: astatine.InteractionResponseChannelMessageWithSource,
				Data: &astatine.InteractionResponseData{
					Content: "Are you comfortable with buttons and other message components?",
					Flags:   1 << 6,
					// Buttons and other components are specified in Components field.
					Components: []astatine.MessageComponent{
						// ActionRow is a container of all buttons within the same row.
						astatine.ActionsRow{
							Components: []astatine.MessageComponent{
								astatine.Button{
									// Label is what the user will see on the button.
									Label: "Yes",
									// Style provides coloring of the button. There are not so many styles tho.
									Style: astatine.SuccessButton,
									// Disabled allows bot to disable some buttons for users.
									Disabled: false,
									// CustomID is a thing telling Discord which data to send when this button will be pressed.
									CustomID: "fd_yes",
								},
								astatine.Button{
									Label:    "No",
									Style:    astatine.DangerButton,
									Disabled: false,
									CustomID: "fd_no",
								},
								astatine.Button{
									Label:    "I don't know",
									Style:    astatine.LinkButton,
									Disabled: false,
									// Link buttons don't require CustomID and do not trigger the gateway/HTTP event
									URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
									Emoji: astatine.ComponentEmoji{
										Name: "🤷",
									},
								},
							},
						},
						// The message may have multiple actions rows.
						astatine.ActionsRow{
							Components: []astatine.MessageComponent{
								astatine.Button{
									Label:    "Discord Developers server",
									Style:    astatine.LinkButton,
									Disabled: false,
									URL:      "https://discord.gg/discord-developers",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"selects": func(s *astatine.Session, i *astatine.InteractionCreate) {
			var response *astatine.InteractionResponse
			switch i.ApplicationCommandData().Options[0].Name {
			case "single":
				response = &astatine.InteractionResponse{
					Type: astatine.InteractionResponseChannelMessageWithSource,
					Data: &astatine.InteractionResponseData{
						Content: "Now let's take a look on selects. This is single item select menu.",
						Flags:   1 << 6,
						Components: []astatine.MessageComponent{
							astatine.ActionsRow{
								Components: []astatine.MessageComponent{
									astatine.SelectMenu{
										// Select menu, as other components, must have a customID, so we set it to this value.
										CustomID:    "select",
										Placeholder: "Choose your favorite programming language 👇",
										Options: []astatine.SelectMenuOption{
											{
												Label: "Go",
												// As with components, this things must have their own unique "id" to identify which is which.
												// In this case such id is Value field.
												Value: "go",
												Emoji: astatine.ComponentEmoji{
													Name: "🦦",
												},
												// You can also make it a default option, but in this case we won't.
												Default:     false,
												Description: "Go programming language",
											},
											{
												Label: "JS",
												Value: "js",
												Emoji: astatine.ComponentEmoji{
													Name: "🟨",
												},
												Description: "JavaScript programming language",
											},
											{
												Label: "Python",
												Value: "py",
												Emoji: astatine.ComponentEmoji{
													Name: "🐍",
												},
												Description: "Python programming language",
											},
										},
									},
								},
							},
						},
					},
				}
			case "multi":
				minValues := 1
				response = &astatine.InteractionResponse{
					Type: astatine.InteractionResponseChannelMessageWithSource,
					Data: &astatine.InteractionResponseData{
						Content: "The tastiest things are left for the end. Let's see how the multi-item select menu works: " +
							"try generating your own stackoverflow search link",
						Flags: 1 << 6,
						Components: []astatine.MessageComponent{
							astatine.ActionsRow{
								Components: []astatine.MessageComponent{
									astatine.SelectMenu{
										CustomID:    "stackoverflow_tags",
										Placeholder: "Select tags to search on StackOverflow",
										// This is where confusion comes from. If you don't specify these things you will get single item select.
										// These fields control the minimum and maximum amount of selected items.
										MinValues: &minValues,
										MaxValues: 3,
										Options: []astatine.SelectMenuOption{
											{
												Label:       "Go",
												Description: "Simple yet powerful programming language",
												Value:       "go",
												// Default works the same for multi-select menus.
												Default: false,
												Emoji: astatine.ComponentEmoji{
													Name: "🦦",
												},
											},
											{
												Label:       "JS",
												Description: "Multiparadigm OOP language",
												Value:       "javascript",
												Emoji: astatine.ComponentEmoji{
													Name: "🟨",
												},
											},
											{
												Label:       "Python",
												Description: "OOP prototyping programming language",
												Value:       "python",
												Emoji: astatine.ComponentEmoji{
													Name: "🐍",
												},
											},
											{
												Label:       "Web",
												Description: "Web related technologies",
												Value:       "web",
												Emoji: astatine.ComponentEmoji{
													Name: "🌐",
												},
											},
											{
												Label:       "Desktop",
												Description: "Desktop applications",
												Value:       "desktop",
												Emoji: astatine.ComponentEmoji{
													Name: "💻",
												},
											},
										},
									},
								},
							},
						},
					},
				}

			}
			err := s.InteractionRespond(i.Interaction, response)
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
	// Components are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *astatine.Session, i *astatine.InteractionCreate) {
		switch i.Type {
		case astatine.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case astatine.InteractionMessageComponent:

			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
	_, err := s.ApplicationCommandCreate(*AppID, *GuildID, &astatine.ApplicationCommand{
		Name:        "buttons",
		Description: "Test the buttons if you got courage",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
	_, err = s.ApplicationCommandCreate(*AppID, *GuildID, &astatine.ApplicationCommand{
		Name: "selects",
		Options: []*astatine.ApplicationCommandOption{
			{
				Type:        astatine.ApplicationCommandOptionSubCommand,
				Name:        "multi",
				Description: "Multi-item select menu",
			},
			{
				Type:        astatine.ApplicationCommandOptionSubCommand,
				Name:        "single",
				Description: "Single-item select menu",
			},
		},
		Description: "Lo and behold: dropdowns are coming",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
