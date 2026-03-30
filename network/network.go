package network

import (
	"memmole/pkg/object"
	"memmole/pkg/logger"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// NetworkPackage 网络包对象
type NetworkPackage struct {
	connections map[int64]*net.Conn
	nextID      int64
	servers     map[int64]*HTTPServer
	serverID    int64
	serverMutex sync.RWMutex
}

// HTTPServer HTTP服务器结构
type HTTPServer struct {
	server   *http.Server
	handlers map[string]http.HandlerFunc
	mux      *http.ServeMux
	running  bool
	port     int
}

// NewNetworkPackage 创建网络包实例
func NewNetworkPackage() *NetworkPackage {
	return &NetworkPackage{
		connections: make(map[int64]*net.Conn),
		nextID:      1,
		servers:     make(map[int64]*HTTPServer),
		serverID:    1,
	}
}

// isPortAvailable 检查端口是否可用
func (np *NetworkPackage) isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// findAvailablePort 查找可用端口
func (np *NetworkPackage) findAvailablePort(startPort int) int {
	for port := startPort; port < startPort+100; port++ {
		if np.isPortAvailable(port) {
			return port
		}
	}
	return -1 // 没有找到可用端口
}

// CreateHTTPServer 创建HTTP服务器
func (np *NetworkPackage) CreateHTTPServer(port int64) *object.Integer {
	np.serverMutex.Lock()
	defer np.serverMutex.Unlock()

	serverID := np.serverID
	np.serverID++

	// 检查端口是否可用，如果不可用则自动选择其他端口
	actualPort := int(port)
	if !np.isPortAvailable(actualPort) {
		logger.Warn("端口 %d 已被占用，正在查找可用端口...", actualPort)
		newPort := np.findAvailablePort(actualPort + 1)
		if newPort == -1 {
			// 如果找不到可用端口，返回错误
			return &object.Integer{Value: -1}
		}
		actualPort = newPort
		logger.Info("自动切换到端口 %d", actualPort)
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", actualPort),
		Handler: mux,
	}

	httpServer := &HTTPServer{
		server:   server,
		handlers: make(map[string]http.HandlerFunc),
		mux:      mux,
		running:  false,
		port:     actualPort,
	}

	np.servers[serverID] = httpServer

	return &object.Integer{Value: serverID}
}

// AddRoute 添加路由处理器
func (np *NetworkPackage) AddRoute(serverID int64, method, path string, handler func(w http.ResponseWriter, r *http.Request)) *object.Integer {
	np.serverMutex.RLock()
	httpServer, exists := np.servers[serverID]
	np.serverMutex.RUnlock()

	if !exists {
		return &object.Integer{Value: -1}
	}

	routeKey := fmt.Sprintf("%s:%s", strings.ToUpper(method), path)
	httpServer.handlers[routeKey] = handler

	// 注册到ServeMux
	httpServer.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) == strings.ToUpper(method) {
			handler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return &object.Integer{Value: 0}
}

// StartServer 启动HTTP服务器
func (np *NetworkPackage) StartServer(serverID int64) *object.String {
	np.serverMutex.RLock()
	httpServer, exists := np.servers[serverID]
	np.serverMutex.RUnlock()

	if !exists {
		return &object.String{Value: "服务器不存在"}
	}

	if httpServer.running {
		return &object.String{Value: "服务器已在运行"}
	}

	httpServer.running = true

	// 在goroutine中启动服务器
	go func() {
		err := httpServer.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("服务器错误: %v", err)
		}
		httpServer.running = false
	}()

	return &object.String{Value: fmt.Sprintf("服务器已启动在端口 %d", httpServer.port)}
}

// StartServerAndWait 启动HTTP服务器并阻塞等待
func (np *NetworkPackage) StartServerAndWait(serverID int64) *object.String {
	np.serverMutex.RLock()
	httpServer, exists := np.servers[serverID]
	np.serverMutex.RUnlock()

	if !exists {
		return &object.String{Value: "服务器不存在"}
	}

	if httpServer.running {
		return &object.String{Value: "服务器已在运行"}
	}

	httpServer.running = true

	// 直接启动服务器（阻塞）
	err := httpServer.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		httpServer.running = false
		return &object.String{Value: fmt.Sprintf("服务器错误: %v", err)}
	}

	httpServer.running = false
	return &object.String{Value: "服务器已停止"}
}

// StopServer 停止HTTP服务器
func (np *NetworkPackage) StopServer(serverID int64) *object.String {
	np.serverMutex.RLock()
	httpServer, exists := np.servers[serverID]
	np.serverMutex.RUnlock()

	if !exists {
		return &object.String{Value: "服务器不存在"}
	}

	if !httpServer.running {
		return &object.String{Value: "服务器未运行"}
	}

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := httpServer.server.Shutdown(ctx)
	if err != nil {
		return &object.String{Value: fmt.Sprintf("停止服务器错误: %v", err)}
	}

	httpServer.running = false
	return &object.String{Value: "服务器已停止"}
}

// IsServerRunning 检查服务器是否运行
func (np *NetworkPackage) IsServerRunning(serverID int64) *object.Boolean {
	np.serverMutex.RLock()
	defer np.serverMutex.RUnlock()

	httpServer, exists := np.servers[serverID]
	if !exists {
		return &object.Boolean{Value: false}
	}

	return &object.Boolean{Value: httpServer.running}
}

// GetServerPort 获取服务器端口
func (np *NetworkPackage) GetServerPort(serverID int64) *object.Integer {
	np.serverMutex.RLock()
	defer np.serverMutex.RUnlock()

	httpServer, exists := np.servers[serverID]
	if !exists {
		return &object.Integer{Value: -1}
	}

	return &object.Integer{Value: int64(httpServer.port)}
}

// HTTPGet 执行HTTP GET请求
func (np *NetworkPackage) HTTPGet(url string) *object.String {
	resp, err := http.Get(url)
	if err != nil {
		return &object.String{Value: fmt.Sprintf("HTTP GET错误: %v", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &object.String{Value: fmt.Sprintf("读取响应错误: %v", err)}
	}

	return &object.String{Value: string(body)}
}

// HTTPPost 执行HTTP POST请求
func (np *NetworkPackage) HTTPPost(url, data string) *object.String {
	resp, err := http.Post(url, "application/json", strings.NewReader(data))
	if err != nil {
		return &object.String{Value: fmt.Sprintf("HTTP POST错误: %v", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &object.String{Value: fmt.Sprintf("读取响应错误: %v", err)}
	}

	return &object.String{Value: string(body)}
}

// TCPConnect 建立TCP连接
func (np *NetworkPackage) TCPConnect(host string, port int64) *object.Integer {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return &object.Integer{Value: -1}
	}

	connID := np.nextID
	np.nextID++
	np.connections[connID] = &conn

	return &object.Integer{Value: connID}
}

// TCPSend 通过TCP发送数据
func (np *NetworkPackage) TCPSend(connID int64, data string) *object.Integer {
	connPtr, exists := np.connections[connID]
	if !exists {
		return &object.Integer{Value: -1}
	}

	conn := *connPtr
	_, err := conn.Write([]byte(data))
	if err != nil {
		return &object.Integer{Value: -1}
	}

	return &object.Integer{Value: int64(len(data))}
}

// TCPReceive 从TCP连接接收数据
func (np *NetworkPackage) TCPReceive(connID int64) *object.String {
	connPtr, exists := np.connections[connID]
	if !exists {
		return &object.String{Value: ""}
	}

	conn := *connPtr
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return &object.String{Value: ""}
	}

	return &object.String{Value: string(buffer[:n])}
}

// TCPClose 关闭TCP连接
func (np *NetworkPackage) TCPClose(connID int64) *object.Integer {
	connPtr, exists := np.connections[connID]
	if !exists {
		return &object.Integer{Value: -1}
	}

	conn := *connPtr
	err := conn.Close()
	if err != nil {
		return &object.Integer{Value: -1}
	}

	delete(np.connections, connID)
	return &object.Integer{Value: 0}
}

// UDPSend 发送UDP数据包
func (np *NetworkPackage) UDPSend(host string, port int64, data string) *object.Integer {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("udp", address)
	if err != nil {
		return &object.Integer{Value: -1}
	}
	defer conn.Close()

	_, err = conn.Write([]byte(data))
	if err != nil {
		return &object.Integer{Value: -1}
	}

	return &object.Integer{Value: int64(len(data))}
}

// ResolveDNS DNS解析
func (np *NetworkPackage) ResolveDNS(hostname string) *object.String {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return &object.String{Value: ""}
	}

	if len(ips) > 0 {
		return &object.String{Value: ips[0].String()}
	}

	return &object.String{Value: ""}
}

// Ping 简单的ping测试
func (np *NetworkPackage) Ping(host string) *object.Integer {
	conn, err := net.DialTimeout("tcp", host+":80", 5*time.Second)
	if err != nil {
		return &object.Integer{Value: -1}
	}
	defer conn.Close()

	return &object.Integer{Value: 0}
}

// NetworkObject 网络对象（用于在CVM环境中表示网络包）
type NetworkObject struct {
	Package *NetworkPackage
}

func (no *NetworkObject) Type() object.ObjectType { return "NETWORK" }
func (no *NetworkObject) Inspect() string         { return "network_package" }
