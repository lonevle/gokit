package net

import (
	"net"
)

// IsPrivateIP 判断给定地址是否为私有/内网 IP 地址
// 包含 RFC 1918 定义的私有地址、回环地址和链路本地地址
// 支持 IPv4 和 IPv6 地址的判断
// 参数:
//   - addr: IP 地址字符串
//
// 返回:
//   - true: 是私有/内网地址
//   - false: 不是私有/内网地址或格式错误
func IsPrivateIP(addr string) bool {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false
	}

	// RFC 1918 私有地址 (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, fc00::/7)
	if ip.IsPrivate() {
		return true
	}

	// 回环地址 (127.0.0.0/8, ::1)
	if ip.IsLoopback() {
		return true
	}

	// 链路本地地址 (169.254.0.0/16, fe80::/10)
	if ip.IsLinkLocalUnicast() {
		return true
	}

	return false
}
