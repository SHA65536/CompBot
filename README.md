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
