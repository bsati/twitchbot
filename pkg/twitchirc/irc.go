package twitchirc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	stopCharacter string = "\r\n"
	sslConn       string = "irc.chat.twitch.tv:6697"
	defaultConn   string = "irc.chat.twitch.tv:6667"
)

// MessageEvent is used to store data about a received message
// to be used by message handlers in further processing
type MessageEvent struct {
	Author       string
	Channel      string
	Message      string
	ReceivedTime time.Time
}

// memberEvent holds information about a JOIN / PART event
type memberEvent struct {
	User         string
	Channel      string
	ReceivedTime time.Time
}

// JoinEvent wraps memberEvent and is supplied to JOIN handlers
type JoinEvent memberEvent

// PartEvent wraps memberEvent and is supplied to PART handlers
type PartEvent memberEvent

// IRCClient provides the main functionality for working with Twitch's IRC
type IRCClient struct {
	Channels        []string
	channelMutex    *sync.RWMutex
	capabilities    []string
	conn            *ircConn
	messageHandlers []func(irc *IRCClient, me *MessageEvent)
	mhMutex         *sync.RWMutex
	joinHandlers    []func(irc *IRCClient, je *JoinEvent)
	jhMutex         *sync.RWMutex
	partHandlers    []func(irc *IRCClient, pe *PartEvent)
	phMutex         *sync.RWMutex
	ratelimiter     *Ratelimiter
}

type ircConn struct {
	ssl     bool
	tlsConn *tls.Conn
	netConn net.Conn
}

func (conn *ircConn) Read(buf []byte) (int, error) {
	if conn.ssl {
		return conn.tlsConn.Read(buf)
	}
	return conn.netConn.Read(buf)
}

func (conn *ircConn) Write(bytes []byte) (int, error) {
	if conn.ssl {
		return conn.tlsConn.Write(bytes)
	}
	return conn.netConn.Write(bytes)
}

// NewClient creates a new IRCClient with given parameters
func NewClient(ssl bool) *IRCClient {
	rateLimiterInput := make(map[string]RateLimitInput)
	rateLimiterInput["chat"] = RateLimitInput{
		Limit:         150,
		ResetInterval: time.Minute,
	}
	client := &IRCClient{
		Channels:     make([]string, 5),
		channelMutex: &sync.RWMutex{},
		conn: &ircConn{
			ssl: ssl,
		},
		capabilities:    make([]string, 0),
		messageHandlers: make([]func(irc *IRCClient, me *MessageEvent), 0),
		mhMutex:         &sync.RWMutex{},
		joinHandlers:    make([]func(irc *IRCClient, je *JoinEvent), 0),
		jhMutex:         &sync.RWMutex{},
		partHandlers:    make([]func(irc *IRCClient, pe *PartEvent), 0),
		phMutex:         &sync.RWMutex{},
		ratelimiter:     NewRatelimiter(rateLimiterInput),
	}
	return client
}

// Connect tries to establish a connection with the selected twitch server
// and starts an infinite loop for message handling
func (irc *IRCClient) Connect(nick string, pass string) error {
	if irc.conn.ssl {
		cfg := &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS10,
		}
		c, err := tls.Dial("tcp", sslConn, cfg)
		if err != nil {
			return err
		}
		irc.conn.tlsConn = c
	} else {
		c, err := net.Dial("tcp", defaultConn)
		if err != nil {
			return err
		}
		irc.conn.netConn = c
	}
	irc.sendString("PASS " + pass)
	irc.sendString("NICK " + nick)
	go irc.mainLoop()
	return nil
}

// Close closes the bot's socket for proper application exiting
func (irc *IRCClient) Close() {
	if irc.conn.ssl {
		irc.conn.tlsConn.Close()
	} else {
		irc.conn.netConn.Close()
	}
}

func (irc *IRCClient) mainLoop() {
	buf := make([]byte, 1024)
	for {
		numRcvd, err := irc.conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("TwitchIRC ERROR - Error receiving bytes from Conn. Error message: " + err.Error())
		}
		if numRcvd > 0 {
			str := strings.Trim(string(buf[:numRcvd]), stopCharacter)
			if str == "PING :tmi.twitch.tv" {
				irc.sendString("PONG :tmi.twitch.tv")
			} else {
				messageSplit := strings.Split(str, " ")
				user := strings.Split(messageSplit[0], "!")[0][1:]
				if messageSplit[1] == "PRIVMSG" {
					messageSplit[3] = messageSplit[3][1:]
					irc.executeMessageHandlers(&MessageEvent{
						Author:       user,
						Channel:      messageSplit[2][1:],
						Message:      strings.Join(messageSplit[3:], " "),
						ReceivedTime: time.Now(),
					})
				} else if messageSplit[1] == "JOIN" {
					irc.executeJoinHandlers(&JoinEvent{
						User:         user,
						Channel:      messageSplit[2][1:],
						ReceivedTime: time.Now(),
					})
				} else if messageSplit[1] == "PART" {
					irc.executePartHandlers(&PartEvent{
						User:         user,
						Channel:      messageSplit[2][1:],
						ReceivedTime: time.Now(),
					})
				}
			}
		}
	}
}

func (irc *IRCClient) executeMessageHandlers(me *MessageEvent) {
	irc.mhMutex.RLock()
	defer irc.mhMutex.RUnlock()
	for _, handler := range irc.messageHandlers {
		go handler(irc, me)
	}
}

func (irc *IRCClient) executeJoinHandlers(je *JoinEvent) {
	irc.jhMutex.RLock()
	defer irc.jhMutex.RUnlock()
	for _, handler := range irc.joinHandlers {
		go handler(irc, je)
	}
}

func (irc *IRCClient) executePartHandlers(pe *PartEvent) {
	irc.phMutex.RLock()
	defer irc.phMutex.RUnlock()
	for _, handler := range irc.partHandlers {
		go handler(irc, pe)
	}
}

func (irc *IRCClient) sendString(input string) (int, error) {
	b, err := irc.ratelimiter.CaptureEndpoint("chat")
	defer b.Release()
	if err != nil {
		return 0, err
	}
	err = b.Increment()
	if err != nil {
		return 0, err
	}
	return irc.conn.Write([]byte(input + stopCharacter))
}

// JoinChannel makes the client join the specified channel
func (irc *IRCClient) JoinChannel(channel string) error {
	irc.channelMutex.Lock()
	defer irc.channelMutex.Unlock()
	_, err := irc.sendString("JOIN #" + channel)
	if err != nil {
		irc.Channels = append(irc.Channels, channel)
	}
	return err
}

// LeaveChannel makes the client leave the specified channel
func (irc *IRCClient) LeaveChannel(channel string) error {
	irc.channelMutex.Lock()
	defer irc.channelMutex.Unlock()
	index := -1
	i := 0
	for index == -1 && i < len(irc.Channels) {
		if irc.Channels[i] == channel {
			index = i
		}
	}
	if index != -1 {
		_, err := irc.sendString("PART #" + channel)
		if err != nil {
			irc.Channels = append(irc.Channels[:i], irc.Channels[i+1:]...)
		}
		return err
	}
	return errors.New("Cannot part from a channel that wasn't joined before")
}

// SendMessage sends the specified message to the specified channel
func (irc *IRCClient) SendMessage(channel string, message string) error {
	_, err := irc.sendString("PRIVMSG #" + channel + " :" + message)
	return err
}

// AddHandler adds the specified handler to the corresponding group (message, join, part)
func (irc *IRCClient) AddHandler(handler interface{}) error {
	if f, ok := handler.(func(*IRCClient, *MessageEvent)); ok {
		irc.mhMutex.Lock()
		defer irc.mhMutex.Unlock()
		irc.messageHandlers = append(irc.messageHandlers, f)
		return nil
	}
	if f, ok := handler.(func(*IRCClient, *JoinEvent)); ok {
		irc.jhMutex.Lock()
		defer irc.jhMutex.Unlock()
		irc.joinHandlers = append(irc.joinHandlers, f)
		return nil
	}
	if f, ok := handler.(func(*IRCClient, *PartEvent)); ok {
		irc.phMutex.Lock()
		defer irc.phMutex.Unlock()
		irc.partHandlers = append(irc.partHandlers, f)
		return nil
	}
	return errors.New("Unknown handler type received")
}

// RequestMembershipCapability requests membership capability (see https://dev.twitch.tv/docs/irc/membership)
func (irc *IRCClient) RequestMembershipCapability() {
	irc.sendString("CAP REQ :twitch.tv/membership")
}
