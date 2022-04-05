package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gos-auto-bot/bot"
	"gos-auto-bot/orchestrator"
	"gos-auto-bot/task"
	"regexp"
)

const (
	PingCommand     = "!ping"
	TidyCommand     = "!tidy"
	HelpCommand     = "!help"
	ScheduleCommand = "!schedule"
	CommandsCommand = "!commands"
	Help            = "!help"

	ModeratorRole = "Moderator"
)

type Handler func(s *discordgo.Session, m *discordgo.MessageCreate)

func CommandsHandler(l logrus.FieldLogger, commands []string) Handler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		commandsStr := ""
		for _, commandStr := range commands {
			if commandsStr != "" {
				commandsStr += ", "
			}
			commandsStr += commandStr
		}

		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Currently configured commands = %s.", commandsStr))
		if err != nil {
			l.WithError(err).Errorf("Unable to send %s response.", CommandsCommand)
			return
		}
	}
}

func HelpHandler(l logrus.FieldLogger, _ string, instructions chan<- orchestrator.Instruction) Handler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		names := []string{"Tin", "Bender"}
		for _, n := range names {
			i := orchestrator.Instruction{
				Name: n,
				Task: task.ImperialCarnivalHelpTask,
			}

			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Scheduling task %s for %s. This may take up to a minute to perform.", task.ImperialCarnivalHelpTask, n))
			if err != nil {
				l.WithError(err).Errorf("Unable to send %s response.", HelpCommand)
				return
			}

			instructions <- i
		}
	}
}

func TidyHandler(l logrus.FieldLogger, guildId string, channelId string) Handler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if !isRole(l)(s, guildId, ModeratorRole, m.Member.Roles) {
			l.Debugf("Unauthorized user %s trying to execute privileged command.", m.Author.Username)
			return
		}

		cleaned := 0
		cleaning := true
		for cleaning {
			ms, err := s.ChannelMessages(channelId, 100, "", "", "")
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve messages for channel.")
				return
			}
			if len(ms) > 0 {
				cleaned += len(ms)
				ids := make([]string, 0)
				for _, cm := range ms {
					ids = append(ids, cm.ID)
				}
				err = s.ChannelMessagesBulkDelete(channelId, ids)
				if err != nil {
					l.WithError(err).Errorf("Unable to bulk delete messages for channel.")
				}
			} else {
				cleaning = false
			}
		}
		_, err := s.ChannelMessageSend(channelId, fmt.Sprintf("Cleaned %d messages at the request of <@%s>.", cleaned, m.Author.ID))
		if err != nil {
			l.WithError(err).Errorf("Unable to send %s response.", TidyCommand)
		}
	}
}

func PingHandler(l logrus.FieldLogger) Handler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			l.WithError(err).Errorf("Unable to send %s response.", PingCommand)
			return
		}
	}
}

func isRole(l logrus.FieldLogger) func(s *discordgo.Session, guildId string, roleName string, memberRoleIds []string) bool {
	return func(s *discordgo.Session, guildId string, roleName string, memberRoleIds []string) bool {
		rs, err := s.GuildRoles(guildId)
		if err != nil {
			l.WithError(err).Errorf("Unable to look up guild roles.")
			return false
		}
		var role *discordgo.Role
		for _, r := range rs {
			if r.Name == roleName {
				role = r
			}
		}
		if role == nil {
			l.Errorf("Unable to locate %s role.", roleName)
			return false
		}
		is := false
		for _, rid := range memberRoleIds {
			if rid == role.ID {
				is = true
			}
		}
		return is
	}
}

func ScheduleHandler(l logrus.FieldLogger, guildId string, channelId string, instructions chan<- orchestrator.Instruction) Handler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if !isRole(l)(s, guildId, ModeratorRole, m.Member.Roles) {
			l.Debugf("Unauthorized user %s trying to execute privileged command.", m.Author.Username)
			return
		}
		name, taskName := parseSchedule(m.Message.Content)
		if name == "BULK" {
			names, err := bot.ReadNames(l)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve list of bots.")
				return
			}
			for _, n := range names {
				i := orchestrator.Instruction{
					Name: n,
					Task: taskName,
				}
				instructions <- i
			}
		} else {
			i := orchestrator.Instruction{
				Name: name,
				Task: taskName,
			}
			instructions <- i
		}

		_, err := s.ChannelMessageSend(channelId, fmt.Sprintf("Scheduling %s for %s at the request of <@%s>.", taskName, name, m.Author.ID))
		if err != nil {
			l.WithError(err).Errorf("Unable to send %s response.", ScheduleCommand)
		}
	}
}

func parseSchedule(message string) (string, string) {
	re := regexp.MustCompile("<(.*)> !schedule ([a-zA-Z 0-9]*)* ([a-zA-Z]*)")
	match := re.FindStringSubmatch(message)
	return match[2], match[3]
}
