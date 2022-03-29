package main

import (
	_ "embed"
	"log"
	"strings"

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
		Args:           args,
		CompsByMessage: make(map[string]*Comp),
		CompsByOwner:   make(map[string]*Comp),
	}
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
		comp := MakeComp("", m.Author.ID, m.Author.String())
		msg, err := s.ChannelMessageSendEmbed(m.ChannelID, comp.Embed())
		if err != nil {
			bot.Logger.Printf("Message could not be sent in %s", m.ChannelID)
			return
		}
		comp.Id = msg.ID
		bot.CompsByMessage[msg.ID] = comp
		s.MessageReactionAdd(m.ChannelID, msg.ID, "🆗")
		bot.Logger.Printf("User %s Created a new comp", m.Author.String())
	}
}

func (bot *CompBot) reactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID || m.ChannelID != bot.Args.Channel || m.Emoji.Name != "🆗" {
		return
	}
	if comp, ok := bot.CompsByMessage[m.MessageID]; ok {
		err := comp.AddUser(m.UserID, m.Member.User.String())
		if err != nil {
			bot.Logger.Print(err)
			return
		}
		_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, comp.Embed())
		if err != nil {
			bot.Logger.Printf("Message could not be edited in %s", m.ChannelID)
			return
		}
		bot.Logger.Printf("User %s Joined Comp %s", m.Member.User.String(), m.MessageID)
	}
}

func (bot *CompBot) reactionRemove(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m.UserID == s.State.User.ID || m.ChannelID != bot.Args.Channel || m.Emoji.Name != "🆗" {
		return
	}
	if comp, ok := bot.CompsByMessage[m.MessageID]; ok {
		err := comp.RemoveUser(m.UserID)
		if err != nil {
			bot.Logger.Print(err)
			return
		}
		_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, comp.Embed())
		if err != nil {
			bot.Logger.Printf("Message could not be edited in %s", m.ChannelID)
			return
		}
		bot.Logger.Printf("User %s Left Comp %s", m.UserID, m.MessageID)
	}
}
