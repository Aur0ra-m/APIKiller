package ahttp

import (
	"net/http"
	"testing"
)

func TestCopyRequest(t *testing.T) {
	type args struct {
		src *http.Request
	}
	request, _ := http.NewRequest("GET", "http://127.0.0.1/list", nil)
	tests := []struct {
		name    string
		args    args
		wantDst *http.Request
	}{
		// TODO: Add test cases.
		{
			name:    "13213",
			args:    args{request},
			wantDst: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//request.Body = aio.TransformReadCloser(request.Body)

			client := http.Client{}
			client.Do(request)
			client.Do(request)
			RequestClone(request)
		})
	}
}
