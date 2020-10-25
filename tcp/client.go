package tcp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	port           string
	logError       *log.Logger
	logInfo        *log.Logger
	Username       string
	Connection     net.Conn
	MessageManager *MessageManager
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
	msg := c.MessageManager.CreateAuthMessage(username)
	c.logInfo.Println("Auth Message Created: ", msg)
	c.logInfo.Println(string(msg))
	c.Connection.Write(msg)
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
			continue
		}
		if length > 0 {
			messageFrame := MakeMessageFrame()
			completed, msgType, msg := messageFrame.Parse(buffer)

			if completed {
				if msgType == 5 {
					fmt.Println("Dealing cards...")
					time.Sleep(time.Second * 3)
					fmt.Println("Your hand:", msg)
				}
				if msgType == 6 {
					c.logInfo.Println("Connection terminated by server")
					c.Connection.Close()
					c.Connection = nil
				}
				if msgType == 7 {
					fmt.Println("Dealing the flop...")
					time.Sleep(time.Second * 3)
					fmt.Println(msg)
				}
				if msgType == 10 {
					c.logInfo.Println("Received standard message from server")
					fmt.Println(msg)
				}
				if Message(msgType) == MESSAGE_ACTION {
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
		}
	}
}

func MakeClient(port string) *Client {
	enableLogging := os.Getenv("ENABLE_LOGGING") == "1"

	c := &Client{
		port:           port,
		logInfo:        log.New(os.Stdout, "INFO:Client\t", log.Ldate|log.Ltime),
		logError:       log.New(os.Stdout, "ERROR:Client\t", log.Ldate|log.Ltime),
		MessageManager: MakeMessageManager(),
	}

	if !enableLogging {
		c.logInfo.SetOutput(ioutil.Discard)
	}

	return c
}
