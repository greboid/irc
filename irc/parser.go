package irc

import "strings"

func (irc *Connection) parseMesage(line string) *Message {
	line = strings.TrimRight(line, "\r\n")
	if len(line) < 5 {
		panic("//todo handle this nicely")
	}
	message := &Message{
		Raw: line,
	}
	if line[0] == '@' {
		if i := strings.Index(line, " "); i > -1 {
			message.Tags = line[1:i]
			line = line[i+1:]
		} else {
			panic("Malformed msg from server")
		}
	}
	if line[0] == ':' {
		if i := strings.Index(line, " "); i > -1 {
			message.Source = line[1:i]
			line = line[i+1:]
		} else {
			panic("Malformed msg from server")
		}
	} else {
		message.Source = irc.ClientConfig.Server
	}
	split := strings.SplitN(line, " :", 2)
	args := strings.Split(split[0], " ")
	message.Verb = strings.ToUpper(args[0])
	params := args[1:]
	if len(split) > 1 {
		params = append(params, split[1])
	}
	message.ParamsArray = params
	message.Params = strings.Join(params, " ")
	return message
}
