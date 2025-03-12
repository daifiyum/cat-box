package outbound

type HysteriaOption struct {
	Base
	Server              string   `proxy:"server"`
	Port                int      `proxy:"port,omitempty"`
	Ports               string   `proxy:"ports,omitempty"`
	Protocol            string   `proxy:"protocol,omitempty"`
	ObfsProtocol        string   `proxy:"obfs-protocol,omitempty"` // compatible with Stash
	Up                  string   `proxy:"up"`
	UpSpeed             int      `proxy:"up-speed,omitempty"` // compatible with Stash
	Down                string   `proxy:"down"`
	DownSpeed           int      `proxy:"down-speed,omitempty"` // compatible with Stash
	Auth                string   `proxy:"auth,omitempty"`
	AuthString          string   `proxy:"auth-str,omitempty"`
	Obfs                string   `proxy:"obfs,omitempty"`
	SNI                 string   `proxy:"sni,omitempty"`
	SkipCertVerify      bool     `proxy:"skip-cert-verify,omitempty"`
	Fingerprint         string   `proxy:"fingerprint,omitempty"`
	ALPN                []string `proxy:"alpn,omitempty"`
	CustomCA            string   `proxy:"ca,omitempty"`
	CustomCAString      string   `proxy:"ca-str,omitempty"`
	ReceiveWindowConn   int      `proxy:"recv-window-conn,omitempty"`
	ReceiveWindow       int      `proxy:"recv-window,omitempty"`
	DisableMTUDiscovery bool     `proxy:"disable-mtu-discovery,omitempty"`
	FastOpen            bool     `proxy:"fast-open,omitempty"`
	HopInterval         int      `proxy:"hop-interval,omitempty"`
}
