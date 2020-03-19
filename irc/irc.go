package irc

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/juju/ratelimit"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func NewIRC(server string, password string, nickname string, useTLS bool, useSasl bool, saslUser, saslPass string,
	debug bool, floodEnabled bool, floodRate time.Duration, floodCapacity int) *Connection {
	log.Print("Creating new IRC")
	connection := &Connection{
		ClientConfig: ClientConfig{
			Server:   server,
			Password: password,
			Nick:     nickname,
			User:     nickname,
			Realname: nickname,
			UseTLS:   useTLS,
		},
		ConnConfig:      DefaultConnectionConfig,
		SASLAuth:        useSasl,
		SASLUser:        saslUser,
		SASLPass:        saslPass,
		Debug:           debug,
		floodProtection: floodEnabled,
		floodRate:       floodRate,
		floodCapacity:   floodCapacity,
	}
	connection.Init()
	return connection
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
		go irc.runRawHandlers(RawMessage{message: msg, out: false})
		message := irc.parseMesage(msg)
		go irc.runInboundHandlers(message)
	}
}

func (irc *Connection) writeLoop() {
	for {
		select {
		case b, ok := <-irc.writeChan:
			if !ok || b == "" || irc.socket == nil {
				break
			}
			go irc.runRawHandlers(RawMessage{message: b, out: true})
			go irc.runOutboundHandlers(b)
			_, err := irc.limitedWriter.Write([]byte(b))
			if err != nil {
				irc.errorChannel <- err
				break
			}
		}
	}
}

func (irc *Connection) miscLoop() {
	keepaliveTicker := time.NewTicker(irc.ConnConfig.KeepAlive)
	for {
		select {
		case <-keepaliveTicker.C:
			if time.Since(irc.lastMessage) >= irc.ConnConfig.KeepAlive {
				irc.SendRawf("PING %d", time.Now().UnixNano())
			}
		case err := <-irc.errorChannel:
			log.Printf("IRC Error occurred: %s", err.Error())
			irc.Finished <- true
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
	irc.listeners = newEventListeners()
	irc.inboundHandlers = make(map[string][]func(*Connection, *Message))
	irc.writeChan = make(chan string, 10)
	irc.errorChannel = make(chan error, 1)
	irc.quitting = make(chan bool, 1)
	irc.signals = make(chan os.Signal, 1)
	irc.Finished = make(chan bool, 1)
	irc.saslFinishedChan = make(chan bool, 1)
	signal.Notify(irc.signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	irc.outboundHandlers = make([]func(*Connection, string), 0)
	irc.rawHandlers = make([]func(*Connection, RawMessage), 0)

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
	if irc.floodProtection {
		irc.limitedWriter = ratelimit.Writer(irc.socket, ratelimit.NewBucket(time.Duration(irc.floodRate*time.Second), int64(irc.floodCapacity)))
	} else {
		irc.limitedWriter = irc.socket
	}
	if err != nil {
		return err
	}

	go irc.readLoop()
	go irc.miscLoop()
	go irc.writeLoop()
	NewErrorHandler().install(irc)
	NewPingHandler().install(irc)
	NewCapabilityHandler().install(irc)
	NewNickHandler(irc.ClientConfig.Nick).install(irc)
	NewDebugHandler(irc.Debug).install(irc)
	NewSASLHandler(irc.SASLAuth, irc.SASLUser, irc.SASLPass).Install(irc)
	if len(irc.ClientConfig.Password) > 0 {
		irc.SendRawf("PASS %s", irc.ClientConfig.Password)
	}
	irc.SendRawf("NICK %s", irc.ClientConfig.Nick)
	irc.SendRawf("USER %s 0 * :%s", irc.ClientConfig.User, irc.ClientConfig.Realname)

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

func (irc *Connection) ConnectAndWaitWithRetry(maxRetries int) error {
	sigWait := make(chan os.Signal, 1)
	signal.Notify(sigWait, os.Interrupt)
	signal.Notify(sigWait, syscall.SIGTERM)
	retryDelay := 0
	retryCount := -1
	for {
		retryCount++
		err := irc.ConnectAndWait()
		if retryCount > maxRetries {
			return errors.New("maximum retries reached")
		}
		irc.Init()
		retryDelay = retryCount*5 + retryDelay
		if retryDelay > 300 {
			retryDelay = 300
		}
		if err != nil {
			log.Printf("Error connecting, retrying in %d", retryDelay)
		} else {
			return nil
		}
		sleep := time.NewTimer(time.Duration(retryDelay) * time.Second)
		select {
		case <-sleep.C:
		//NOOP
		case <-sigWait:
			break
		}
	}
}
