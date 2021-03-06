package tcp

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	AWAITING_LENGTH = iota
	AWAITING_TYPE
	AWAITING_DATA
	STD_TOKEN = ":"
)

type MessageFrame struct {
	state         int
	tempMessage   []byte
	messageLength int
	messageType   int
	logInfo       *log.Logger
}

type MessageHandler interface {
	Serialize() string
	Deserialize() []byte
}

func (m *MessageFrame) Parse(data []byte) (bool, int, string) {
	var completed bool
	var message string
	var currentIdx int

	for !completed {
		if m.state == AWAITING_LENGTH {
			m.logInfo.Println("Parsing Message Length")
			found, consumed, parsedMessage := m.readToToken(data, currentIdx, STD_TOKEN)
			currentIdx += consumed

			if found {
				m.logInfo.Println("Message Length Found: ", parsedMessage, "|", len(parsedMessage))
				m.messageLength = int(parsedMessage[0])
				m.logInfo.Println("Length:", m.messageLength)
				m.state = AWAITING_TYPE
			}
		}

		if m.state == AWAITING_TYPE {
			m.logInfo.Println("Parsing message type")
			found, consumed, parsedMessage := m.readToToken(data, currentIdx, STD_TOKEN)
			currentIdx += consumed

			if found {
				m.messageType = int(parsedMessage[0])
				m.state = AWAITING_DATA
			}
		}

		if m.state == AWAITING_DATA {
			m.logInfo.Println("Parsing message data")
			found, consumed, parsedMessage := m.readToLength(data, currentIdx, m.messageLength)
			currentIdx += consumed

			if found {
				message = parsedMessage
				completed = true
			}
		}

		if currentIdx > len(data) {
			break
		}
	}

	return completed, m.messageType, message

}

func (m *MessageFrame) ParseProto(data []byte) {

}

func (m *MessageFrame) readToToken(msg []byte, offset int, token string) (bool, int, []byte) {
	var completed bool
	var parsedMessage []byte
	var consumed int

	idx := m.findToken(msg, offset, token)

	if idx == -1 {
		m.logInfo.Println("Token was not found")
		m.store(msg[offset:])
		consumed = len(msg) - offset
	} else {
		m.logInfo.Println("Reading token from index: ", idx, "starting at offset:", offset)
		parsedMessage = append(m.getTempMessage(), msg[offset:idx]...)
		m.logInfo.Println("New parsed message:", parsedMessage)
		completed = true
		consumed = idx - offset + 1
	}

	return completed, consumed, parsedMessage
}

func (m *MessageFrame) readToLength(data []byte, offset int, length int) (bool, int, string) {
	var consumed int
	var parsed string
	var completed bool

	m.logInfo.Println("ReadToLength: ", length, " | offset:", offset, "->", data[offset])

	current := m.getTempMessage()
	m.logInfo.Println("Current message: ", current)
	remaining := len(data) - offset
	toParse := (length - len(current)) - 1

	if remaining >= toParse {
		parsed = string(data[offset:toParse])
		m.logInfo.Println("Parsed:", parsed, "|", len(parsed))
		consumed = toParse
		completed = true
	} else {
		m.store(current)
		m.store(data[offset:])
		consumed = len(data) - offset
	}

	return completed, consumed, parsed
}

func (m *MessageFrame) findToken(data []byte, offset int, token string) int {
	m.logInfo.Println("Searching for token at offset:", offset)
	idx := -1
	for i := offset; i < len(data); i++ {
		if string(data[i]) == token {
			m.logInfo.Println("Found token at index:", i)
			return i
		}
	}

	return idx
}

func (m *MessageFrame) store(data []byte) []byte {
	m.tempMessage = append(m.tempMessage, data...)
	return m.tempMessage
}

func (m *MessageFrame) getTempMessage() []byte {
	temp := m.tempMessage
	m.tempMessage = []byte{}
	return temp
}

func MakeMessageFrame() *MessageFrame {
	enableLogging := os.Getenv("ENABLE_LOGGING") == "1"

	m := &MessageFrame{
		logInfo: log.New(os.Stdout, fmt.Sprintf("INFO:MessageFrame:\t"), log.Ldate|log.Ltime),
	}

	if !enableLogging {
		m.logInfo.SetOutput(ioutil.Discard)
	}

	return m
}
