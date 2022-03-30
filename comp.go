package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

//go:embed embedTemplate.json
var template []byte

type User struct {
	Id   string
	Name string
}

type Comp struct {
	Id    string          // Message ID
	Owner User            // Owner ID
	Users map[string]bool // Participants by ID
	Chron []User          // Participants Chronologically
}

func MakeComp(msg, owner_id, owner_name string) *Comp {
	return &Comp{
		Id:    msg,
		Owner: User{owner_id, owner_name},
		Users: make(map[string]bool),
		Chron: make([]User, 0),
	}
}

func (c *Comp) AddUser(id, name string) error {
	if _, ok := c.Users[id]; ok {
		return errors.New("comp: user already in comp.")
	}
	c.Users[id] = true
	c.Chron = append(c.Chron, User{id, name})
	return nil
}

func (c *Comp) RemoveUser(id string) error {
	if _, ok := c.Users[id]; !ok {
		return errors.New("comp: user not in comp.")
	}
	delete(c.Users, id)
	for i, user := range c.Chron {
		if id == user.Id {
			c.Chron = append(c.Chron[:i], c.Chron[i+1:]...)
			break
		}
	}
	return nil
}

func (c *Comp) Embed() *discordgo.MessageSend {
	embed := &discordgo.MessageSend{}
	json.Unmarshal(template, &embed)
	embed.Embeds[0].Title = fmt.Sprintf(embed.Embeds[0].Title, c.Owner.Name)
	embed.Embeds[0].Description = fmt.Sprintf("**%v/5 have volunteered!**\n%s", len(c.Users), c.nameList())
	return embed
}

func (c *Comp) nameList() string {
	res := ""
	for i, user := range c.Chron {
		res += fmt.Sprintf("%v - %s\n", i+1, user.Name)
	}
	return res
}
