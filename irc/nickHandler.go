package irc

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type nickHandler struct {
	preferred         string
	current           string
	letters           []rune
	checkingPreferred bool
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
	go h.monitorNickname(c)
}

func (h *nickHandler) nicknameChanged(_ *EventManager, c *Connection, m *Message) {
	sourceNick := strings.SplitN(m.Source, "!", 2)[0]
	destNick := m.Params[0]
	if strings.HasPrefix(sourceNick, h.current) {
		log.Printf("Nickname changed: %s", destNick)
		h.current = destNick
	} else if sourceNick == h.preferred {
		log.Printf("Regained preferred nickname: %s", destNick)
		c.SendRawf("NICK %s", h.preferred)
	}
}

func (h *nickHandler) nicknameCollision(em *EventManager, c *Connection, m *Message) {
	h.nicknameInUse(em, c, m)
}

func (h *nickHandler) nicknameInUse(em *EventManager, c *Connection, _ *Message) {
	if !h.checkingPreferred {
		log.Printf("Nickname in use %s", h.current)
		h.updateNickname(em, c, fmt.Sprintf("%s%d", h.current, rand.Intn(10)))
	} else {
		h.checkingPreferred = false
	}
}

func (h *nickHandler) erroneusNickame(em *EventManager, c *Connection, _ *Message) {
	log.Printf("Erroneous nickname (%s), randomising.", h.current)
	h.updateNickname(em, c, h.randSeq(8))
}

func (h *nickHandler) updateNickname(_ *EventManager, c *Connection, newNickname string) {
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

func (h *nickHandler) monitorNickname(c *Connection) {
	ticker := time.NewTicker(60000 * time.Millisecond)
	sigWait := make(chan os.Signal, 1)
	signal.Notify(sigWait, os.Interrupt)
	signal.Notify(sigWait, syscall.SIGTERM)
	for {
		select {
		case <-sigWait:
			return
		case <-ticker.C:
			if h.current != h.preferred {
				h.checkingPreferred = true
				c.SendRawf("NICK %s", h.preferred)
			}
		}
	}
}
