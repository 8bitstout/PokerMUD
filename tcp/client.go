package tcp

import (
	"bufio"
	"fmt"
	"github.com/8bitstout/pokermud"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	port       string
	logError   *log.Logger
	logInfo    *log.Logger
	Username   string
	Connection net.Conn
}

func (c *Client) IsConnected() bool {
	return c.Connection != nil
}

func (c *Client) Authenticate() {
	if !c.IsConnected() {
		c.Connect()
		return
	}

	fmt.Println("Enter your username:")
	fmt.Print(">> ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	msg := &pokermud.Action{
		PlayerName: username,
	}

	out, _ := proto.Marshal(msg)
	buffer := []byte{1}
	buffer = append(buffer, out[:]...)
	c.Connection.Write(buffer)
	reader.Reset(c.Connection)
	response, _ := reader.ReadString('\n')
	fmt.Println(response)
}

func (c *Client) Connect() {
	connection, err := net.Dial("tcp", c.port)
	if err != nil {
		c.logError.Println(err)
	}

	c.logInfo.Println("Connected to server")
	c.Connection = connection
	c.Authenticate()
	for {
		buffer := make([]byte, DEFAULT_BUFFER_SIZE)
		length, err := c.Connection.Read(buffer)
		if err != nil {
			c.logInfo.Println(err)
		}
		msg, msgType := ParseMessage(buffer, length)
		if msgType == 5 {
			fmt.Println("Dealing cards...")
			time.Sleep(time.Second * 5)
			fmt.Println(fmt.Println("Your hand:", string(msg)))
		}
		if msgType == 6 {
			c.logInfo.Println("Connection terminated by server")
			break
		}
		if msgType == 7 {
			fmt.Println("Dealing the flop...")
			time.Sleep(time.Second * 3)
			fmt.Println(fmt.Println(string(msg)))
		}
	}
	c.Connection = nil
}

func MakeClient(port string) *Client {
	return &Client{
		port:     port,
		logInfo:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		logError: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime),
	}
}
