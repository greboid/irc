package irc

import (
	"reflect"
	"sync"
	"testing"
)

func Test_capabilityHandler_parseCapabilities(t *testing.T) {
	type fields struct {
		available map[capabilityStruct]bool
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
		want   map[capabilityStruct]bool
	}{
		{
			name: "simple",
			args: args{tokenised: []string{"account-notify"}},
			want: map[capabilityStruct]bool{capabilityStruct{name: "account-notify", values: ""}: true},
		},
		{
			name: "simple + value",
			args: args{tokenised: []string{"sasl=PLAIN,EXTERNAL"}},
			want: map[capabilityStruct]bool{capabilityStruct{name: "sasl", values: "PLAIN,EXTERNAL"}: true},
		},
		{
			name: "domain",
			args: args{tokenised: []string{"draft/chathistory"}},
			want: map[capabilityStruct]bool{capabilityStruct{name: "draft/chathistory", values: ""}: true},
		},
		{
			name: "domain + value",
			args: args{tokenised: []string{"draft/languages=13,en,~bs,~de,~el,~en-AU,~es,~fr-FR,~no,~pl,~pt-BR,~ro,~tr-TR,~zh-CN"}},
			want: map[capabilityStruct]bool{capabilityStruct{name: "draft/languages", values: "13,en,~bs,~de,~el,~en-AU,~es,~fr-FR,~no,~pl,~pt-BR,~ro,~tr-TR,~zh-CN"}: true},
		},
		{
			name: "simple + multiple values",
			args: args{tokenised: []string{"sts=duration=2765100,port=6697"}},
			want: map[capabilityStruct]bool{capabilityStruct{name: "sts", values: "duration=2765100,port=6697"}: true},
		},
		{
			name: "multiple",
			args: args{tokenised: []string{"sasl=PLAIN,EXTERNAL", "server-time", "sts=duration=2765100,port=6697"}},
			want: map[capabilityStruct]bool{capabilityStruct{name: "sasl", values: "PLAIN,EXTERNAL"}: true, capabilityStruct{name: "server-time", values: ""}: true, capabilityStruct{name: "sts", values: "duration=2765100,port=6697"}: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := &capabilityHandler{
				available: tt.fields.available,
				wanted:    tt.fields.wanted,
				acked:     tt.fields.acked,
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
