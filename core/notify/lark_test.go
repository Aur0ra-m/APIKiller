package notify

import (
	"APIKiller/core/data"
	"context"
	"testing"
)

func TestLark_Notify(t *testing.T) {
	type fields struct {
		webhookUrl string
		secret     string
		signature  string
		timestamp  int64
	}
	type args struct {
		ctx  context.Context
		item *data.DataItem
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				webhookUrl: "https://open.feishu.cn/open-apis/bot/v2/hook/f658fb89-d83b-40cb-9f56-375caf862e4d",
				secret:     "8YGvp4SgoP9ozUpO8d38Mh",
				signature:  "",
				timestamp:  0,
			},
			args: args{
				ctx:  nil,
				item: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			l := &Lark{
				webhookUrl: tt.fields.webhookUrl,
				secret:     tt.fields.secret,
				signature:  tt.fields.signature,
				timestamp:  tt.fields.timestamp,
			}
			l.init()
			l.Notify(tt.args.ctx, tt.args.item)
		})
	}
}
