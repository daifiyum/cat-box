package outbound

// HTTP
type HTTPOptions struct {
	Method  string              `proxy:"method,omitempty"`
	Path    []string            `proxy:"path,omitempty"`
	Headers map[string][]string `proxy:"headers,omitempty"`
}

// HTTP2
type HTTP2Options struct {
	Host []string `proxy:"host,omitempty"`
	Path string   `proxy:"path,omitempty"`
}

// gRPC
type GrpcOptions struct {
	GrpcServiceName string `proxy:"grpc-service-name,omitempty"`
}

// WebSocket
type WSOptions struct {
	Path                     string            `proxy:"path,omitempty"`
	Headers                  map[string]string `proxy:"headers,omitempty"`
	MaxEarlyData             int               `proxy:"max-early-data,omitempty"`
	EarlyDataHeaderName      string            `proxy:"early-data-header-name,omitempty"`
	V2rayHttpUpgrade         bool              `proxy:"v2ray-http-upgrade,omitempty"`
	V2rayHttpUpgradeFastOpen bool              `proxy:"v2ray-http-upgrade-fast-open,omitempty"`
}
