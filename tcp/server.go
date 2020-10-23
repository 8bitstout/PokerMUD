package tcp

import (
	"fmt"
	"github.com/8bitstout/pokermud"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"os"
	"strings"
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

type Messenger interface {
	BroadcastMessage()
	SendMessage()
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

type MessageManager struct {
	players []pokermud.Player
}

func (m *MessageManager) BroadcastMessage(message []byte) {
	for _, p := range m.players {
		go m.SendMessage(message, p)
	}
}

func (m *MessageManager) BroadcastMessageFromPlayer(message []byte, sender pokermud.Player) {
	for _, receiver := range m.players {
		if sender.Connection.LocalAddr() != receiver.Connection.LocalAddr() {
			go m.SendMessage(message, receiver)
		}
	}
}

func (m *MessageManager) SendMessage(message []byte, p pokermud.Player) {
	p.Connection.Write(message)
}

func ParseMessage(messageBuffer []byte) ([]byte, Message) {
	length := int(messageBuffer[0]) + 1
	message := messageBuffer[1:length]
	messageType := Message(messageBuffer[length-1])

	return message[:len(message)-1], messageType
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

		msg, mtype := ParseMessage(data)

		if mtype == MESSAGE_AUTHENTICATE {
			a := &pokermud.Action{}
			err := proto.Unmarshal(msg, a)
			if err != nil {
				s.logError.Println("failed to unmarshal protobuf")
				log.Fatal(err)
			}
			s.players[a.GetPlayerName()] = pokermud.MakePlayer(strings.TrimSuffix(a.GetPlayerName(), "\n"), c)

			c.Write([]byte("Welcome, " + a.GetPlayerName()))
			s.playersReady = true
			s.connections++
			return
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
	go s.SendMessages()
	fmt.Println(s.messages)
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
			s.SendPlayerStackUpdate(game)
			time.Sleep(time.Second * 3)
			game.ForEachPlayer(s.RequestPlayerAction)
			game.DealFlop()
			game.ForEachPlayer(s.SendBoardUpdate)
			time.Sleep(time.Second * 3)
			game.ForEachPlayer(s.DisconnectPlayer)
			break
		}
	}
	s.logInfo.Println("This game has ended")
}

func (s *Server) SendMessages() {
	s.logInfo.Println("Listening for channel messages")
	for {
		s.logInfo.Println("New message")
		msg := <-s.messages
		s.logInfo.Println("Writing channel message to: ", msg.receivers[0].LocalAddr())
		s.logInfo.Println("Message type: ", msg.message[0])
		s.logInfo.Println(string(msg.message))
		r := msg.receivers[0]
		r.Write(msg.message)
	}
	s.logInfo.Println("Closed channel")
}

func (s *Server) CreateStandardMessage(msg string) []byte {
	s.logInfo.Println("Creating standard network message")
	msg += "10"
	length := byte(len(msg))
	buffer := []byte{length}
	buffer = append(buffer, []byte(msg)...)
	s.logInfo.Println(buffer)
	return buffer
}

func (s *Server) BroadcastMessage(g *pokermud.Game, msg string, excludeAddress interface{}) {
	buffer := s.CreateStandardMessage(msg)
	for _, p := range g.Players {
		if p.Connection.LocalAddr() != excludeAddress {
			p.Connection.Write(buffer)
		}
	}
}

func (s *Server) RequestPlayerAction(p *pokermud.Player) {
	s.logInfo.Println("Requesting action from: ", p.Name)
	s.messages <- CMessage{
		receivers: []net.Conn{p.Connection},
		message:   s.CreateStandardMessage("The action is on you..."),
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
				message:   s.CreateStandardMessage(fmt.Sprintf("Waiting for %s to act", p.Name)),
			}
		}
	}
	response := make([]byte, DEFAULT_BUFFER_SIZE)
	length, _ := p.Connection.Read(response)
	msg = string(response[:length])
	s.BroadcastMessage(s.game, msg, p.Connection.LocalAddr())
}

func (s *Server) DisconnectPlayer(p *pokermud.Player) {
	s.logInfo.Println("Disconnecting player:", p.Name)
	buffer := []byte{6}
	buffer = append(buffer, []byte("EOF")...)
	p.Connection.Write(buffer)
	delete(s.players, p.Name)
	p.Connection.Close()
	s.connections--
}

func (s *Server) SendCardsToPlayer(p *pokermud.Player) {
	s.logInfo.Println("Sending card message,", p.Hand.String(), "to:", p.Name)
	buffer := []byte{5}
	buffer = append(buffer, []byte(p.Hand.String())...)
	p.Connection.Write(buffer)
}

func (s *Server) SendPlayerStackUpdate(g *pokermud.Game) {
	stacks := ""
	for _, p := range g.Players {
		s.logInfo.Println(len(p.Name))
		stacks += fmt.Sprintf("%s: $%v | ", p.Name, p.Value)
	}
	s.logInfo.Println("Sending update of player stacks")
	s.logInfo.Println(stacks)
	s.BroadcastMessage(g, stacks, nil)
}

func (s *Server) SendBoardUpdate(p *pokermud.Player) {
	s.logInfo.Println("Send a board update message to: ", p.Name)
	board := ""
	for _, c := range s.game.Board.Cards {
		board += c.Name + " "
	}
	s.logInfo.Println("Board:", board)
	buffer := []byte{7}
	buffer = append(buffer, []byte(board)...)
	p.Connection.Write(buffer)
}

func MakeServer(port string) *Server {
	return &Server{
		port:     port,
		logInfo:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		logError: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime),
		players:  make(map[string]*pokermud.Player),
		messages: make(chan CMessage),
	}
}
