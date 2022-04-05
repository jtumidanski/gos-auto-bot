package discord

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gos-auto-bot/command"
	"gos-auto-bot/orchestrator"
	"os"
	"regexp"
	"sync"
)

func NewWorker(fl logrus.FieldLogger) func(ctx context.Context, wg *sync.WaitGroup, instructions chan<- orchestrator.Instruction) {
	l := fl.WithField("thread", "discord")

	ownerId, ok := os.LookupEnv("DISCORD_OWNER_ID")
	if !ok {
		l.Fatalf("Unable to lookup DISCORD_OWNER_ID")
	}
	guildId, ok := os.LookupEnv("DISCORD_GUILD_ID")
	if !ok {
		l.Fatalf("Unable to lookup DISCORD_GUILD_ID")
	}
	username, ok := os.LookupEnv("DISCORD_USERNAME")
	if !ok {
		l.Fatalf("Unable to lookup DISCORD_USERNAME")
	}
	channel, ok := os.LookupEnv("DISCORD_CHANNEL")
	if !ok {
		l.Fatalf("Unable to lookup DISCORD_CHANNEL")
	}
	token, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		l.Fatalf("Unable to lookup DISCORD_TOKEN")
	}

	return func(ctx context.Context, wg *sync.WaitGroup, instructions chan<- orchestrator.Instruction) {
		wg.Add(1)

		l.Infof("Discord bot initializing.")

		dg, err := discordgo.New("Bot " + token)
		if err != nil {
			l.WithError(err).Warnf("Error creating Discord session.")
			return
		}

		commandMap := make(map[string]command.Handler)
		commandMap[command.PingCommand] = command.PingHandler(l)
		commandMap[command.TidyCommand] = command.TidyHandler(l, guildId, channel)
		commandMap[command.HelpCommand] = command.HelpHandler(l, ownerId, instructions)
		commandMap[command.ScheduleCommand] = command.ScheduleHandler(l, guildId, channel, instructions)
		commandMap[command.CommandsCommand] = command.CommandsHandler(l, []string{command.PingCommand, command.TidyCommand, command.HelpCommand, command.ScheduleCommand})

		dg.AddHandler(messageCreate(l)(username, commandMap))

		dg.Identify.Intents = discordgo.IntentsGuildMessages

		err = dg.Open()
		if err != nil {
			l.WithError(err).Warnf("Error opening connection.")
			return
		}

		_, err = dg.ChannelMessageSend(channel, "Reporting for duty.")
		if err != nil {
			l.WithError(err).Warnf("Unable to report in for duty.")
			return
		}

		<-ctx.Done()

		l.Infof("Discord bot shutting down.")

		_, err = dg.ChannelMessageSend(channel, "Leaving my duties.")
		if err != nil {
			l.WithError(err).Warnf("Unable to report termination.")
			return
		}

		err = dg.Close()
		if err != nil {
			l.WithError(err).Warnf("Unable to close gracefully.")
			return
		}

		wg.Done()
	}
}

func messageCreate(l logrus.FieldLogger) func(username string, commandHandlers map[string]command.Handler) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(username string, commandHandlers map[string]command.Handler) func(s *discordgo.Session, m *discordgo.MessageCreate) {
		return func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.Author.ID == username {
				return
			}
			if !mentioned(username, m.Message) {
				return
			}

			l.Debugf("Mentioned in message %s.", m.Content)

			c := parseCommand(m.Content)
			if handler, ok := commandHandlers[c]; ok {
				handler(s, m)
			}
		}
	}
}

func parseCommand(message string) string {
	re := regexp.MustCompile("<(.*)> (![a-zA-Z]*)")
	match := re.FindStringSubmatch(message)
	return match[2]
}

func mentioned(username string, message *discordgo.Message) bool {
	for _, mention := range message.Mentions {
		if mention.Username == username {
			return true
		}
	}
	return false
}
