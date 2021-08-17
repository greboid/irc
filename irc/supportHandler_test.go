package irc

import (
	"log"
	"reflect"
	"testing"

	"github.com/ergochat/irc-go/ircmsg"
)

type noopLogger struct{}

func (n noopLogger) Debugf(_ string, _ ...interface{}) {}

func (n noopLogger) Infof(_ string, _ ...interface{}) {}

func (n noopLogger) Warnf(_ string, _ ...interface{}) {}

func (n noopLogger) Errorf(_ string, _ ...interface{}) {}

func (n noopLogger) Panicf(_ string, _ ...interface{}) {}

func TestNewSupportHandler(t *testing.T) {
	value := NewSupportHandler()
	if value.conn != nil {
		t.Error("NewSupportHandler() connection not nil")
	}
	if value.values != nil && len(value.values) != 0 {
		t.Error("NewSupportHandler() values not empty")
	}
}

func Test_supportParser_install(t *testing.T) {
	conn := &Connection{
		logger: noopLogger{},
	}
	value := supportParser{}
	value.install(conn)
	if value.conn != conn {
		t.Error("NewSupportHandler() connection not set")
	}
	// TODO: Test 005 hook adding
}

func Test_supportParser_remove(t *testing.T) {
	type args struct {
		s []supportedValue
		i int
	}
	tests := []struct {
		name string
		args args
		want []supportedValue
	}{
		{
			name: "",
			args: args{
				s: []supportedValue{{name: "1", value: "1"}, {name: "2", value: "1"}, {name: "3", value: "1"}, {name: "4", value: "1"}},
				i: 0,
			},
			want: []supportedValue{{name: "2", value: "1"}, {name: "3", value: "1"}, {name: "4", value: "1"}},
		},
		{
			name: "",
			args: args{
				s: []supportedValue{{name: "1", value: "1"}, {name: "2", value: "1"}, {name: "3", value: "1"}, {name: "4", value: "1"}},
				i: 1,
			},
			want: []supportedValue{{name: "1", value: "1"}, {name: "3", value: "1"}, {name: "4", value: "1"}},
		},
		{
			name: "",
			args: args{
				s: []supportedValue{{name: "1", value: "1"}, {name: "2", value: "1"}, {name: "3", value: "1"}, {name: "4", value: "1"}},
				i: 2,
			},
			want: []supportedValue{{name: "1", value: "1"}, {name: "2", value: "1"}, {name: "4", value: "1"}},
		},
		{
			name: "",
			args: args{
				s: []supportedValue{{name: "1", value: "1"}, {name: "2", value: "1"}, {name: "3", value: "1"}, {name: "4", value: "1"}},
				i: 3,
			},
			want: []supportedValue{{name: "1", value: "1"}, {name: "2", value: "1"}, {name: "3", value: "1"}},
		},
		{
			name: "",
			args: args{
				s: []supportedValue{{name: "1", value: "1"}, {name: "2", value: "1"}, {name: "3", value: "1"}, {name: "4", value: "1"}},
				i: 4,
			},
			want: []supportedValue{{name: "1", value: "1"}, {name: "2", value: "1"}, {name: "3", value: "1"}, {name: "4", value: "1"}},
		},
		{
			name: "",
			args: args{
				s: []supportedValue{{name: "1", value: "1"}, {name: "2", value: "1"}, {name: "3", value: "1"}, {name: "4", value: "4"}},
				i: -1,
			},
			want: []supportedValue{{name: "1", value: "1"}, {name: "2", value: "1"}, {name: "3", value: "1"}, {name: "4", value: "4"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &supportParser{}
			if got := h.remove(tt.args.s, tt.args.i); !supvalequals(got, tt.want) {
				t.Errorf("remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func supvalequals(a, b []supportedValue) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	matches := len(a)
	for i := range a {
		if supvalContains(b, a[i]) {
			matches--
		}
	}
	return matches == 0
}

func supvalContains(a []supportedValue, cap supportedValue) bool {
	for i := range a {
		if a[i] == cap {
			return true
		}
	}
	return false
}

func Test_supportParser_tokenise(t *testing.T) {
	type fields struct {
		values []supportedValue
		conn   *Connection
	}
	type args struct {
		input []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantValues []supportedValue
	}{
		{
			name: "empty input",
			fields: fields{
				values: nil,
				conn:   nil,
			},
			args: args{
				input: []string{},
			},
			wantValues: []supportedValue{},
		},
		{
			name: "broken input",
			fields: fields{
				values: nil,
				conn:   nil,
			},
			args: args{
				input: []string{"12", "3=4", "5=6"},
			},
			wantValues: []supportedValue{
				{name: "3", value: "4"},
				{name: "5", value: "6"},
			},
		},
		{
			name: "normal input",
			fields: fields{
				values: nil,
				conn:   nil,
			},
			args: args{
				input: []string{"1=2", "3=4", "5=6"},
			},
			wantValues: []supportedValue{
				{name: "1", value: "2"},
				{name: "3", value: "4"},
				{name: "5", value: "6"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &supportParser{
				values: tt.fields.values,
				conn:   tt.fields.conn,
			}
			if gotValues := h.tokenise(tt.args.input); !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("tokenise() = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func Test_supportParser_handleSupport1(t *testing.T) {
	connection := Connection{}
	type fields struct {
		values []supportedValue
		conn   *Connection
	}
	type args struct {
		em *EventManager
		c  *Connection
		m  *ircmsg.Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wanted []supportedValue
	}{
		{
			name: "null message",
			fields: fields{
				values: nil,
				conn:   &connection,
			},
			args: args{
				em: nil,
				c:  &connection,
				m:  nil,
			},
			wanted: nil,
		},
		{
			name: "null params",
			fields: fields{
				values: nil,
				conn:   &connection,
			},
			args: args{
				em: nil,
				c:  &connection,
				m: &ircmsg.Message{},
			},
			wanted: nil,
		},
		{
			name: "empty message",
			fields: fields{
				values: nil,
				conn:   &connection,
			},
			args: args{
				em: nil,
				c:  &connection,
				m: &ircmsg.Message{},
			},
			wanted: nil,
		},
		{
			name: "multiple values",
			fields: fields{
				values: nil,
				conn:   &connection,
			},
			args: args{
				em: nil,
				c:  &connection,
				m:  mustParse(":server 005 nick AWAYLEN=500 CHANLIMIT=#:100 :are supported by this server"),
			},
			wanted: []supportedValue{
				{
					name:  "AWAYLEN",
					value: "500",
				},
				{
					name:  "CHANLIMIT",
					value: "#:100",
				},
			},
		},
		{
			name: "empty value",
			fields: fields{
				values: nil,
				conn:   &connection,
			},
			args: args{
				em: nil,
				c:  &connection,
				m:  mustParse(":server 005 nick EXCEPTS= :are supported by this server"),
			},
			wanted: []supportedValue{
				{
					name:  "EXCEPTS",
					value: "",
				},
			},
		},
		{
			name: "handle remove",
			fields: fields{
				values: []supportedValue{
					{
						name:  "AWAYLEN",
						value: "500",
					},
				},
				conn: &connection,
			},
			args: args{
				em: nil,
				c:  &connection,
				m:  mustParse(":server 005 nick -AWAYLEN=500 :are supported by this server"),
			},
			wanted: []supportedValue{},
		},
		{
			name: "handle remove without value",
			fields: fields{
				values: []supportedValue{},
				conn:   &connection,
			},
			args: args{
				em: nil,
				c:  &connection,
				m:  mustParse(":server 005 nick -AWAYLEN=500 :are supported by this server"),
			},
			wanted: []supportedValue{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &supportParser{
				values: tt.fields.values,
				conn:   tt.fields.conn,
			}
			if h.handleSupport(tt.args.em, tt.args.c, tt.args.m); !reflect.DeepEqual(h.values, tt.wanted) {
				t.Errorf("tokenise() = %v, want %v", h.values, tt.wanted)
			}
		})
	}
}

func mustParse(line string) *ircmsg.Message {
	message, err := ircmsg.ParseLine(line)
	if err != nil {
		log.Panicf("Unable to parse line: %s", err)
	}
	return &message
}