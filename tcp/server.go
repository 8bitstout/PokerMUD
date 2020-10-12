package tcp

import (
	"fmt"
	"github.com/8bitstout/pokermud"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"os"
)

const (
	COMMAND_STOP                 = "stop"
	COMMAND_PLAYERS              = "players"
	DEFAULT_BUFFER_SIZE          = 8000
	MESSAGE_ACTION       Message = 0
	MESSAGE_AUTHENTICATE Message = 1
)

type Message int8

func (m Message) MakeMessage() proto.Message {
	return [...]proto.Message{
		&pokermud.Action{},
		&pokermud.Action{},
	}[m]
}

type Server struct {
	port         string
	logInfo      *log.Logger
	logError     *log.Logger
	connections  int
	players      map[string]*pokermud.Player
	playersReady bool
	game         *pokermud.Game
}

func ParseMessage(messageBuffer []byte, length int) ([]byte, Message) {
	messageType := Message(messageBuffer[0])
	protoMessage := messageBuffer[1:length]

	return protoMessage, messageType
}

func (s *Server) handleConnection(c net.Conn) {
	s.logInfo.Println("New connection")
	//player := pokermud.MakePlayer(fmt.Sprint("Player", s.players), c)
	data := make([]byte, DEFAULT_BUFFER_SIZE)
	for {
		length, err := c.Read(data)

		if err != nil {
			s.logError.Println(err)
			return
		}

		msg, mtype := ParseMessage(data, length)

		if mtype == MESSAGE_AUTHENTICATE {
			a := &pokermud.Action{}
			err := proto.Unmarshal(msg, a)
			if err != nil {
				log.Fatal(err)
			}
			s.players[a.GetPlayerName()] = pokermud.MakePlayer(a.GetPlayerName(), c)

			c.Write([]byte("Welcome, " + a.GetPlayerName()))
			s.playersReady = true
		}
		fmt.Println(msg)
		s.logInfo.Println("Parsed Message")
		fmt.Println(msg)
		return
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
	go s.listen(l)
	s.StartGame()
}

func (s *Server) listen(l net.Listener) {
	for {
		c, err := l.Accept()
		if err != nil {
			s.logError.Println(err)
			return
		}
		s.playersReady = false
		go s.handleConnection(c)
		s.connections++
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
			game.DealFlop()
			game.ForEachPlayer(s.SendBoardUpdate)
			game.ForEachPlayer(s.DisconnectPlayer)
			break
		}
	}
}

func (s *Server) DisconnectPlayer(p *pokermud.Player) {
	s.logInfo.Println("Disconnecting player:", p.Name)
	buffer := []byte{6}
	buffer = append(buffer, []byte("EOF")...)
	p.Connection.Write(buffer)
	delete(s.players, p.Name)
}

func (s *Server) SendCardsToPlayer(p *pokermud.Player) {
	s.logInfo.Println("Sending card message,", p.Hand.String(), "to:", p.Name)
	buffer := []byte{5}
	buffer = append(buffer, []byte(p.Hand.String())...)
	p.Connection.Write(buffer)
}

func (s *Server) SendBoardUpdate(p *pokermud.Player) {
	s.logInfo.Println("Send a board update message to: ", p.Name)
	board := ""
	for _, c := range s.game.Board.Cards {
		board += c.Name + " "
	}
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
	}
}
