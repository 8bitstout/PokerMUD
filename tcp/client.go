package tcp

import (
	"bufio"
	"fmt"
	"github.com/8bitstout/pokermud"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	port            string
	logError        *log.Logger
	logInfo         *log.Logger
	Username        string
	Connection      net.Conn
	MessageManager  *MessageManager
	canChooseAction bool
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
				switch msgType {
				// Cards are dealt to player
				case 2:
					{
						fmt.Println("Dealing cards...")
						time.Sleep(time.Second * 3)
						fmt.Println("Your hand:", msg)
					}
				// Client connection terminated by server
				case 6:
					{
						c.logInfo.Println("Connection terminated by server")
						c.Connection.Close()
						c.Connection = nil
					}
				// Flop is dealt
				case 7:
					{
						fmt.Println("Dealing the flop...")
						time.Sleep(time.Second * 3)
						fmt.Println(msg)
					}
				// Any plain text sent from the server
				case 10:
					{
						c.logInfo.Println("Received standard message from server")
						fmt.Println(msg)
					}
				// Request player action
				case 11:
					{
						c.logInfo.Println("Server requested action from client")
						fmt.Println("Enter one of the following commands to send an action\n1. fold\n2. check\n3. bet (amount e.g 10)")
						c.canChooseAction = true

						for c.canChooseAction {
							reader := bufio.NewReader(os.Stdin)
							fmt.Print(">> ")
							i, _ := reader.ReadString('\n')
							msg := c.MessageManager.CreateMessage(i, 11)
							c.MessageManager.SendMessage(msg, pokermud.MakePlayer(c.Username, c.Connection))
							reader.Reset(os.Stdin)
							c.canChooseAction = false
						}

					}
				// Player ran out of time to choose an action
				case 12:
					{
						c.logInfo.Println("Player ran out of time to choose action and server terminating player action")
						c.canChooseAction = false
					}
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
