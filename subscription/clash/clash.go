package clash

import (
	"context"
	"fmt"
	"strings"

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
		}
		outbounds = append(outbounds, outbound)
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
