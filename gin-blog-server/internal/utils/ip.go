package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"log/slog"
	"net"
	"strings"
	"xojoc.pw/useragent"
)

type ipUtil struct{}

var IP = new(ipUtil)

// GetIpAddress 获取用户发送请求的 IP 地址
// 如果服务器不经过代理, 直接把自己 IP 暴露出去, c.Request.RemoteAddr 就可以直接获取 IP
// 目前流行的架构中, 请求经过服务器前基本会经过代理 (Nginx 最常见), 此时直接获取 IP 拿到的是代理服务器的 IP
func (*ipUtil) GetIpAddress(c *gin.Context) (ipAddress string) {
	// c.ClientIP() 获取的是代理服务器的 IP (Nginx)
	// X-Real-IP: Nginx 服务代理, 本项目明确使用 Nginx 作代理, 因此优先获取这个
	ipAddress = c.Request.Header.Get("X-Real-IP")

	// X-Forwarded-For 经过 HTTP 代理或负载均衡服务器时会添加该项
	// X-Forwarded-For 格式: client1,proxy1,proxy2
	// 一般情况下，第一个 IP 为客户端真实 IP，后面的为经过的代理服务器 IP
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ips := c.Request.Header.Get("X-Forwarded-For") // "ip1,ip2,ip3"
		splitIps := strings.Split(ips, ",")            // ["ip1", "ip2", "ip3"]
		if len(splitIps) > 0 {
			ipAddress = splitIps[0]
		}
	}

	// Pdoxy-Client-IP: Apache 服务代理
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.Header.Get("Proxy-Client-IP")
	}

	// WL-Proxy-Client-IP: Weblogic 服务代理
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.Header.Get("WL-Proxy-Client-IP")
	}

	// RemoteAddr: 发出请求的远程主机的 IP 地址 (经过代理会设置为代理机器的 IP)
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.RemoteAddr
	}

	// 检测到是本机 IP, 读取其局域网 IP 地址
	if strings.HasPrefix(ipAddress, "127.0.0.1") || strings.HasPrefix(ipAddress, "[::1]") {
		ip, err := externalIp()
		if err != nil {
			slog.Error("GetIpAddress, externalIp, err: ", err)
		}
		ipAddress = ip.String()
	}

	if ipAddress != "" && len(ipAddress) > 15 {
		if strings.Index(ipAddress, ",") > 0 {
			ipAddress = ipAddress[:strings.Index(ipAddress, ",")]
		}
	}
	return ipAddress
}

// 获取 IP 来源
// https://github.com/lionsoul2014/ip2region
var vIndex []byte // 缓存 VectorIndex 索引, 减少一次固定的 IO 操作

// GetIpSource 查询给定 IP 地址的地理位置
// 参数:
//
//	ipAddress - 需要查询的 IP 地址（字符串形式）
//
// 返回:
//   - 如果查询成功，返回 IP 地址对应的地理位置信息（如：国家|区域|省份|城市|ISP）
//   - 如果查询失败，返回空字符串 ""
func (*ipUtil) GetIpSource(ipAddress string) string {
	// 设置 IP 数据库的路径
	var dbPath = "../assets/ip2region.xdb" // ip 数据库文件

	// 缓存 VectorIndex 索引，减少每次查询时的 IO 操作
	// 如果 vIndex 为空，则加载数据库的向量索引
	if vIndex == nil {
		var err error
		// 加载向量索引文件，vIndex 是全局缓存，用来加速查询
		vIndex, err = xdb.LoadVectorIndexFromFile(dbPath)
		if err != nil {
			// 如果加载失败，记录错误日志并返回空字符串
			slog.Error(fmt.Sprintf("failed to load vector index from `%s` : %s\n", dbPath, err))
			return ""
		}
	}

	// 使用加载的向量索引创建搜索器（searcher），用于进行 IP 查询
	searcher, err := xdb.NewWithVectorIndex(dbPath, vIndex)
	if err != nil {
		// 如果创建搜索器失败，记录错误日志并返回空字符串
		slog.Error("failed to create searcher with vector index: ", err)
		return ""
	}
	defer searcher.Close()

	// 调用 searcher 的 SearchByStr 方法，通过 IP 地址查询地理位置
	// 返回格式：国家|区域|省份|城市|ISP，若查询结果无数据，后面的字段为 0
	region, err := searcher.SearchByStr(ipAddress)
	if err != nil {
		// 如果查询失败，记录错误日志并返回空字符串
		slog.Error(fmt.Sprintf("failed to search ip(%s): %s\n", ipAddress, err))
		return ""
	}
	// 返回查询到的地理位置字符串
	return region
}

// GetIpSourceSimpleIdle 获取 IP 简易信息, 例如: "江苏省苏州市 电信"
func (i *ipUtil) GetIpSourceSimpleIdle(ipAddress string) string {
	region := i.GetIpSource(ipAddress) // 国家|区域|省份|城市|ISP

	// 检测到是内网, 直接返回 "内网IP"
	// 0|0|0|内网IP|内网IP
	if strings.Contains(region, "内网IP") {
		return "内网IP"
	}

	// 一般无法获取到区域
	// 中国|0|江苏省|苏州市|电信
	ipSource := strings.Split(region, "|")
	if ipSource[0] != "中国" && ipSource[0] != "0" {
		return ipSource[0]
	}
	if ipSource[2] == "0" {
		ipSource[2] = ""
	}
	if ipSource[3] == "0" {
		ipSource[3] = ""
	}
	if ipSource[4] == "0" {
		ipSource[4] = ""
	}
	if ipSource[2] == "" && ipSource[3] == "" && ipSource[4] == "" {
		return ipSource[0]
	}
	return ipSource[2] + ipSource[3] + " " + ipSource[4]
}

func (*ipUtil) GetUserAgent(c *gin.Context) *useragent.UserAgent {
	return useragent.Parse(c.Request.UserAgent())
}

// externalIp 非 127.0.0.1 的局域网 IP
// 返回:
//   - net.IP：返回服务器的局域网 IP 地址（如：192.168.x.x）
//   - error：如果获取失败，返回相应的错误信息
func externalIp() (net.IP, error) {
	// 获取服务器的所有网络接口列表
	ifaces, err := net.Interfaces()
	if err != nil {
		// 如果获取网络接口失败，返回错误
		return nil, err
	}

	// 遍历每个网络接口
	for _, iface := range ifaces {
		// 如果网络接口不在活动状态（即接口未启用），跳过
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// 跳过环回接口（如：127.0.0.1）
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的所有地址（包括 IPv4、IPv6）
		addrs, err := iface.Addrs()
		if err != nil {
			// 如果获取地址失败，返回错误
			return nil, err
		}

		// 遍历接口上的所有地址
		for _, addr := range addrs {
			// 获取地址中的 IP 地址
			ip := getIpFromAddr(addr)
			if ip == nil {
				// 如果没有有效的 IP 地址，继续检查下一个地址
				continue
			}
			// 如果找到了有效的 IP 地址，返回该 IP 地址
			return ip, nil
		}
	}

	// 如果没有找到有效的 IP 地址，返回错误
	return nil, errors.New("connected to the network")
}

// getIpFromAddr 从 net.Addr 地址中提取出有效的 IP 地址
// 参数:
//
//	addr - 网络接口的地址（可能是 *net.IPNet 或 *net.IPAddr 类型）
//
// 返回:
//   - net.IP：提取出的 IP 地址
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	// 根据 addr 的类型进行类型断言，提取 IP 地址
	switch v := addr.(type) {
	case *net.IPNet:
		// 如果是 *net.IPNet 类型，获取其中的 IP 地址
		ip = v.IP
	case *net.IPAddr:
		// 如果是 *net.IPAddr 类型，获取其中的 IP 地址
		ip = v.IP
	}

	// 如果没有 IP 地址，或者是环回地址（127.0.0.1），返回 nil
	if ip == nil || ip.IsLoopback() {
		return nil
	}

	// 如果是 IPv6 地址，则转换为 IPv4 地址
	ip = ip.To4()
	// 如果不是 IPv4 地址，返回 nil
	if ip == nil {
		return nil
	}

	// 返回有效的 IPv4 地址
	return ip
}
