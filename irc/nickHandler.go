package irc

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type nickHandler struct {
	preferred string
	current   string
	letters   []rune
}

func NewNickHandler(preferredNickname string) *nickHandler {
	return &nickHandler{
		preferred: preferredNickname,
		current:   preferredNickname,
		letters:   []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	}
}

func (h *nickHandler) install(c *Connection) {
	c.AddInboundHandler("432", h.erroneusNickame)
	c.AddInboundHandler("433", h.nicknameInUse)
	c.AddInboundHandler("436", h.nicknameCollision)
	c.AddInboundHandler("NICK", h.nicknameChanged)
}

func (h *nickHandler) nicknameChanged(c *Connection, m *Message) {
	log.Printf("Nickname changed: %s", m.Params[0])
	h.current = m.Params[0]
}

func (h *nickHandler) nicknameCollision(c *Connection, m *Message) {
	h.nicknameInUse(c, m)
}

func (h *nickHandler) nicknameInUse(c *Connection, _ *Message) {
	log.Printf("Nickname in use %s", h.current)
	h.updateNickname(c, fmt.Sprintf("%s%d", h.current, rand.Intn(10)))
}

func (h *nickHandler) erroneusNickame(c *Connection, _ *Message) {
	log.Printf("Erroneous nickname (%s), randomising.", h.current)
	h.updateNickname(c, h.randSeq(8))
}

func (h *nickHandler) updateNickname(c *Connection, newNickname string) {
	log.Printf("Changing nickname: %s", newNickname)
	h.current = newNickname
	c.SendRawf("NICK :%s", h.current)
}

func (h *nickHandler) randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = h.letters[rand.Intn(len(h.letters))]
	}
	return string(b)
}
