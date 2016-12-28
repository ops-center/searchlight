package net

import (
	"bytes"
	"net"
	"strings"
)

type IPRange struct {
	from net.IP
	to   net.IP
}

func NewIPRange(from string, to string) IPRange {
	return IPRange{net.ParseIP(from).To4(), net.ParseIP(to).To4()}
}

var privateIPRanges []IPRange = []IPRange{
	NewIPRange("10.0.0.0", "10.255.255.255"),
	NewIPRange("172.16.0.0", "172.31.255.255"),
	NewIPRange("192.168.0.0", "192.168.255.255"),
}

func IsPrivateIP(ip net.IP) bool {
	for _, r := range privateIPRanges {
		if bytes.Compare(ip, r.from) >= 0 && bytes.Compare(ip, r.to) <= 0 {
			return true
		}
	}
	return false
}

func routableIPs() []net.IP {
	result := make([]net.IP, 0, 2)
	ifaces, err := net.Interfaces()
	if err != nil {
		return result
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		if strings.HasPrefix(iface.Name, "docker") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			result = append(result, ip)
		}
	}
	return result
}

func GetInternalIP() string {
	ips := routableIPs()
	for _, ip := range ips {
		if IsPrivateIP(ip) {
			return ip.String()
		}
	}
	return ""
}

func GetExternalIPs() []string {
	ret := make([]string, 0)
	ips := routableIPs()
	for _, ip := range ips {
		if !IsPrivateIP(ip) {
			ret = append(ret, ip.String())
		}
	}
	return ret
}
