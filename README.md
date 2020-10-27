# PokerMUD
PokerMUD is an implementation of Texas Hold 'em Poker written in Go.
It is meant to be a MUD-like game in that the client is the command line.
The reason for this is that I thought it would be cool if you could have
a simple lightweight Poker game to play with friends.

## Motivation
The motivation for writing this game was primarily for the challenge of writing a complete game.
Games are interesting to write because they have many edge cases, have clear real world
models that make you think about proper software design, and I wanted to learn more
about network programming using raw TCP sockets (otherwise I would have done something simple
and used an abstraction like websockets)

## Design
### Deck
The deck is a struct that has a slice of Cards. At the beginning of each new game,
a new Card slice is created with 52 cards and shuffled. At the beginning of a round,
cards will be removed from the deck and dealt to Players, giving them a hand.
### Player
Players are structs that have two important properties: a tcp connection
and a slice of 2 cards representing a hand. The tcp connection is used to manage
messaging between the client and the server. 
### Hand Ranking
Hand ranking is simple. For every time we need to compare player hands, we take a single player,
iterate from the highest possible hand to the lowest possible hand, determine their hand rank, and
finally compare all the player's hand ranks once they have been determined. If players have the
same hand rank, we look to the kicker value of their hand to determine a winner.
### Networking
PokerMUD uses TCP to ensure that no packet loss occurs. Every network message has the following message framing:
```
<MessageLength int>:<MessageType int>:<MessageContent string>
```
MessageLength tells us how far into the buffer to read so that we can succesfully parse our network message. MessageType
tells us how we should handle the message. A standard message, 0, is just text that is broadcasted to players. An
authentication request message, 1, tells us we must authenticate a new connection. An action message, 2, notifies the client
that the player must make a in-game action.