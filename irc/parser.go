package irc

import (
	"strings"
)

func (irc *Connection) parseMessage(line string) *Message {
	line = strings.TrimRight(line, "\r\n")
	message := &Message{Raw: line}

	// Parse IRCv3 tags
	if len(line) > 0 && line[0] == '@' {
		if i := strings.Index(line, " "); i > -1 {
			message.Tags = line[1:i]
			line = line[i+1:]
		} else {
			return nil
		}
	}

	// Parse the message prefix
	if len(line) > 0 && line[0] == ':' {
		if i := strings.Index(line, " "); i > -1 {
			message.Source = line[1:i]
			line = line[i+1:]
		} else {
			return nil
		}
	} else {
		message.Source = irc.ClientConfig.Server
	}

	if len(line) == 0 {
		return nil
	}

	// Parse the command and its arguments
	split := strings.SplitN(line, " :", 2)
	args := strings.Split(split[0], " ")
	message.Verb = strings.ToUpper(args[0])
	params := args[1:]
	if len(split) > 1 {
		params = append(params, split[1])
	}
	message.Params = params
	return message
}
