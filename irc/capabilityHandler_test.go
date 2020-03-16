package irc

import (
	"reflect"
	"sync"
	"testing"
)

func Test_capabilityHandler_parseCapabilities(t *testing.T) {
	type fields struct {
		available map[CapabilityStruct]bool
		wanted    map[string]bool
		acked     map[string]bool
		listing   bool
		requested bool
		finished  bool
		mutex     *sync.Mutex
	}
	type args struct {
		tokenised []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*CapabilityStruct
	}{
		{
			name: "simple",
			args: args{tokenised: []string{"account-notify"}},
			want: map[string]*CapabilityStruct{
				"account-notify": {name: "account-notify", values: ""},
			},
		},
		{
			name: "simple + value",
			args: args{tokenised: []string{"sasl=PLAIN,EXTERNAL"}},
			want: map[string]*CapabilityStruct{
				"sasl": {name: "sasl", values: "PLAIN,EXTERNAL"},
			},
		},
		{
			name: "domain",
			args: args{tokenised: []string{"draft/chathistory"}},
			want: map[string]*CapabilityStruct{
				"draft/chathistory": {name: "draft/chathistory", values: ""},
			},
		},
		{
			name: "domain + value",
			args: args{tokenised: []string{"draft/languages=13,en,~bs,~de,~el,~en-AU,~es,~fr-FR,~no,~pl,~pt-BR,~ro,~tr-TR,~zh-CN"}},
			want: map[string]*CapabilityStruct{
				"draft/languages": {name: "draft/languages", values: "13,en,~bs,~de,~el,~en-AU,~es,~fr-FR,~no,~pl,~pt-BR,~ro,~tr-TR,~zh-CN"},
			},
		},
		{
			name: "simple + multiple values",
			args: args{tokenised: []string{"sts=duration=2765100,port=6697"}},
			want: map[string]*CapabilityStruct{
				"sts": {name: "sts", values: "duration=2765100,port=6697"},
			},
		},
		{
			name: "multiple",
			args: args{tokenised: []string{"sasl=PLAIN,EXTERNAL", "server-time", "sts=duration=2765100,port=6697"}},
			want: map[string]*CapabilityStruct{
				"sasl":        {name: "sasl", values: "PLAIN,EXTERNAL"},
				"server-time": {name: "server-time", values: ""},
				"sts":         {name: "sts", values: "duration=2765100,port=6697"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := &capabilityHandler{
				listing:   tt.fields.listing,
				requested: tt.fields.requested,
				finished:  tt.fields.finished,
				mutex:     tt.fields.mutex,
			}
			if got := ca.parseCapabilities(tt.args.tokenised); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCapabilities() = %v, want %v", got, tt.want)
			}
		})
	}
}
