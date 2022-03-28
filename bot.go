package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//go:embed embedTemplate.json
var template []byte

type Config struct {
	Token   string
	Prefix  string
	Channel string
}

type Comp struct {
	Id    string            // Message ID
	Owner string            // Owner ID
	Users map[string]string // Participants
	Title string            // Embed Title
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
		bot.createComp(s, m)
	}
}

// Creates a new comp and registers it
func (bot *CompBot) createComp(s *discordgo.Session, m *discordgo.MessageCreate) {
	var embed *discordgo.MessageEmbed
	var send *discordgo.MessageSend
	var comp *Comp
	embed = &discordgo.MessageEmbed{}
	json.Unmarshal(template, &embed)
	embed.Title = fmt.Sprintf(embed.Title, m.Author.Username)
	send = &discordgo.MessageSend{
		Embed: embed,
		TTS:   false,
	}
	msg, err := s.ChannelMessageSendComplex(m.ChannelID, send)
	if err != nil {
		bot.Logger.Printf("Message could not be sent in %s", m.ChannelID)
		return
	}
	comp = &Comp{msg.ID, m.Author.ID, make(map[string]string), embed.Title}
	bot.CompsByMessage[msg.ID] = comp
	s.MessageReactionAdd(m.ChannelID, msg.ID, "🆗")
}

func (bot *CompBot) reactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID || m.ChannelID != bot.Args.Channel || m.Emoji.Name != "🆗" {
		return
	}
	if comp, ok := bot.CompsByMessage[m.MessageID]; ok {
		comp.Users[m.UserID] = m.Member.User.Username
		embed := &discordgo.MessageEmbed{}
		json.Unmarshal(template, &embed)
		embed.Title = comp.Title
		embed.Description = fmt.Sprintf("**%v have volunteered!**\n%s", len(comp.Users), dictToList(comp.Users))
		_, err := s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, embed)
		if err != nil {
			bot.Logger.Printf("Message could not be edited in %s", m.ChannelID)
			return
		}
	}
}

func dictToList(d map[string]string) string {
	keys := make([]string, 0, len(d))
	for k := range d {
		keys = append(keys, d[k])
	}
	return strings.Join(keys, ",")
}
