# echo-chamber

echo-chamber allows you to type keywords and reads twitch chat messages for those keywords. 

* Elasticsearch for storing messages
* Go go-twitch-irc and relaybroker to read messages from many chats at once
* Deployable with Docker

## Development

Running echo-chamber requires docker and docker-compose. To run it locally remove the ".dist" from chamber.env.dist and add your clientId from twitch.
After that all you have to type is 
```
docker-compose up
```

