package async

import "testing"

func TestAsyncCheckEngine_Start(t *testing.T) {
	type fields struct {
		httpAPI string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "test",
			fields: fields{httpAPI: "http://api.ceye.io/v1/records?token=0920449a5ed8b9db7a287a66a6632498&type=http"},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &AsyncCheckEngine{
				httpAPI: tt.fields.httpAPI,
			}
			e.Start()
		})
	}
}
