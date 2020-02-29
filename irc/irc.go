package irc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/greboid/irc/config"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func NewIRC(config *config.Config) *Connection {
	log.Print("Creating new IRC")
	return &Connection{
		ClientConfig: ClientConfig{
			Server:   config.Server,
			Password: config.Password,
			Nick:     config.Nickname,
			User:     config.Nickname,
			Realname: config.Nickname,
			UseTLS:   config.TLS,
		},
		ConnConfig: DefaultConnectionConfig,
	}
}

func (irc *Connection) readLoop() {
	rb := bufio.NewReaderSize(irc.socket, 8192+512)
	for {
		msg, err := rb.ReadString('\n')
		if err != nil {
			irc.errorChannel <- err
			break
		}
		irc.lastMessage = time.Now()
		message := irc.parseMesage(msg)
		go irc.runCallbacks(message)
	}
}

func (irc *Connection) writeLoop() {
	keepaliveTicker := time.NewTicker(irc.ConnConfig.KeepAlive)
	for {
		select {
		case <-keepaliveTicker.C:
			if time.Since(irc.lastMessage) >= irc.ConnConfig.KeepAlive {
				irc.SendRawf("PING %d", time.Now().UnixNano())
			}
		case b, ok := <-irc.writeChan:
			if !ok || b == "" || irc.socket == nil {
				break
			}
			_, err := irc.socket.Write([]byte(b))
			if err != nil {
				irc.errorChannel <- err
				break
			}
		case <-irc.signals:
			go irc.doQuit()
		case <-irc.quitting:
			go irc.doQuit()
		}
	}
}

func (irc *Connection) doQuit() {
	irc.SendRaw("QUIT")
	select {
	case <-time.After(2 * time.Second):
		irc.Finished <- true
	}
}

func (irc *Connection) Quit() {
	irc.quitting <- true
}

func (irc *Connection) SendRaw(line string) {
	if !strings.HasSuffix(line, "\r\n") {
		line = line + "\r\n"
	}
	irc.writeChan <- line
}

func (irc *Connection) SendRawf(formatLine string, args ...interface{}) {
	irc.SendRaw(fmt.Sprintf(formatLine, args...))
}

func (irc *Connection) Init() {
	log.Print("Initialising IRC")
	irc.callbacks = make(map[string][]func(*Connection, *Message))
	irc.capabilityHandler = capabilityHandler{}
	irc.writeChan = make(chan string, 10)
	irc.quitting = make(chan bool, 1)
	irc.signals = make(chan os.Signal, 1)
	irc.Finished = make(chan bool, 1)
	signal.Notify(irc.signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	irc.initialised = true
}

func (irc *Connection) Connect() error {
	log.Printf("Connecting to IRC: %s", irc.ClientConfig.Server)
	var err error
	if irc.ClientConfig.UseTLS {
		dialer := &net.Dialer{Timeout: irc.ConnConfig.Timeout}
		irc.socket, err = tls.DialWithDialer(dialer, "tcp", irc.ClientConfig.Server, nil)
	} else {
		irc.socket, err = net.DialTimeout("tcp", irc.ClientConfig.Server, irc.ConnConfig.Timeout)
	}
	if err != nil {
		return err
	}

	go irc.readLoop()
	go irc.writeLoop()
	irc.capabilityHandler.install(irc)
	irc.nickHandler.install(irc)
	if len(irc.ClientConfig.Password) > 0 {
		irc.SendRawf("PASS %s", irc.ClientConfig.Password)
	}
	irc.SendRaw("CAP LS 302")
	irc.SendRawf("NICK %s", irc.ClientConfig.Nick)
	irc.SendRawf("USER %s 0 * :%s", irc.ClientConfig.User, irc.ClientConfig.Realname)

	irc.AddCallbacks(defaultCallbacks)

	return nil
}

func (irc *Connection) Wait() {
	log.Print("Waiting for IRC to finish")
	<-irc.Finished
	close(irc.writeChan)
	_ = irc.socket.Close()
	log.Print("IRC Finished")
}

func (irc *Connection) ConnectAndWait() error {
	if !irc.initialised {
		irc.Init()
	}
	err := irc.Connect()
	if err != nil {
		return err
	}
	irc.Wait()
	return nil
}
