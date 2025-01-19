package toolbox

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/silentrc/toolbox/ipv6"
	"net"
	"net/http"
	"strconv"
	"strings"
)

var ipv4Search *xdb.Searcher
var ipv6Search *ipv6.Dat

type ipUtils struct {
}

func (u *utils) NewIP() *ipUtils {
	return &ipUtils{}
}

func (i *ipUtils) Init(ipv4Path, ipv6Path string) {
	vIndex, err := xdb.LoadVectorIndexFromFile(ipv4Path)
	if err != nil {
		fmt.Printf("failed to load vector index from `%s`: %s\n", ipv4Path, err)
		return
	}
	searcher, err := xdb.NewWithVectorIndex(ipv4Path, vIndex)
	if err != nil {
		fmt.Printf("failed to create searcher with vector index: %s\n", err)
		return
	}
	ipv4Search = searcher
	ipv6Client, err := ipv6.NewIPv6(ipv6Path)
	if err != nil {
		fmt.Printf(" ipv6.NewIPv6 err: %s\n", err)
		return
	}
	ipv6Search = ipv6Client
}

func (i *ipUtils) Search(ip string) string {
	switch i.IPv4orIPv6(ip) {
	case 4:
		str, err := ipv4Search.SearchByStr(ip)
		if err != nil {
			return "无法解析"
		}
		return str
	case 6:
		res := ipv6Search.Find(ip)
		return fmt.Sprintf("%v-%v", res.Country, res.Area)
	default:
		return "非法地址"
	}
}

func (i *ipUtils) IPv4orIPv6(ip string) int {
	res := net.ParseIP(ip)
	if res != nil && strings.Contains(ip, ".") {
		return 4
	}
	if res != nil && strings.Contains(ip, ":") {
		return 6
	}
	return 0
}

// ClientPublicIP 尽最大努力实现获取客户端公网 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !IsLocalIp(ip) {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" && !IsLocalIp(ip) {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if !IsLocalIp(ip) {
			return ip
		}
	}
	return ""
}

func IsLocalIp(ip string) bool {
	/*
		局域网（intranet）的IP地址范围包括：
		10．0．0．0／8－－10．0．0．0～10．255．255．255（A类）
		172．16．0．0／12－172．16．0．0－172．31．255．255（B类）
		192．168．0．0／16－－192．168．0．0～192．168．255．255（C类）
	*/
	ipAddr := strings.Split(ip, ".")

	if strings.EqualFold(ipAddr[0], "10") {
		return true
	} else if strings.EqualFold(ipAddr[0], "172") {
		addr, _ := strconv.Atoi(ipAddr[1])
		if addr >= 16 && addr < 31 {
			return true
		}
	} else if strings.EqualFold(ipAddr[0], "192") && strings.EqualFold(ipAddr[1], "168") {
		return true
	}
	return false
}

func (h *httpUtils) GetRealIP(r *http.Request) string {
	ip := ClientPublicIP(r)
	if ip == "" {
		ip = ClientIp(r)
	}
	return ip
}

// ClientIp 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIp(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return ip
	}
	return ""
}
