package irc

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func (irc *IRCConnection) readLoop() {
	rb := bufio.NewReaderSize(irc.socket, 512)
	for {
			msg, err := rb.ReadString('\n')
			if err != nil {
				irc.errorChannel <- err
				break
			}
			irc.lastMessage = time.Now()
			log.Printf("-> %v", msg)
			message := irc.parseMesage(msg)
			go irc.runCallbacks(message)
	}
}

func (irc *IRCConnection) writeLoop() {
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
			log.Printf("<- %v", b)
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

func (irc *IRCConnection) doQuit() {
	irc.SendRaw("QUIT")
	select {
	case <-time.After(2 * time.Second):
		irc.Finished <- true
	}
}

func (irc *IRCConnection) Quit() {
	irc.quitting <- true
}

func (irc *IRCConnection) SendRaw(line string) {
	if !strings.HasSuffix(line, "\r\n") {
		line = line + "\r\n"
	}
	irc.writeChan <- line
}

func (irc *IRCConnection) SendRawf(formatLine string, args ...interface{}) {
	irc.SendRaw(fmt.Sprintf(formatLine, args...))
}

func (irc *IRCConnection) Init() {
	irc.callbacks = make(map[string]map[int]func(*IRCConnection, *Message))
	irc.writeChan = make(chan string, 10)
	irc.quitting = make(chan bool, 1)
	irc.signals = make(chan os.Signal, 1)
	irc.Finished = make(chan bool, 1)
	signal.Notify(irc.signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	irc.initialised = true
}

func (irc *IRCConnection) Connect() error {
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

	if len(irc.ClientConfig.Password) > 0 {
		irc.SendRawf("PASS %s", irc.ClientConfig.Password)
	}
	irc.SendRawf("NICK %s", irc.ClientConfig.Nick)
	irc.SendRawf("USER %s 0.0.0.0 0.0.0.0 :%s", irc.ClientConfig.User, irc.ClientConfig.User)

	irc.AddCallbacks(defaultCallbacks)

	return nil
}

func (irc *IRCConnection) Wait() {
	<-irc.Finished
	close(irc.writeChan)
	_ = irc.socket.Close()
}

func (irc *IRCConnection) ConnectAndWait() error {
	if !irc.initialised {
		return errors.New("Must be initialised first")
	}
	err := irc.Connect()
	if err != nil {
		return err
	}
	irc.Wait()
	return nil
}
