package outbound

type TrojanOption struct {
	Base
	Server            string         `proxy:"server"`
	Port              int            `proxy:"port"`
	Password          string         `proxy:"password"`
	ALPN              []string       `proxy:"alpn,omitempty"`
	SNI               string         `proxy:"sni,omitempty"`
	SkipCertVerify    bool           `proxy:"skip-cert-verify,omitempty"`
	Fingerprint       string         `proxy:"fingerprint,omitempty"`
	UDP               bool           `proxy:"udp,omitempty"`
	Network           string         `proxy:"network,omitempty"`
	RealityOpts       RealityOptions `proxy:"reality-opts,omitempty"`
	GrpcOpts          GrpcOptions    `proxy:"grpc-opts,omitempty"`
	WSOpts            WSOptions      `proxy:"ws-opts,omitempty"`
	SSOpts            TrojanSSOption `proxy:"ss-opts,omitempty"`
	ClientFingerprint string         `proxy:"client-fingerprint,omitempty"`
}

// TrojanSSOption from https://github.com/p4gefau1t/trojan-go/blob/v0.10.6/tunnel/shadowsocks/config.go#L5
type TrojanSSOption struct {
	Enabled  bool   `proxy:"enabled,omitempty"`
	Method   string `proxy:"method,omitempty"`
	Password string `proxy:"password,omitempty"`
}
