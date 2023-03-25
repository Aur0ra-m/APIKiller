package ahttp

import (
	"bytes"
	"net/http"
	"testing"
)

func TestModifyPostFormParameter(t *testing.T) {
	type args struct {
		req       *http.Request
		paramName string
		newValue  string
	}
	request, _ := http.NewRequest("POST", "https://localhost", bytes.NewBuffer([]byte("test=123")))
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				req:       request,
				paramName: "test",
				newValue:  "hacker",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifyPostFormParam(tt.args.req, tt.args.paramName, tt.args.newValue)
		})
	}
}
