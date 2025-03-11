package outbound

// Reality
type RealityOptions struct {
	PublicKey   string `proxy:"public-key,omitempty"`
	ShortID     string `proxy:"short-id,omitempty"`
	Fingerprint string `proxy:"fingerprint,omitempty"`
}
