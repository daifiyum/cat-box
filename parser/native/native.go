package native

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

func NewNativeURIParser(content string) ([]option.Outbound, error) {
	var outbounds []option.Outbound

	for _, proxyRaw := range strings.Split(content, "\n") {
		proxyRaw = strings.TrimSpace(proxyRaw)
		schemeIndex := strings.Index(proxyRaw, "://")
		if schemeIndex == -1 {
			continue
		}

		schema := strings.ToLower(proxyRaw[:schemeIndex])

		var parserFunc func(string) (option.Outbound, error)

		switch schema {
		case "ss":
			parserFunc = newSSNativeParser
		case "tuic":
			parserFunc = newTuicNativeParser
		case "vmess":
			parserFunc = newVMessNativeParser
		case "vless":
			parserFunc = newVLESSNativeParser
		case "trojan":
			parserFunc = newTrojanNativeParser
		case "hysteria":
			parserFunc = newHysteriaNativeParser
		case "hysteria2":
			parserFunc = newHysteria2NativeParser
		default:
			continue
		}

		outbound, err := parserFunc(proxyRaw)
		if err != nil {
			// 所有节点如有解析错误，停止解析，并返回nil
			return nil, fmt.Errorf("[%s]: %w", proxyRaw, err)
		}
		outbounds = append(outbounds, outbound)
	}
	return outbounds, nil
}

func newSSNativeParser(content string) (option.Outbound, error) {
	outbound := option.Outbound{
		Type: C.TypeShadowsocks,
	}
	result, err := url.Parse(content)
	if err != nil {
		return outbound, E.New("invalid ss uri")
	}
	outbound.Tag = result.Fragment
	options := option.ShadowsocksOutboundOptions{}
	options.Server = result.Hostname()
	options.ServerPort = StringToUint16(result.Port())
	if password, _ := result.User.Password(); password != "" {
		options.Method = result.User.Username()
		options.Password = password
	} else {
		userAndPassword := Base64Safe(result.User.Username())
		userAndPasswordParts := strings.Split(userAndPassword, ":")
		if len(userAndPasswordParts) != 2 {
			return option.Outbound{}, E.New("bad user info")
		}
		options.Method, options.Password = userAndPasswordParts[0], userAndPasswordParts[1]
	}
	plugin := result.Query().Get("plugin")
	if index := strings.Index(plugin, ";"); index != -1 {
		if strings.Contains(plugin[:index], "obfs") {
			options.Plugin = "obfs-local"
		} else {
			options.Plugin = plugin[:index]
		}
		options.PluginOptions = plugin[index+1:]
	}
	outbound.ShadowsocksOptions = options
	return outbound, nil
}

func newTuicNativeParser(content string) (option.Outbound, error) {
	outbound := option.Outbound{
		Type: C.TypeTUIC,
	}
	result, err := url.Parse(content)
	if err != nil {
		return outbound, E.New("invalid tuic uri")
	}
	outbound.Tag = result.Fragment
	options := option.TUICOutboundOptions{}
	TLSOptions := option.OutboundTLSOptions{
		Enabled: true,
		ECH:     &option.OutboundECHOptions{},
		UTLS:    &option.OutboundUTLSOptions{},
		Reality: &option.OutboundRealityOptions{},
	}
	options.UUID = result.User.Username()
	passwd, _ := result.User.Password()
	options.Password = passwd
	options.Server = result.Hostname()
	TLSOptions.ServerName = result.Hostname()
	options.ServerPort = StringToUint16(result.Port())
	for key := range result.Query() {
		value := result.Query().Get(key)
		switch key {
		case "congestion_control":
			if value != "cubic" {
				options.CongestionControl = value
			}
		case "udp_relay_mode":
			options.UDPRelayMode = value
		case "udp_over_stream":
			if value == "true" || value == "1" {
				options.UDPOverStream = true
			}
		case "zero_rtt_handshake", "reduce_rtt":
			if value == "true" || value == "1" {
				options.ZeroRTTHandshake = true
			}
		case "heartbeat_interval":
			options.Heartbeat = option.Duration(StringToInt64(value))
		case "sni":
			TLSOptions.ServerName = value
		case "insecure", "skip-cert-verify", "allow_insecure":
			if value == "1" || value == "true" {
				TLSOptions.Insecure = true
			}
		case "disable_sni":
			if value == "1" || value == "true" {
				TLSOptions.DisableSNI = true
			}
		case "tfo", "tcp-fast-open", "tcp_fast_open":
			if value == "1" || value == "true" {
				options.TCPFastOpen = true
			}
		case "alpn":
			TLSOptions.ALPN = strings.Split(value, ",")
		}
	}
	if options.UDPOverStream {
		options.UDPRelayMode = ""
	}
	options.TLS = &TLSOptions
	outbound.TUICOptions = options
	return outbound, nil
}

func newVMessNativeParser(content string) (option.Outbound, error) {
	outbound := option.Outbound{
		Type: C.TypeVMess,
	}
	splitedArr := strings.Split(content, "://")
	content = Base64Safe(splitedArr[1])
	var proxyRaw map[string]any
	err := json.Unmarshal([]byte(content), &proxyRaw)
	if err != nil {
		return outbound, E.New("invalid vmess uri")
	}
	proxy := ConvertToStrings(proxyRaw)
	outbound.Type = C.TypeVMess
	options := option.VMessOutboundOptions{}
	TLSOptions := option.OutboundTLSOptions{
		ECH:     &option.OutboundECHOptions{},
		UTLS:    &option.OutboundUTLSOptions{},
		Reality: &option.OutboundRealityOptions{},
	}

	for key, value := range proxy {
		switch key {
		case "ps":
			outbound.Tag = value
		case "add":
			options.Server = value
			TLSOptions.ServerName = value
		case "port":
			options.ServerPort = StringToUint16(value)
		case "id":
			options.UUID = value
		case "scy":
			options.Security = value
		case "aid":
			options.AlterId, _ = strconv.Atoi(value)
		case "packet_encoding":
			options.PacketEncoding = value
		case "xudp":
			if value == "1" || value == "true" {
				options.PacketEncoding = "xudp"
			}
		case "tls":
			if value == "1" || value == "true" || value == "tls" {
				TLSOptions.Enabled = true
			}
		case "insecure", "skip-cert-verify":
			if value == "1" || value == "true" {
				TLSOptions.Insecure = true
			}
		case "fp":
			TLSOptions.UTLS.Enabled = true
			TLSOptions.UTLS.Fingerprint = value
		case "net":
			var Transport *option.V2RayTransportOptions
			Transport = &option.V2RayTransportOptions{
				Type: "",
				WebsocketOptions: option.V2RayWebsocketOptions{
					Headers: map[string]option.Listable[string]{},
				},
				HTTPOptions: option.V2RayHTTPOptions{
					Host:    option.Listable[string]{},
					Headers: map[string]option.Listable[string]{},
				},
				GRPCOptions: option.V2RayGRPCOptions{},
			}
			switch value {
			case "ws":
				Transport.Type = C.V2RayTransportTypeWebsocket
				if host, exists := proxy["host"]; exists && host != "" {
					for _, headerStr := range strings.Split(fmt.Sprint("Host:", host), "\n") {
						key, valueRaw := SplitKeyValueWithColon(headerStr)
						value := []string{}
						for _, item := range strings.Split(valueRaw, ",") {
							value = append(value, TrimBlank(item))
						}
						Transport.WebsocketOptions.Headers[key] = value
					}
				}
				if path, exists := proxy["path"]; exists && path != "" {
					reg := regexp.MustCompile(`^(.*?)(?:\?ed=(\d*))?$`)
					result := reg.FindStringSubmatch(path)
					Transport.WebsocketOptions.Path = result[1]
					if result[2] != "" {
						Transport.WebsocketOptions.EarlyDataHeaderName = "Sec-WebSocket-Protocol"
						Transport.WebsocketOptions.MaxEarlyData = StringToUint32(result[2])
					}
				}
			case "h2":
				Transport.Type = C.V2RayTransportTypeHTTP
				TLSOptions.Enabled = true
				if host, exists := proxy["host"]; exists && host != "" {
					Transport.HTTPOptions.Host = []string{host}
				}
				if path, exists := proxy["path"]; exists && path != "" {
					Transport.HTTPOptions.Path = path
				}
			case "tcp":
				if tType, exists := proxy["type"]; exists {
					if tType == "http" {
						Transport.Type = C.V2RayTransportTypeHTTP
						if method, exists := proxy["method"]; exists {
							Transport.HTTPOptions.Method = method
						}
						if host, exists := proxy["host"]; exists && host != "" {
							Transport.HTTPOptions.Host = []string{host}
						}
						if path, exists := proxy["path"]; exists && path != "" {
							Transport.HTTPOptions.Path = path
						}
						if headers, exists := proxy["headers"]; exists {
							for _, header := range strings.Split(headers, "\n") {
								reg := regexp.MustCompile(`^[ \t]*?(\S+?):[ \t]*?(\S+?)[ \t]*?$`)
								result := reg.FindStringSubmatch(header)
								key := result[1]
								value := []string{}
								for _, item := range strings.Split(result[2], ",") {
									value = append(value, TrimBlank(item))
								}
								Transport.HTTPOptions.Headers[key] = value
							}
						}
					} else {
						Transport = nil
					}
				}
			case "grpc":
				Transport.Type = C.V2RayTransportTypeGRPC
				if host, exists := proxy["host"]; exists && host != "" {
					Transport.GRPCOptions.ServiceName = host
				}
			default:
				Transport = nil
			}
			options.Transport = Transport

		case "tfo", "tcp-fast-open", "tcp_fast_open":
			if value == "1" || value == "true" {
				options.TCPFastOpen = true
			}
		}
	}
	options.TLS = &TLSOptions
	outbound.VMessOptions = options
	return outbound, nil
}

func newVLESSNativeParser(content string) (option.Outbound, error) {
	outbound := option.Outbound{
		Type: C.TypeVLESS,
	}
	result, err := url.Parse(content)
	if err != nil {
		return outbound, E.New("invalid vless uri")
	}
	outbound.Tag = result.Fragment
	options := option.VLESSOutboundOptions{}
	TLSOptions := option.OutboundTLSOptions{
		ECH:     &option.OutboundECHOptions{},
		UTLS:    &option.OutboundUTLSOptions{},
		Reality: &option.OutboundRealityOptions{},
	}
	options.UUID = result.User.Username()
	options.Server = result.Hostname()
	TLSOptions.ServerName = result.Hostname()
	options.ServerPort = StringToUint16(result.Port())
	proxy := map[string]string{}
	for key := range result.Query() {
		value := result.Query().Get(key)
		switch key {
		case "key", "alpn", "seed", "path", "host":
			proxy[key] = DecodeURIComponent(value)
		default:
			proxy[key] = value
		}
	}
	for key, value := range proxy {
		switch key {
		case "type":
			var Transport *option.V2RayTransportOptions
			Transport = &option.V2RayTransportOptions{
				Type: "",
				WebsocketOptions: option.V2RayWebsocketOptions{
					Headers: map[string]option.Listable[string]{},
				},
				HTTPOptions: option.V2RayHTTPOptions{
					Host:    option.Listable[string]{},
					Headers: map[string]option.Listable[string]{},
				},
				GRPCOptions: option.V2RayGRPCOptions{},
			}
			switch value {
			case "kcp":
				return outbound, E.New("unsupported transport type: kcp")
			case "ws":
				Transport.Type = C.V2RayTransportTypeWebsocket
				if host, exists := proxy["host"]; exists && host != "" {
					for _, header := range strings.Split(fmt.Sprint("Host:", host), "\n") {
						reg := regexp.MustCompile(`^[ \t]*?(\S+?):[ \t]*?(\S+?)[ \t]*?$`)
						result := reg.FindStringSubmatch(header)
						key := result[1]
						value := []string{}
						for _, item := range strings.Split(result[2], ",") {
							value = append(value, TrimBlank(item))
						}
						Transport.WebsocketOptions.Headers[key] = value
					}
				}
				if path, exists := proxy["path"]; exists && path != "" {
					reg := regexp.MustCompile(`^(.*?)(?:\?ed=(\d*))?$`)
					result := reg.FindStringSubmatch(path)
					Transport.WebsocketOptions.Path = result[1]
					if result[2] != "" {
						Transport.WebsocketOptions.EarlyDataHeaderName = "Sec-WebSocket-Protocol"
						Transport.WebsocketOptions.MaxEarlyData = StringToUint32(result[2])
					}
				}
			case "http":
				Transport.Type = C.V2RayTransportTypeHTTP
				if host, exists := proxy["host"]; exists && host != "" {
					Transport.HTTPOptions.Host = strings.Split(host, ",")
				}
				if path, exists := proxy["path"]; exists && path != "" {
					Transport.HTTPOptions.Path = path
				}
			case "grpc":
				Transport.Type = C.V2RayTransportTypeGRPC
				if serviceName, exists := proxy["serviceName"]; exists && serviceName != "" {
					Transport.GRPCOptions.ServiceName = serviceName
				}
			default:
				Transport = nil
			}
			options.Transport = Transport
		case "security":
			if value == "tls" {
				TLSOptions.Enabled = true
			} else if value == "reality" {
				TLSOptions.Enabled = true
				TLSOptions.Reality.Enabled = true
			}
		case "insecure", "skip-cert-verify":
			if value == "1" || value == "true" {
				TLSOptions.Insecure = true
			}
		case "serviceName", "sni", "peer":
			TLSOptions.ServerName = value
		case "alpn":
			TLSOptions.ALPN = strings.Split(value, ",")
		case "fp":
			TLSOptions.UTLS.Enabled = true
			TLSOptions.UTLS.Fingerprint = value
		case "flow":
			if value == "xtls-rprx-vision" {
				options.Flow = "xtls-rprx-vision"
			}
		case "pbk":
			TLSOptions.Reality.PublicKey = value
		case "sid":
			TLSOptions.Reality.ShortID = value
		case "tfo", "tcp-fast-open", "tcp_fast_open":
			if value == "1" || value == "true" {
				options.TCPFastOpen = true
			}
		}
	}
	options.TLS = &TLSOptions
	outbound.VLESSOptions = options
	return outbound, nil
}

func newTrojanNativeParser(content string) (option.Outbound, error) {
	outbound := option.Outbound{
		Type: C.TypeTrojan,
	}
	result, err := url.Parse(content)
	if err != nil {
		return outbound, E.New("invalid trojan uri")
	}
	outbound.Tag = result.Fragment
	options := option.TrojanOutboundOptions{}
	TLSOptions := option.OutboundTLSOptions{
		Enabled: true,
		ECH:     &option.OutboundECHOptions{},
		UTLS:    &option.OutboundUTLSOptions{},
		Reality: &option.OutboundRealityOptions{},
	}
	options.Server = result.Hostname()
	TLSOptions.ServerName = result.Hostname()
	options.ServerPort = StringToUint16(result.Port())
	options.Password = result.User.Username()
	proxy := map[string]string{}
	for key := range result.Query() {
		value := result.Query().Get(key)
		proxy[key] = DecodeURIComponent(value)
	}
	for key, value := range proxy {
		switch key {
		case "insecure", "allowInsecure", "skip-cert-verify":
			if value == "1" || value == "true" {
				TLSOptions.Insecure = true
			}
		case "serviceName", "sni", "peer":
			TLSOptions.ServerName = value
		case "alpn":
			TLSOptions.ALPN = strings.Split(value, ",")
		case "fp":
			TLSOptions.UTLS.Enabled = true
			TLSOptions.UTLS.Fingerprint = value
		case "type":
			var Transport *option.V2RayTransportOptions
			Transport = &option.V2RayTransportOptions{
				Type: "",
				WebsocketOptions: option.V2RayWebsocketOptions{
					Headers: map[string]option.Listable[string]{},
				},
				HTTPOptions: option.V2RayHTTPOptions{
					Host:    option.Listable[string]{},
					Headers: map[string]option.Listable[string]{},
				},
				GRPCOptions: option.V2RayGRPCOptions{},
			}
			switch value {
			case "ws":
				Transport.Type = C.V2RayTransportTypeWebsocket
				if host, exists := proxy["host"]; exists && host != "" {
					for _, header := range strings.Split(fmt.Sprint("Host:", host), "\n") {
						reg := regexp.MustCompile(`^[ \t]*?(\S+?):[ \t]*?(\S+?)[ \t]*?$`)
						result := reg.FindStringSubmatch(header)
						key := result[1]
						value := []string{}
						for _, item := range strings.Split(result[2], ",") {
							value = append(value, TrimBlank(item))
						}
						Transport.WebsocketOptions.Headers[key] = value
					}
				}
				if path, exists := proxy["path"]; exists && path != "" {
					reg := regexp.MustCompile(`^(.*?)(?:\?ed=(\d*))?$`)
					result := reg.FindStringSubmatch(path)
					Transport.WebsocketOptions.Path = result[1]
					if result[2] != "" {
						Transport.WebsocketOptions.EarlyDataHeaderName = "Sec-WebSocket-Protocol"
						Transport.WebsocketOptions.MaxEarlyData = StringToUint32(result[2])
					}
				}
			case "grpc":
				Transport.Type = C.V2RayTransportTypeGRPC
				if serviceName, exists := proxy["grpc-service-name"]; exists && serviceName != "" {
					Transport.GRPCOptions.ServiceName = serviceName
				}
			default:
				Transport = nil
			}
			options.Transport = Transport
		case "tfo", "tcp-fast-open", "tcp_fast_open":
			if value == "1" || value == "true" {
				options.TCPFastOpen = true
			}
		}
	}
	options.TLS = &TLSOptions
	outbound.TrojanOptions = options
	return outbound, nil
}

func newHysteriaNativeParser(content string) (option.Outbound, error) {
	outbound := option.Outbound{
		Type: C.TypeHysteria,
	}
	result, err := url.Parse(content)
	if err != nil {
		return outbound, E.New("invalid hysteria uri")
	}
	outbound.Tag = result.Fragment
	options := option.HysteriaOutboundOptions{}
	TLSOptions := option.OutboundTLSOptions{
		Enabled: true,
		ECH:     &option.OutboundECHOptions{},
		UTLS:    &option.OutboundUTLSOptions{},
		Reality: &option.OutboundRealityOptions{},
	}
	options.Server = result.Hostname()
	TLSOptions.ServerName = result.Hostname()
	options.ServerPort = StringToUint16(result.Port())
	for key := range result.Query() {
		value := result.Query().Get(key)
		switch key {
		case "auth":
			options.AuthString = value
		case "peer", "sni":
			TLSOptions.ServerName = value
		case "alpn":
			TLSOptions.ALPN = strings.Split(value, ",")
		case "ca":
			TLSOptions.CertificatePath = value
		case "ca_str":
			TLSOptions.Certificate = strings.Split(value, "\n")
		case "up":
			options.Up = value
		case "up_mbps":
			options.UpMbps, _ = strconv.Atoi(value)
		case "down":
			options.Down = value
		case "down_mbps":
			options.DownMbps, _ = strconv.Atoi(value)
		case "obfs", "obfsParam":
			options.Obfs = value
		case "insecure", "skip-cert-verify":
			if value == "1" || value == "true" {
				TLSOptions.Insecure = true
			}
		case "tfo", "tcp-fast-open", "tcp_fast_open":
			if value == "1" || value == "true" {
				options.TCPFastOpen = true
			}
		}
	}
	options.TLS = &TLSOptions
	outbound.HysteriaOptions = options
	return outbound, nil
}

func newHysteria2NativeParser(content string) (option.Outbound, error) {
	outbound := option.Outbound{
		Type: C.TypeHysteria2,
	}
	result, err := url.Parse(content)
	if err != nil {
		return outbound, E.New("invalid hysteria2 uri")
	}
	outbound.Tag = result.Fragment
	options := option.Hysteria2OutboundOptions{
		Obfs: &option.Hysteria2Obfs{},
	}
	TLSOptions := option.OutboundTLSOptions{
		Enabled: true,
		ECH:     &option.OutboundECHOptions{},
		UTLS:    &option.OutboundUTLSOptions{},
		Reality: &option.OutboundRealityOptions{},
	}
	options.ServerPort = uint16(443)
	options.Server = result.Hostname()
	TLSOptions.ServerName = result.Hostname()
	options.Password = result.User.Username()
	if result.Port() != "" {
		options.ServerPort = StringToUint16(result.Port())
	}
	for key := range result.Query() {
		value := result.Query().Get(key)
		switch key {
		case "up":
			options.UpMbps, _ = strconv.Atoi(value)
		case "down":
			options.DownMbps, _ = strconv.Atoi(value)
		case "obfs":
			if value == "salamander" {
				options.Obfs.Type = "salamander"
			}
		case "obfs-password":
			options.Obfs.Password = value
		case "insecure", "skip-cert-verify":
			if value == "1" || value == "true" {
				TLSOptions.Insecure = true
			}
		}
	}
	options.TLS = &TLSOptions
	outbound.Hysteria2Options = options
	return outbound, nil
}
