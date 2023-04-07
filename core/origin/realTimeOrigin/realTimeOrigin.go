package realTimeOrigin

import (
	"APIKiller/core/ahttp"
	"APIKiller/core/aio"
	"APIKiller/core/origin"
	logger "APIKiller/logger"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"net/http"
)

type RealTimeOrigin struct {
	address string
	port    string
}

func (o *RealTimeOrigin) LoadOriginRequest() {
	logger.Infoln("[Load Request] load request from real time origin")
	if o.address == "" || o.port == "" {
		panic("Config error: have not set address or port properly")
	}

	// start to listen
	go func() {
		logger.Infoln(fmt.Sprintf("starting proxy: listen at %s:%s", o.address, o.port))
		l, err := net.Listen("tcp", o.address+":"+o.port)
		if err != nil {
			panic(err)
		}

		proxy := proxyN()
		http.Serve(l, proxy)
	}()

}

func NewRealTimeOrigin() *RealTimeOrigin {
	logger.Infoln("[Origin] real-time origin")

	// get config
	address := viper.GetString("app.origin.realTime.address")
	port := viper.GetString("app.origin.realTime.port")

	return &RealTimeOrigin{
		address: address,
		port:    port,
	}
}

// proxyN
//
//	@Description: Get httpItem objects through goproxy project
//	@param httpItemQueue
//	@return *goproxy.ProxyHttpServer
func proxyN() *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()

	setCA(caCert, caKey) //defined in this file

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			// transform body
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				logger.Errorln("read req body failed, ", err.Error())
				return nil, nil
			}

			bodyStr := string(body)

			req.Body = aio.TransformReadCloser(bytes.NewBuffer([]byte(bodyStr)))

			return req, nil
		},
	)

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		// filter

		// switch aio.ReadCloser to util.repeatReadCloser
		if ctx.Resp.Body != nil {
			ctx.Resp.Body = aio.TransformReadCloser(ctx.Resp.Body)
		}

		// copy request and response
		request := ahttp.RequestClone(ctx.Req)
		response := ahttp.ResponseClone(ctx.Resp, request)

		// transport ctx.Req via channel
		origin.TransferItemQueue <- &origin.TransferItem{
			Req:  request,
			Resp: response,
		}

		return resp
	})

	// handle https connection
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	return proxy
}

func setCA(caCert, caKey []byte) error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}

var caCert = []byte(`-----BEGIN CERTIFICATE-----
MIIF9DCCA9ygAwIBAgIJAODqYUwoVjJkMA0GCSqGSIb3DQEBCwUAMIGOMQswCQYD
VQQGEwJJTDEPMA0GA1UECAwGQ2VudGVyMQwwCgYDVQQHDANMb2QxEDAOBgNVBAoM
B0dvUHJveHkxEDAOBgNVBAsMB0dvUHJveHkxGjAYBgNVBAMMEWdvcHJveHkuZ2l0
aHViLmlvMSAwHgYJKoZIhvcNAQkBFhFlbGF6YXJsQGdtYWlsLmNvbTAeFw0xNzA0
MDUyMDAwMTBaFw0zNzAzMzEyMDAwMTBaMIGOMQswCQYDVQQGEwJJTDEPMA0GA1UE
CAwGQ2VudGVyMQwwCgYDVQQHDANMb2QxEDAOBgNVBAoMB0dvUHJveHkxEDAOBgNV
BAsMB0dvUHJveHkxGjAYBgNVBAMMEWdvcHJveHkuZ2l0aHViLmlvMSAwHgYJKoZI
hvcNAQkBFhFlbGF6YXJsQGdtYWlsLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIP
ADCCAgoCggIBAJ4Qy+H6hhoY1s0QRcvIhxrjSHaO/RbaFj3rwqcnpOgFq07gRdI9
3c0TFKQJHpgv6feLRhEvX/YllFYu4J35lM9ZcYY4qlKFuStcX8Jm8fqpgtmAMBzP
sqtqDi8M9RQGKENzU9IFOnCV7SAeh45scMuI3wz8wrjBcH7zquHkvqUSYZz035t9
V6WTrHyTEvT4w+lFOVN2bA/6DAIxrjBiF6DhoJqnha0SZtDfv77XpwGG3EhA/qoh
hiYrDruYK7zJdESQL44LwzMPupVigqalfv+YHfQjbhT951IVurW2NJgRyBE62dLr
lHYdtT9tCTCrd+KJNMJ+jp9hAjdIu1Br/kifU4F4+4ZLMR9Ueji0GkkPKsYdyMnq
j0p0PogyvP1l4qmboPImMYtaoFuYmMYlebgC9LN10bL91K4+jLt0I1YntEzrqgJo
WsJztYDw543NzSy5W+/cq4XRYgtq1b0RWwuUiswezmMoeyHZ8BQJe2xMjAOllASD
fqa8OK3WABHJpy4zUrnUBiMuPITzD/FuDx4C5IwwlC68gHAZblNqpBZCX0nFCtKj
YOcI2So5HbQ2OC8QF+zGVuduHUSok4hSy2BBfZ1pfvziqBeetWJwFvapGB44nIHh
WKNKvqOxLNIy7e+TGRiWOomrAWM18VSR9LZbBxpJK7PLSzWqYJYTRCZHAgMBAAGj
UzBRMB0GA1UdDgQWBBR4uDD9Y6x7iUoHO+32ioOcw1ICZTAfBgNVHSMEGDAWgBR4
uDD9Y6x7iUoHO+32ioOcw1ICZTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEB
CwUAA4ICAQAaCEupzGGqcdh+L7BzhX7zyd7yzAKUoLxFrxaZY34Xyj3lcx1XoK6F
AqsH2JM25GixgadzhNt92JP7vzoWeHZtLfstrPS638Y1zZi6toy4E49viYjFk5J0
C6ZcFC04VYWWx6z0HwJuAS08tZ37JuFXpJGfXJOjZCQyxse0Lg0tuKLMeXDCk2Y3
Ba0noeuNyHRoWXXPyiUoeApkVCU5gIsyiJSWOjhJ5hpJG06rQNfNYexgKrrraEin
o0jmEMtJMx5TtD83hSnLCnFGBBq5lkE7jgXME1KsbIE3lJZzRX1mQwUK8CJDYxye
i6M/dzSvy0SsPvz8fTAlprXRtWWtJQmxgWENp3Dv+0Pmux/l+ilk7KA4sMXGhsfr
bvTOeWl1/uoFTPYiWR/ww7QEPLq23yDFY04Q7Un0qjIk8ExvaY8lCkXMgc8i7sGY
VfvOYb0zm67EfAQl3TW8Ky5fl5CcxpVCD360Bzi6hwjYixa3qEeBggOixFQBFWft
8wrkKTHpOQXjn4sDPtet8imm9UYEtzWrFX6T9MFYkBR0/yye0FIh9+YPiTA6WB86
NCNwK5Yl6HuvF97CIH5CdgO+5C7KifUtqTOL8pQKbNwy0S3sNYvB+njGvRpR7pKV
BUnFpB/Atptqr4CUlTXrc5IPLAqAfmwk5IKcwy3EXUbruf9Dwz69YA==
-----END CERTIFICATE-----`)

var caKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
