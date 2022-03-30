package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token   string
	Prefix  string
	Channel string
}

type CompBot struct {
	Session        *discordgo.Session
	Logger         *log.Logger
	Cooldowns      *CooldownManager
	Args           Config
	CompsByMessage map[string]*Comp
	CompsByOwner   map[string]*Comp
}

func MakeCompBot(args Config) *CompBot {
	disc, err := discordgo.New("Bot " + args.Token)
	if err != nil {
		log.Fatal(err)
	}
	bot := &CompBot{
		Session:        disc,
		Logger:         log.Default(),
		Cooldowns:      MakeCoolownManager(),
		Args:           args,
		CompsByMessage: make(map[string]*Comp),
		CompsByOwner:   make(map[string]*Comp),
	}

	bot.Cooldowns.NewCooldown("create", time.Minute*5)
	bot.Cooldowns.NewCooldown("react", time.Second*3)

	disc.AddHandler(bot.messageCreate)
	disc.AddHandler(bot.reactionAdd)
	disc.AddHandler(bot.reactionRemove)
	disc.Identify.Intents =
		discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageReactions
	return bot
}

func (bot *CompBot) Start() {
	err := bot.Session.Open()
	if err != nil {
		log.Fatal(err)
	}
	bot.Logger.Println("Bot Successfully connected")
}

func (bot *CompBot) Stop() {
	bot.Session.Close()
	bot.Logger.Println("Bot disconnected")
}

// Handles Messages
func (bot *CompBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.ChannelID != bot.Args.Channel {
		return
	}
	if strings.HasPrefix(m.Content, bot.Args.Prefix) {
		if ok, err := bot.Cooldowns.IsAllowed("create", m.Author.ID); ok {
			bot.createComp(s, m)
			err := bot.Cooldowns.SetObject("create", m.Author.ID)
			if err != nil {
				bot.Logger.Print(err)
			}
			return
		} else {
			if err != nil {
				bot.Logger.Print(err)
			} else {
				bot.Logger.Printf("User %s is create on cooldown", m.Author.String())
			}
			return
		}
	}
}

// Handles reaction adds
func (bot *CompBot) reactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID || m.ChannelID != bot.Args.Channel || m.Emoji.Name != "🆗" {
		return
	}
	if comp, ok := bot.CompsByMessage[m.MessageID]; ok {
		if ok, err := bot.Cooldowns.IsAllowed("react", m.UserID); ok {
			err = bot.joinComp(s, m, comp)
			if err != nil {
				bot.Logger.Print(err)
				return
			}
			err = bot.Cooldowns.SetObject("react", m.UserID)
			if err != nil {
				bot.Logger.Print(err)
				return
			}
		} else {
			if err != nil {
				bot.Logger.Print(err)
				return
			} else {
				bot.Logger.Printf("User %s is react on cooldown", m.UserID)
				return
			}
		}
	}
}

// Handles reaction removes
func (bot *CompBot) reactionRemove(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m.UserID == s.State.User.ID || m.ChannelID != bot.Args.Channel || m.Emoji.Name != "🆗" {
		return
	}
	if comp, ok := bot.CompsByMessage[m.MessageID]; ok {
		if ok, err := bot.Cooldowns.IsAllowed("react", m.UserID); ok {
			err = bot.leaveComp(s, m, comp)
			if err != nil {
				bot.Logger.Print(err)
				return
			}
			err = bot.Cooldowns.SetObject("react", m.UserID)
			if err != nil {
				bot.Logger.Print(err)
				return
			}
		} else {
			if err != nil {
				bot.Logger.Print(err)
				return
			} else {
				bot.Logger.Printf("User %s is react on cooldown", m.UserID)
				return
			}
		}
	}
}

func (bot *CompBot) createComp(s *discordgo.Session, m *discordgo.MessageCreate) {
	comp := MakeComp("", m.Author.ID, m.Author.String())
	msg, err := s.ChannelMessageSendComplex(m.ChannelID, comp.Embed())
	if err != nil {
		bot.Logger.Printf("Message could not be sent in %s", m.ChannelID)
		return
	}
	comp.Id = msg.ID
	bot.CompsByMessage[msg.ID] = comp
	s.MessageReactionAdd(m.ChannelID, msg.ID, "🆗")
	bot.Logger.Printf("User %s Created a new comp", m.Author.String())
}

func (bot *CompBot) joinComp(s *discordgo.Session, m *discordgo.MessageReactionAdd, c *Comp) error {
	err := c.AddUser(m.UserID, m.Member.User.String())
	if err != nil {
		return err
	}
	bot.Logger.Printf("User %s Joined Comp %s", m.Member.User.String(), m.MessageID)
	if len(c.Users) == CompSize {
		s.ChannelMessageDelete(m.ChannelID, c.Id)
		msg, err := s.ChannelMessageSendComplex(m.ChannelID, c.Embed())
		if err != nil {
			return fmt.Errorf("Message could not be sent in %s", m.ChannelID)
		}
		c.Id = msg.ID
		bot.CompsByMessage[msg.ID] = c
		s.MessageReactionAdd(m.ChannelID, msg.ID, "🆗")
		bot.Logger.Printf("Comp %s is ready!", msg.ID)
	} else {
		_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, c.Embed().Embeds[0])
		if err != nil {
			return fmt.Errorf("Message could not be edited in %s", m.ChannelID)
		}
	}
	return nil
}

func (bot *CompBot) leaveComp(s *discordgo.Session, m *discordgo.MessageReactionRemove, c *Comp) error {
	err := c.RemoveUser(m.UserID)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, c.Embed().Embeds[0])
	if err != nil {
		return fmt.Errorf("Message could not be edited in %s", m.ChannelID)
	}
	bot.Logger.Printf("User %s Left Comp %s", m.UserID, m.MessageID)
	return nil
}
