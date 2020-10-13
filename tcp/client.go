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
	for c.Connection != nil {
		buffer := make([]byte, DEFAULT_BUFFER_SIZE)
		length, err := c.Connection.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		if length == 0 {
			c.logInfo.Println("Breaking loop")
			break
		}
		if length > 0 {
			msg, msgType := ParseMessage(buffer, length)
			if msgType == 5 {
				fmt.Println("Dealing cards...")
				time.Sleep(time.Second * 3)
				fmt.Println("Your hand:", string(msg))
			}
			if msgType == 6 {
				c.logInfo.Println("Connection terminated by server")
				c.Connection.Close()
				c.Connection = nil
			}
			if msgType == 7 {
				fmt.Println("Dealing the flop...")
				time.Sleep(time.Second * 3)
				fmt.Println(string(msg))
			}
			if msgType == 10 {
				c.logInfo.Println("Received standard message from server")
				fmt.Println(string(msg))
			}
			if msgType == MESSAGE_ACTION {
				c.logInfo.Println("Action requested from server")
				fmt.Println("Enter one of the following commands to send an action\n1. fold\n2. check\n3. bet (amount e.g 10)")
				reader := bufio.NewReader(os.Stdin)
				fmt.Print(">> ")
				response := ""
				switch cmd, _ := reader.ReadString('\n'); cmd {
				case "fold\n":
					response = fmt.Sprint(c.Username, "folded their hand")
				case "check\n":
					response = fmt.Sprint(c.Username, "checked their hand")
				case "bet\n":
					{
						fmt.Println("Enter a bet size (e.g. 10)")
						fmt.Print(">> ")
						reader.Reset(os.Stdin)
						i, _ := reader.ReadString('\n')
						response = fmt.Sprintf("%s bet $%s", c.Username, i)
					}
				default:
					fmt.Println("Command not recognized")
				}
				os.Stdin.Close()
				c.Connection.Write([]byte(response))
			}
		}

		buffer = make([]byte, DEFAULT_BUFFER_SIZE)
	}
}

func MakeClient(port string) *Client {
	return &Client{
		port:     port,
		logInfo:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		logError: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime),
	}
}
