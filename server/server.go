package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	COMMAND_STOP    = "stop"
	COMMAND_PLAYERS = "players"
)

type Server struct {
	port        string
	logInfo     *log.Logger
	logError    *log.Logger
	connections int
	players     map[string]string
}

func (s *Server) authenticateUser(c net.Conn) {
	s.logInfo.Println("Attempting to authenticate new connection")

	for {
		c.Write([]byte("Please enter a username\n"))
		netData, _ := bufio.NewReader(c).ReadString('\n')
		input := strings.TrimSpace(netData)

		if _, ok := s.players[input]; !ok {
			s.players[input] = input
			c.Write([]byte(fmt.Sprintf("Welcome, " + input + "\n")))
			return
		}
		c.Write([]byte(fmt.Sprint("Username taken, please choose a different username\n")))
	}
}

func (s *Server) handleConnection(c net.Conn) {
	s.logInfo.Println("New connection")
	s.authenticateUser(c)

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		s.logInfo.Println("Listening for input")
		if err != nil {
			s.logError.Println(err)
			s.connections--
			return
		}

		s.logInfo.Println("-> ", netData)
		input := strings.TrimSpace(netData)
		input = strings.ToLower(input)
		s.logInfo.Println("User command:", input)

		if input == COMMAND_STOP {
			break
		}

		if input == "player_disconnected" {
			s.logInfo.Println("A player has disconnected")
			s.connections--
		}

		if input == COMMAND_PLAYERS {
			s.logInfo.Println("User requested number of server connections")
			c.Write([]byte(fmt.Sprint("There are", s.connections, "users connected\n")))
		}

		if input == "list players" {
			s.logInfo.Println("User requested to list players")
			response := ""
			for _, player := range s.players {
				response += player + ","
			}
			s.logInfo.Println(response)
			c.Write([]byte(fmt.Sprintf(response + "\n")))
		}

		s.logInfo.Println("Response sent")
	}
	s.logInfo.Println("Exiting TCP server")
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

	for {
		c, err := l.Accept()
		if err != nil {
			s.logError.Println(err)
			return
		}
		go s.handleConnection(c)
		s.connections++
	}
}

func MakeServer(port string) *Server {
	return &Server{
		port:     port,
		logInfo:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		logError: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime),
		players:  make(map[string]string),
	}
}
