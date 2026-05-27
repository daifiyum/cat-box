package outbound

type AnyTLSOption struct {
	Base
	Server                   string     `proxy:"server"`
	Port                     int        `proxy:"port"`
	Password                 string     `proxy:"password"`
	ALPN                     []string   `proxy:"alpn,omitempty"`
	SNI                      string     `proxy:"sni,omitempty"`
	ECHOpts                  ECHOptions `proxy:"ech-opts,omitempty"`
	ClientFingerprint        string     `proxy:"client-fingerprint,omitempty"`
	SkipCertVerify           bool       `proxy:"skip-cert-verify,omitempty"`
	Fingerprint              string     `proxy:"fingerprint,omitempty"`
	Certificate              string     `proxy:"certificate,omitempty"`
	PrivateKey               string     `proxy:"private-key,omitempty"`
	UDP                      bool       `proxy:"udp,omitempty"`
	IdleSessionCheckInterval int        `proxy:"idle-session-check-interval,omitempty"`
	IdleSessionTimeout       int        `proxy:"idle-session-timeout,omitempty"`
	MinIdleSession           int        `proxy:"min-idle-session,omitempty"`
}
