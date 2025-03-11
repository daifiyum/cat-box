package outbound

type VmessOption struct {
	Base
	Server              string         `proxy:"server"`
	Port                int            `proxy:"port"`
	UUID                string         `proxy:"uuid"`
	AlterID             int            `proxy:"alterId"`
	Cipher              string         `proxy:"cipher"`
	UDP                 bool           `proxy:"udp,omitempty"`
	Network             string         `proxy:"network,omitempty"`
	TLS                 bool           `proxy:"tls,omitempty"`
	ALPN                []string       `proxy:"alpn,omitempty"`
	SkipCertVerify      bool           `proxy:"skip-cert-verify,omitempty"`
	Fingerprint         string         `proxy:"fingerprint,omitempty"`
	ServerName          string         `proxy:"servername,omitempty"`
	RealityOpts         RealityOptions `proxy:"reality-opts,omitempty"`
	HTTPOpts            HTTPOptions    `proxy:"http-opts,omitempty"`
	HTTP2Opts           HTTP2Options   `proxy:"h2-opts,omitempty"`
	GrpcOpts            GrpcOptions    `proxy:"grpc-opts,omitempty"`
	WSOpts              WSOptions      `proxy:"ws-opts,omitempty"`
	PacketAddr          bool           `proxy:"packet-addr,omitempty"`
	XUDP                bool           `proxy:"xudp,omitempty"`
	PacketEncoding      string         `proxy:"packet-encoding,omitempty"`
	GlobalPadding       bool           `proxy:"global-padding,omitempty"`
	AuthenticatedLength bool           `proxy:"authenticated-length,omitempty"`
	ClientFingerprint   string         `proxy:"client-fingerprint,omitempty"`
}
