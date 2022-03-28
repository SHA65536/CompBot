# CompBot
A Discord bot for competitive recruitment 

## Installation
Clone this repository and navigate into the created folder:
```
git clone https://github.com/SHA65536/CompBot.git
cd CompBot
```
Now install the dependencies:
```
go get .
```
Now you will need to create the envirnment variables for the bot.

You could either export them on your own, or create a file named ".env" containing the following values:
```
TOKEN="YOUR_TOKEN"
PREFIX="!comp"
CHANNEL="YOUR_CHANNEL_ID"
```

## Usage
To run the bot just run 
```
go run main.go
```
To create a comp, use the keyword define earlier "!comp"
The Bot will create a message telling everyone you are looking for partners. Now users wishing to join the comp have to click the 🆗 reaction to join. Clicking the 🆗 reaction again will remove a user from the comp.