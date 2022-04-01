# CompBot
A Discord bot for competitive recruitment 

## Installation
Clone this repository and navigate into the created folder:
```shell
git clone https://github.com/SHA65536/CompBot.git
cd CompBot
```
Now install the dependencies:
```shell
go get .
```
Now you will need to create the envirnment variables for the bot.

You could either export them on your own, or create a file named ".env" containing the following values:
```shell
TOKEN="YOUR_TOKEN"
CHANNEL="YOUR_CHANNEL_ID"
PREFIX="!prefix" #Optional. Default: !comp
#Cooldown format is like "300ms", "1.5h" or "2h45m". 
#Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
CREATE_CD="5m" #Optional. Default: 5m
REACT_CD="3s" #Optional. Default: 3s
```

## Usage
To run the bot just run 
```shell
go run .
```
To create a comp, use the keyword define earlier "!comp"
The Bot will create a message telling everyone you are looking for partners. Now users wishing to join the comp have to click the 🆗 reaction to join. Clicking the 🆗 reaction again will remove a user from the comp.

## Custom Messages
To change the text inside the comps, you need to modify `emptyEmbed.json` and `fullEmbed.json`

**emptyEmbed.json**
```json
{
	"content": "@everyone", //This is what's above the embed.
	"embeds": [
		{
			"color": 12846604,
			// %v here is the number of people who joined.
			// %s here is the numbered list of people who joined.
			"description": "**%v/5 have volunteered!**\n%s",
			"footer": {
				"text": "Press 🆗 To Join!"
			},
			"author":{
				// %s here is the Creator's name.
				"name": "%s Is Orginaizing an Attack Squad!"
			},
			// %s here is the Creator's name, this field will be overridden
			// if the user supplied a comp title when creating.
			"title": "%s Is Orginaizing an Attack Squad!"
		}
	]
}
```
**fullEmbed.json**
```json
{
	// %s here is the list of mentions
	"content": "%s",
	"embeds": [
		{
			"color": 65280,
			"footer": {
				"text": "Press 🆗 To Join!"
			},
			"author":{
				// %s here is the Creator's name.
				"name": "%s's Comp Is Ready!"
			},
			//this field will be overridden if the user supplied a comp title when creating.
			"title": "The Comp Is Ready!",
			// %s here is the numbered list of people who joined.
			"description":"5/5 Are ready!\n%s"
		}
	]
}
```