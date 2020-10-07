package client

import (
	"bufio"
	"fmt"
	server2 "github.com/8bitstout/pokermud/server"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	port     string
	logError *log.Logger
	logInfo  *log.Logger
}

func (c *Client) Connect() {
	connection, err := net.Dial("tcp", c.port)
	if err != nil {
		c.logError.Println(err)
	}

	c.logInfo.Println("Connected to server")

	msg, err := bufio.NewReader(connection).ReadString('\n')

	if err != nil {
		c.logError.Println(err)
		return
	}

	fmt.Println(msg)

	for {
		c.logInfo.Println("New Reader")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, err := reader.ReadString('\n')

		if err != nil {
			c.logInfo.Println(err)
		}

		fmt.Fprintf(connection, text+"\n")

		msg, err := bufio.NewReader(connection).ReadString('\n')

		if err != nil {
			c.logInfo.Println(err)
		}
		c.logInfo.Println("Message received: ", msg)
		fmt.Print("->:", msg)
		if strings.TrimSpace(text) == server2.COMMAND_STOP {
			c.logInfo.Println("TCP client exiting")
			return
		}

	}
	fmt.Fprintf(connection, "player_disconnected\n")
}

func MakeClient(port string) *Client {
	return &Client{
		port:     port,
		logInfo:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		logError: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime),
	}
}
