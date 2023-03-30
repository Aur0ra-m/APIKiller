package realtime

import (
	"APIKiller/pkg/config"
	"APIKiller/pkg/logger"
	"APIKiller/pkg/origin"
	"crypto/tls"
	"crypto/x509"
	"github.com/elazarl/goproxy"
	"net"
)

type RealTime struct {
}

func (r *RealTime) LoadOriginRequest(cfg *config.OriginConfig, httpItemQueue chan *origin.TransferItem) {
	logger.Info("[Load Request] load request from real time origin")
	conf := cfg.RealTime
	address := conf["address"]
	port := conf["port"]
	if address == "" || port == "" {
		panic("Config error: address or port is empty\n")
	}

	go func() {
		logger.Infof("starting proxy: listen at %s:%s", address, port)
		l, err := net.Listen("tcp", address+":"+port)
		if err != nil {
			panic(err)
		}

	}()
}

func NewRealTimeOrigin() *RealTime {
	logger.Info("[Origin] real-time origin")
	return &RealTime{}
}

// proxyN
//
//	@Description: Get httpItem objects through goproxy project
//	@param httpItemQueue
//	@return *goproxy.ProxyHttpServer
func proxyN(httpItemQueue chan *origin.TransferItem) *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()

}

func setCA(caCert, caKey []byte) error {
	var (
		goproxyCa tls.Certificate
		err       error
	)

	if goproxyCa, err = tls.X509KeyPair(caCert, caKey); err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}

	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}

}
