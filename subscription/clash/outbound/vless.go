package outbound

type VlessOption struct {
	Base
	Server            string            `proxy:"server"`
	Port              int               `proxy:"port"`
	UUID              string            `proxy:"uuid"`
	Flow              string            `proxy:"flow,omitempty"`
	TLS               bool              `proxy:"tls,omitempty"`
	ALPN              []string          `proxy:"alpn,omitempty"`
	UDP               bool              `proxy:"udp,omitempty"`
	PacketAddr        bool              `proxy:"packet-addr,omitempty"`
	XUDP              bool              `proxy:"xudp,omitempty"`
	PacketEncoding    string            `proxy:"packet-encoding,omitempty"`
	Network           string            `proxy:"network,omitempty"`
	RealityOpts       RealityOptions    `proxy:"reality-opts,omitempty"`
	HTTPOpts          HTTPOptions       `proxy:"http-opts,omitempty"`
	HTTP2Opts         HTTP2Options      `proxy:"h2-opts,omitempty"`
	GrpcOpts          GrpcOptions       `proxy:"grpc-opts,omitempty"`
	WSOpts            WSOptions         `proxy:"ws-opts,omitempty"`
	WSPath            string            `proxy:"ws-path,omitempty"`
	WSHeaders         map[string]string `proxy:"ws-headers,omitempty"`
	SkipCertVerify    bool              `proxy:"skip-cert-verify,omitempty"`
	Fingerprint       string            `proxy:"fingerprint,omitempty"`
	ServerName        string            `proxy:"servername,omitempty"`
	ClientFingerprint string            `proxy:"client-fingerprint,omitempty"`
}
