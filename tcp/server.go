package tcp

import (
	"fmt"
	"github.com/8bitstout/pokermud"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"os"
	"time"
)

const (
	COMMAND_STOP                 = "stop"
	COMMAND_PLAYERS              = "players"
	DEFAULT_BUFFER_SIZE          = 4096
	MESSAGE_ACTION       Message = 0
	MESSAGE_AUTHENTICATE Message = 1
)

type Message int

func (m Message) MakeMessage() proto.Message {
	return [...]proto.Message{
		&pokermud.Action{},
		&pokermud.Action{},
	}[m]
}

type CMessage struct {
	receivers []net.Conn
	message   []byte
}

type Server struct {
	port           string
	logInfo        *log.Logger
	logError       *log.Logger
	connections    int
	players        map[string]*pokermud.Player
	playersReady   bool
	game           *pokermud.Game
	messages       chan CMessage
	messageManager *MessageManager
}

func (s *Server) handleConnection(c net.Conn) {
	s.logInfo.Println("New connection")

	data := make([]byte, DEFAULT_BUFFER_SIZE)
	for {
		_, err := c.Read(data)

		if err != nil {
			s.logError.Println(err)
			return
		}

		msgFrame := MakeMessageFrame()

		completed, msgType, msg := msgFrame.Parse(data)

		if completed {
			if Message(msgType) == MESSAGE_AUTHENTICATE {
				s.logInfo.Println("Authenticating new user: ", msg)
				if _, ok := s.players[msg]; ok {
					s.logInfo.Println("Player, ", msg, "already exists. Terminating new connection")
					m := s.messageManager.CreateStandardMessage("Username already taken! Disconnecting")
					c.Write(m)
					return
				}
				p := pokermud.MakePlayer(msg, c)
				s.players[msg] = p
				s.messageManager.AddPlayer(p)
				s.messageManager.BroadcastMessageFromPlayer(
					s.messageManager.CreateStandardMessage(fmt.Sprint(p.Name, "has joined the game!")),
					p)
				s.messageManager.SendMessage(
					s.messageManager.CreateStandardMessage(
						fmt.Sprint("Welcome to the table, ", p.Name)),
					p)
				s.playersReady = true
				s.connections++
			}
			if msgType == 10 {
				s.logInfo.Println("Server received standard message: ", msg)
			}
		}
		s.logInfo.Println("Parsed Message")
		fmt.Println(msg)
	}
	s.logInfo.Println("Exiting TCP tcp")
	c.Write([]byte("EOF"))
	c.Close()
}

func (s *Server) Start() {
	l, err := net.Listen("tcp", s.port)
	if err != nil {
		s.logError.Println(err)
		return
	}
	s.logInfo.Println("Server running at: ", s.port)

	defer l.Close()

	go s.StartGame()

	for {
		c, err := l.Accept()
		if err != nil {
			s.logError.Println(err)
			return
		}
		s.playersReady = false
		go s.handleConnection(c)
	}
}

func (s *Server) StartGame() {
	s.logInfo.Println("Waiting for players to connect...")
	for {
		if s.connections > 1 && s.playersReady {
			s.logInfo.Println("A game has been started")
			var players []*pokermud.Player
			for _, p := range s.players {
				players = append(players, p)
			}
			game := pokermud.MakeGame(players)
			s.game = game
			game.DealPlayers()
			game.ForEachPlayer(s.SendCardsToPlayer)
			s.BroadcastPlayerStacks(game)
			time.Sleep(time.Second * 3)
			//game.ForEachPlayer(s.RequestPlayerAction)
			game.DealFlop()
			s.BroadcastCommunityCards()
			time.Sleep(time.Second * 3)
			game.ForEachPlayer(s.DisconnectPlayer)
			break
		}
	}
	s.logInfo.Println("This game has ended")
}

func (s *Server) RequestPlayerAction(p *pokermud.Player) {
	s.logInfo.Println("Requesting action from: ", p.Name)
	s.messages <- CMessage{
		receivers: []net.Conn{p.Connection},
		message:   s.messageManager.CreateStandardMessage("The action is on you..."),
	}
	var buffer []byte
	msg := "client action0"
	buffer = append(buffer, byte(len(msg)))
	buffer = append(buffer, []byte(msg)...)
	fmt.Println("ACTION BUFF SIZE")
	fmt.Println(len(buffer))
	s.messages <- CMessage{
		receivers: []net.Conn{p.Connection},
		message:   buffer,
	}
	fmt.Println(s.messages)
	for _, player := range s.game.Players {
		if p.Connection.LocalAddr() != player.Connection.LocalAddr() {
			s.messages <- CMessage{
				receivers: []net.Conn{player.Connection},
				message:   s.messageManager.CreateStandardMessage(fmt.Sprintf("Waiting for %s to act", p.Name)),
			}
		}
	}
	response := make([]byte, DEFAULT_BUFFER_SIZE)
	length, _ := p.Connection.Read(response)
	msg = string(response[:length])
}

// DisconnectPlayer takes a Player and ends their connection and
// reduces the number of player connections in the Server property
func (s *Server) DisconnectPlayer(p *pokermud.Player) {
	s.logInfo.Println("Disconnecting player:", p.Name)
	s.messageManager.SendMessage(s.messageManager.CreateDisconnectMessage(), p)
	delete(s.players, p.Name)
	p.Connection.Close()
	s.connections--
}

// SendCardsToPlayer takes a Player and sends that player a network message
// containing two cards that cah be parsed to use as a Hand
func (s *Server) SendCardsToPlayer(p *pokermud.Player) {
	s.logInfo.Println("Sending card message,", p.Hand.String(), "to:", p.Name)
	msg := s.messageManager.CreateHandMessage(p.Hand)
	s.messageManager.SendMessage(msg, p)
}

// BroadcastPlayerStack sends a network message to every connection
// which contains every player name and their respective chip stack value
func (s *Server) BroadcastPlayerStacks(g *pokermud.Game) {
	stacks := ""
	for _, p := range g.Players {
		s.logInfo.Println(len(p.Name))
		stacks += fmt.Sprintf("%s: $%v | ", p.Name, p.Value)
	}
	s.logInfo.Println("Sending update of player stacks")
	s.logInfo.Println(stacks)
	msg := s.messageManager.CreateStandardMessage(stacks)
	s.messageManager.BroadcastMessage(msg)
}

// BroadcastCommunityCards sends a network message to every connection
// which contains the current community cards in play
func (s *Server) BroadcastCommunityCards() {
	s.logInfo.Println("Sending a board update to all clients")
	board := ""
	for _, c := range s.game.Board.Cards {
		board += c.Name + " "
	}
	s.logInfo.Println("Board:", board)
	msg := s.messageManager.CreateStandardMessage(board)
	s.messageManager.BroadcastMessage(msg)
}

func MakeServer(port string) *Server {
	return &Server{
		port:           port,
		logInfo:        log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		logError:       log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime),
		players:        make(map[string]*pokermud.Player),
		messages:       make(chan CMessage),
		messageManager: MakeMessageManager(),
	}
}
