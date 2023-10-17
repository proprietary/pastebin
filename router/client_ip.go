package router

import (
	"log"
	"net/http"
	"net/netip"
)

func knownCloudflare() []netip.Prefix {
	// See: https://www.cloudflare.com/ips/
	cloudflareIps := [...]string{
		"173.245.48.0/20",
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"141.101.64.0/18",
		"108.162.192.0/18",
		"190.93.240.0/20",
		"188.114.96.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
		"162.158.0.0/15",
		"104.16.0.0/13",
		"104.24.0.0/14",
		"172.64.0.0/13",
		"131.0.72.0/22",
		"2400:cb00::/32",
		"2606:4700::/32",
		"2803:f800::/32",
		"2405:b500::/32",
		"2405:8100::/32",
		"2a06:98c0::/29",
		"2c0f:f248::/32",
	}
	dst := make([]netip.Prefix, len(cloudflareIps))
	for i, ipString := range cloudflareIps {
		dst[i] = netip.MustParsePrefix(ipString)
	}
	return dst
}

var knownCloudflareNetworks []netip.Prefix

func init() {
	knownCloudflareNetworks = knownCloudflare()
}

func getClientIp(req *http.Request) netip.Addr {
	// normal client
	addrPort := netip.MustParseAddrPort(req.RemoteAddr)
	addr := addrPort.Addr()
	// correctly receive client IP if this service is running behind Cloudflare
	cfClientIp := req.Header.Get("CF-Connecting-IP")
	if len(cfClientIp) > 0 {
		var isTrustworthy bool = false
		// make sure this is not spoofed
		if addr.IsPrivate() || addr.IsLoopback() {
			isTrustworthy = true
		} else {
			for _, prefix := range knownCloudflareNetworks {
				if prefix.Contains(addr) {
					isTrustworthy = true
				}
			}
		}
		if isTrustworthy {
			parsedCfClientIp, err := netip.ParseAddr(cfClientIp)
			if err != nil {
				// basically should never happen
				return addr
			}
			return parsedCfClientIp
		} else {
			// client tried to forge a Cloudflare request
			log.Printf(`Client with IP=%q tried to forge a Cloudflare request by setting the header "CF-Connecting-IP"`,
				addr.String())
			return addr
		}
	}
	return addr
}
