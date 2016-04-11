package main

import "net"

type IPExpander struct {
	IP    net.IP
	IPNet *net.IPNet
}

func NewIPExpander(cidr string) (*IPExpander, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	ip := make([]byte, 4)
	copy(ip, ipnet.IP)
	dec(ip)
	return &IPExpander{
		IP:    ip,
		IPNet: ipnet,
	}, nil
}

func (i *IPExpander) Next() (net.IP, bool) {
	inc(i.IP)
	if !i.IPNet.Contains(i.IP) {
		return nil, false
	}
	return i.IP, true
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func dec(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]--
		if ip[j] > 0 {
			break
		}
	}
}
