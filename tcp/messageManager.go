package tcp

import (
	"errors"
	"fmt"
	"github.com/8bitstout/pokermud"
	"log"
	"net"
)

type Messenger interface {
	BroadcastMessage()
	SendMessage()
}

type MessageManager struct {
	players map[net.Addr]*pokermud.Player
}

func (m *MessageManager) BroadcastMessage(message []byte) {
	for _, p := range m.players {
		go m.SendMessage(message, p)
	}
}

func (m *MessageManager) BroadcastMessageFromPlayer(message []byte, sender *pokermud.Player) {
	for _, receiver := range m.players {
		if sender.Connection.LocalAddr() != receiver.Connection.LocalAddr() {
			go m.SendMessage(message, receiver)
		}
	}
}

func (m *MessageManager) SendMessage(message []byte, p *pokermud.Player) {
	p.Connection.Write(message)
}

func (m *MessageManager) AddPlayer(p *pokermud.Player) error {
	addr := p.Connection.LocalAddr()
	if _, ok := m.players[addr]; ok {
		return errors.New(fmt.Sprint("Player ", p.Name, " is already in the players list"))
	}

	m.players[addr] = p
	return nil
}

func (m *MessageManager) RemovePlayer(p *pokermud.Player) {
	delete(m.players, p.Connection.LocalAddr())
}

func (m *MessageManager) CreateStandardMessage(msg string) []byte {
	return m.CreateMessage(msg, 10)
}

func (m *MessageManager) CreateAuthMessage(username string) []byte {
	return m.CreateMessage(username, 1)
}

func (m *MessageManager) CreateDisconnectMessage() []byte {
	return m.CreateMessage("Connection terminated by server", 6)
}

func (m *MessageManager) CreateMessage(msg string, msgType int) []byte {
	var buff []byte
	buff = append(buff, byte(len(msg)))
	buff = append(buff, byte(':'))
	buff = append(buff, byte(msgType))
	buff = append(buff, byte(':'))
	buff = append(buff, []byte(msg)...)
	buff[0] = byte(len(buff))
	log.Println("Message:", buff)
	return buff
}

func MakeMessageManager() *MessageManager {
	return &MessageManager{
		players: make(map[net.Addr]*pokermud.Player),
	}
}
