package clash

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	O "github.com/daifiyum/cat-box/subscription/clash/outbound"
	"github.com/daifiyum/cat-box/subscription/clash/structure"
	"github.com/goccy/go-yaml"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/format"
	"github.com/sagernet/sing/common/json/badoption"
	N "github.com/sagernet/sing/common/network"
)

type Config struct {
	Proxies []map[string]any `yaml:"proxies"`
}

func ParseClashSubscription(_ context.Context, content string) ([]option.Outbound, error) {
	var config Config
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		return nil, E.Cause(err, "parse clash config")
	}

	decoder := structure.NewDecoder(structure.Option{TagName: "proxy", WeaklyTypedInput: true})

	var outbounds []option.Outbound
	for _, proxyMapping := range config.Proxies {
		proxyType, ok := proxyMapping["type"].(string)
		if !ok {
			return nil, fmt.Errorf("missing type")
		}

		proxyName, ok := proxyMapping["name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing name")
		}

		var outbound option.Outbound
		outbound.Tag = proxyName
		switch proxyType {
		case "ss":
			ssOption := &O.ShadowSocksOption{}
			err = decoder.Decode(proxyMapping, ssOption)
			if err != nil {
				return nil, err
			}
			outbound.Type = C.TypeShadowsocks
			outbound.Options = &option.ShadowsocksOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     ssOption.Server,
					ServerPort: uint16(ssOption.Port),
				},
				Password:      ssOption.Password,
				Method:        clashShadowsocksCipher(ssOption.Cipher),
				Plugin:        clashPluginName(ssOption.Plugin),
				PluginOptions: clashPluginOptions(ssOption.Plugin, ssOption.PluginOpts),
				Network:       clashNetworks(ssOption.UDP),
			}
		case "vmess":
			vmessOption := &O.VmessOption{}
			err = decoder.Decode(proxyMapping, vmessOption)
			if err != nil {
				return nil, err
			}
			outbound.Type = C.TypeVMess
			outbound.Options = &option.VMessOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     vmessOption.Server,
					ServerPort: uint16(vmessOption.Port),
				},
				UUID:     vmessOption.UUID,
				Security: vmessOption.Cipher,
				AlterId:  vmessOption.AlterID,
				OutboundTLSOptionsContainer: option.OutboundTLSOptionsContainer{
					TLS: &option.OutboundTLSOptions{
						Enabled:    vmessOption.TLS,
						ServerName: vmessOption.ServerName,
						Insecure:   vmessOption.SkipCertVerify,
					},
				},
				Transport: clashTransport(vmessOption.Network, vmessOption.HTTPOpts, vmessOption.HTTP2Opts, vmessOption.GrpcOpts, vmessOption.WSOpts),
				Network:   clashNetworks(vmessOption.UDP),
			}
		case "trojan":
			trojanOption := &O.TrojanOption{}
			err = decoder.Decode(proxyMapping, trojanOption)
			if err != nil {
				return nil, err
			}
			outbound.Type = C.TypeTrojan
			outbound.Options = &option.TrojanOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     trojanOption.Server,
					ServerPort: uint16(trojanOption.Port),
				},
				Password: trojanOption.Password,
				OutboundTLSOptionsContainer: option.OutboundTLSOptionsContainer{
					TLS: &option.OutboundTLSOptions{
						Enabled:    true,
						ALPN:       trojanOption.ALPN,
						ServerName: trojanOption.SNI,
						Insecure:   trojanOption.SkipCertVerify,
					},
				},
				Transport: clashTransport(trojanOption.Network, O.HTTPOptions{}, O.HTTP2Options{}, trojanOption.GrpcOpts, trojanOption.WSOpts),
				Network:   clashNetworks(trojanOption.UDP),
			}
		case "hysteria":
			hysteriaOption := &O.HysteriaOption{}
			err = decoder.Decode(proxyMapping, hysteriaOption)
			if err != nil {
				return nil, err
			}
			outbound.Type = C.TypeHysteria
			outbound.Options = &option.HysteriaOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     hysteriaOption.Server,
					ServerPort: uint16(hysteriaOption.Port),
				},
				Up:                  hyBandwidth(hysteriaOption.Up),
				Down:                hyBandwidth(hysteriaOption.Down),
				Obfs:                hysteriaOption.Obfs,
				AuthString:          hysteriaOption.AuthString,
				ReceiveWindowConn:   uint64(hysteriaOption.ReceiveWindowConn),
				ReceiveWindow:       uint64(hysteriaOption.ReceiveWindow),
				DisableMTUDiscovery: hysteriaOption.DisableMTUDiscovery,
				OutboundTLSOptionsContainer: option.OutboundTLSOptionsContainer{
					TLS: &option.OutboundTLSOptions{
						Enabled:    true,
						ALPN:       hysteriaOption.ALPN,
						ServerName: hysteriaOption.SNI,
						Insecure:   hysteriaOption.SkipCertVerify,
					},
				},
			}
		case "hysteria2":
			hysteria2Option := &O.Hysteria2Option{}
			err = decoder.Decode(proxyMapping, hysteria2Option)
			if err != nil {
				return nil, err
			}
			outbound.Type = C.TypeHysteria2
			outbound.Options = &option.Hysteria2OutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     hysteria2Option.Server,
					ServerPort: uint16(hysteria2Option.Port),
				},
				UpMbps:   hy2Bandwidth(hysteria2Option.Up),
				DownMbps: hy2Bandwidth(hysteria2Option.Down),
				Obfs: &option.Hysteria2Obfs{
					Type:     hysteria2Option.Obfs,
					Password: hysteria2Option.ObfsPassword,
				},
				Password: hysteria2Option.Password,
				OutboundTLSOptionsContainer: option.OutboundTLSOptionsContainer{
					TLS: &option.OutboundTLSOptions{
						Enabled:    true,
						ALPN:       hysteria2Option.ALPN,
						ServerName: hysteria2Option.SNI,
						Insecure:   hysteria2Option.SkipCertVerify,
					},
				},
			}
		case "vless":
			vlessOption := &O.VlessOption{}
			err = decoder.Decode(proxyMapping, vlessOption)
			if err != nil {
				return nil, err
			}
			outbound.Type = C.TypeVLESS
			outbound.Options = &option.VLESSOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     vlessOption.Server,
					ServerPort: uint16(vlessOption.Port),
				},
				UUID:    vlessOption.UUID,
				Flow:    vlessOption.Flow,
				Network: clashNetworks(vlessOption.UDP),
				OutboundTLSOptionsContainer: option.OutboundTLSOptionsContainer{
					TLS: &option.OutboundTLSOptions{
						Enabled:    vlessOption.TLS,
						ALPN:       vlessOption.ALPN,
						ServerName: vlessOption.ServerName,
						Insecure:   vlessOption.SkipCertVerify,
						Reality: &option.OutboundRealityOptions{
							Enabled:   isRealityOptionsValid(vlessOption.RealityOpts),
							PublicKey: vlessOption.RealityOpts.PublicKey,
							ShortID:   vlessOption.RealityOpts.ShortID,
						},
						UTLS: &option.OutboundUTLSOptions{
							Enabled:     vlessOption.ClientFingerprint != "",
							Fingerprint: vlessOption.ClientFingerprint,
						},
					},
				},
				Transport: clashTransport(vlessOption.Network, vlessOption.HTTPOpts, vlessOption.HTTP2Opts, vlessOption.GrpcOpts, vlessOption.WSOpts),
			}
		}

		// Filter unsupported protocols
		if outbound.Options != nil {
			outbounds = append(outbounds, outbound)
		}
	}
	if len(outbounds) > 0 {
		return outbounds, nil
	}
	return nil, E.New("no servers found")
}

func clashShadowsocksCipher(cipher string) string {
	switch cipher {
	case "dummy":
		return "none"
	}
	return cipher
}

func clashNetworks(udpEnabled bool) option.NetworkList {
	if !udpEnabled {
		return N.NetworkTCP
	}
	return ""
}

func clashPluginName(plugin string) string {
	switch plugin {
	case "obfs":
		return "obfs-local"
	}
	return plugin
}

type shadowsocksPluginOptionsBuilder map[string]any

func (o shadowsocksPluginOptionsBuilder) Build() string {
	var opts []string
	for key, value := range o {
		if value == nil {
			continue
		}
		opts = append(opts, format.ToString(key, "=", value))
	}
	return strings.Join(opts, ";")
}

func clashPluginOptions(plugin string, opts map[string]any) string {
	options := make(shadowsocksPluginOptionsBuilder)
	switch plugin {
	case "obfs":
		options["obfs"] = opts["mode"]
		options["obfs-host"] = opts["host"]
	case "v2ray-plugin":
		options["mode"] = opts["mode"]
		options["tls"] = opts["tls"]
		options["host"] = opts["host"]
		options["path"] = opts["path"]
	}
	return options.Build()
}

func clashTransport(network string, httpOpts O.HTTPOptions, h2Opts O.HTTP2Options, grpcOpts O.GrpcOptions, wsOpts O.WSOptions) *option.V2RayTransportOptions {
	switch network {
	case "http":
		var headers map[string]badoption.Listable[string]
		for key, values := range httpOpts.Headers {
			if headers == nil {
				headers = make(map[string]badoption.Listable[string])
			}
			headers[key] = values
		}
		return &option.V2RayTransportOptions{
			Type: C.V2RayTransportTypeHTTP,
			HTTPOptions: option.V2RayHTTPOptions{
				Method:  httpOpts.Method,
				Path:    clashStringList(httpOpts.Path),
				Headers: headers,
			},
		}
	case "h2":
		return &option.V2RayTransportOptions{
			Type: C.V2RayTransportTypeHTTP,
			HTTPOptions: option.V2RayHTTPOptions{
				Path: h2Opts.Path,
				Host: h2Opts.Host,
			},
		}
	case "grpc":
		return &option.V2RayTransportOptions{
			Type: C.V2RayTransportTypeGRPC,
			GRPCOptions: option.V2RayGRPCOptions{
				ServiceName: grpcOpts.GrpcServiceName,
			},
		}
	case "ws":
		var headers map[string]badoption.Listable[string]
		for key, value := range wsOpts.Headers {
			if headers == nil {
				headers = make(map[string]badoption.Listable[string])
			}
			headers[key] = []string{value}
		}
		return &option.V2RayTransportOptions{
			Type: C.V2RayTransportTypeWebsocket,
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path:                wsOpts.Path,
				Headers:             headers,
				MaxEarlyData:        uint32(wsOpts.MaxEarlyData),
				EarlyDataHeaderName: wsOpts.EarlyDataHeaderName,
			},
		}
	default:
		return nil
	}
}

func clashStringList(list []string) string {
	if len(list) > 0 {
		return list[0]
	}
	return ""
}

func hyBandwidth(v string) string {
	for _, r := range v {
		if unicode.IsLetter(r) {
			return v
		}
	}
	return v + " Mbps"
}

func hy2Bandwidth(input string) int {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return 0
	}

	num, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}

	if len(parts) == 1 {
		return int(num)
	}

	unit := parts[1]

	conversion := map[string]float64{
		"bps":  1.0 / 1_000_000,
		"Bps":  8.0 / 1_000_000,
		"Kbps": 1.0 / 1_000,
		"KBps": 8.0 / 1_000,
		"Mbps": 1,
		"MBps": 8,
		"Gbps": 1_000,
		"GBps": 8_000,
		"Tbps": 1_000_000,
		"TBps": 8_000_000,
	}

	if factor, exists := conversion[unit]; exists {
		return int(num * factor)
	}

	return int(num)
}

func isRealityOptionsValid(opts O.RealityOptions) bool {
	return opts.PublicKey != "" && opts.ShortID != ""
}
