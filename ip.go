package toolbox

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var SearchIP *xdb.Searcher

type ipUtils struct {
}

func (u *utils) NewIP() *ipUtils {
	return &ipUtils{}
}

func (i *ipUtils) Init() {
	var dbPath = "./conf/data/ip2region.xdb"
	f, err := os.Open(dbPath)
	if err != nil {
		fmt.Printf("read file to load content from `%s`: %s\n", dbPath, err)
		return
	}
	// 1、从 dbPath 加载整个 xdb 到内存
	//cBuff, err := LoadContentFromFile(dbPath)
	cBuff, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("read file to load content from `%s`: %s\n", dbPath, err)
		return
	}

	// 2、用全局的 cBuff 创建完全基于内存的查询对象。
	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		fmt.Printf("failed to create searcher with content: %s\n", err)
		return
	}
	defer f.Close()
	SearchIP = searcher
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

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
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
