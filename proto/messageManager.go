package tcp

import "github.com/8bitstout/pokermud"

type Messenger interface {
	BroadcastMessage()
	SendMessage()
}

type MessageManager struct {
	players []*pokermud.Player
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
