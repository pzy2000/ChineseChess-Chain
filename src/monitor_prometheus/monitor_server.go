package monitor_prometheus

import (
	loggers "chainmaker_web/src/logger"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	log = loggers.GetLogger(loggers.MODULE_WEB)
)

// MonitorServer Server
type MonitorServer struct {
	httpServer *http.Server
	port       int
}

// NewMonitorServer 实力化
func NewMonitorServer(nameSpace string, port int) *MonitorServer {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return &MonitorServer{
		httpServer: &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second, // 设置读取请求头部的超时时间，例如 5 秒
		},
		port: port,
	}
}

// Start 启动
func (s *MonitorServer) Start() error {
	if s.httpServer != nil {
		endPoint := fmt.Sprintf(":%d", s.port)
		conn, err := net.Listen("tcp", endPoint)
		if err != nil {
			return fmt.Errorf("TCP Listen failed, %s", err.Error())
		}
		go func() {
			err = s.httpServer.Serve(conn)
			if err != nil {
				log.Warnf("http server failed , %s", err.Error())
			}
		}()
	}
	return nil
}
