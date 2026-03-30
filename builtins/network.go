package builtins

import (
	"memmole/pkg/ast"
	"memmole/pkg/parser"
	"memmole/network"
	"memmole/pkg/object"
	"memmole/pkg/logger"
	"fmt"
	"net/http"
)

// EvalNetworkCall 处理网络包调用
func EvalNetworkCall(exp *ast.MemberAccessExpression, args []object.Object, env *parser.Environment) object.Object {
	// 获取网络包对象
	networkObj, exists := env.Get("network")
	if !exists {
		return &object.String{Value: "网络包未找到"}
	}

	networkPackage, ok := networkObj.(*network.NetworkObject)
	if !ok {
		return &object.String{Value: "无效的网络包对象"}
	}

	np := networkPackage.Package

	switch exp.Member.Value {
	// 客户端功能
	case "http_get":
		if len(args) != 1 {
			return &object.String{Value: "http_get需要1个参数: URL"}
		}
		url, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "http_get参数必须是字符串"}
		}
		return np.HTTPGet(url.Value)

	case "http_post":
		if len(args) != 2 {
			return &object.String{Value: "http_post需要2个参数: URL, data"}
		}
		url, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "http_post第一个参数必须是字符串"}
		}
		data, ok := args[1].(*object.String)
		if !ok {
			return &object.String{Value: "http_post第二个参数必须是字符串"}
		}
		return np.HTTPPost(url.Value, data.Value)

	case "tcp_connect":
		if len(args) != 2 {
			return &object.String{Value: "tcp_connect需要2个参数: host, port"}
		}
		host, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "tcp_connect第一个参数必须是字符串"}
		}
		port, ok := args[1].(*object.Integer)
		if !ok {
			return &object.String{Value: "tcp_connect第二个参数必须是整数"}
		}
		return np.TCPConnect(host.Value, port.Value)

	case "tcp_send":
		if len(args) != 2 {
			return &object.String{Value: "tcp_send需要2个参数: conn_id, data"}
		}
		connID, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "tcp_send第一个参数必须是整数"}
		}
		data, ok := args[1].(*object.String)
		if !ok {
			return &object.String{Value: "tcp_send第二个参数必须是字符串"}
		}
		return np.TCPSend(connID.Value, data.Value)

	case "tcp_receive":
		if len(args) != 1 {
			return &object.String{Value: "tcp_receive需要1个参数: conn_id"}
		}
		connID, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "tcp_receive参数必须是整数"}
		}
		return np.TCPReceive(connID.Value)

	case "tcp_close":
		if len(args) != 1 {
			return &object.String{Value: "tcp_close需要1个参数: conn_id"}
		}
		connID, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "tcp_close参数必须是整数"}
		}
		return np.TCPClose(connID.Value)

	case "udp_send":
		if len(args) != 3 {
			return &object.String{Value: "udp_send需要3个参数: host, port, data"}
		}
		host, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "udp_send第一个参数必须是字符串"}
		}
		port, ok := args[1].(*object.Integer)
		if !ok {
			return &object.String{Value: "udp_send第二个参数必须是整数"}
		}
		data, ok := args[2].(*object.String)
		if !ok {
			return &object.String{Value: "udp_send第三个参数必须是字符串"}
		}
		return np.UDPSend(host.Value, port.Value, data.Value)

	case "resolve_dns":
		if len(args) != 1 {
			return &object.String{Value: "resolve_dns需要1个参数: hostname"}
		}
		hostname, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "resolve_dns参数必须是字符串"}
		}
		return np.ResolveDNS(hostname.Value)

	case "ping":
		if len(args) != 1 {
			return &object.String{Value: "ping需要1个参数: host"}
		}
		host, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "ping参数必须是字符串"}
		}
		return np.Ping(host.Value)

	// 服务器功能
	case "create_server":
		if len(args) != 1 {
			return &object.String{Value: "create_server需要1个参数: port"}
		}
		port, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "create_server参数必须是整数"}
		}
		result := np.CreateHTTPServer(port.Value)
		
		// 检查是否创建成功
		if result.Value == -1 {
			logger.Error("无法创建服务器：在端口 %d 附近找不到可用端口", port.Value)
			return &object.String{Value: fmt.Sprintf("创建服务器失败：无法找到可用端口")}
		}
		
		// 获取实际使用的端口
		actualPort := np.GetServerPort(result.Value)
		logger.Info("创建HTTP服务器成功，服务器ID: %d，实际端口: %d", result.Value, actualPort.Value)
		return result

	case "add_route":
		if len(args) != 5 {
			return &object.String{Value: "add_route需要5个参数: server_id, method, path, content_type, content"}
		}
		serverID, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "add_route第一个参数必须是整数"}
		}
		method, ok := args[1].(*object.String)
		if !ok {
			return &object.String{Value: "add_route第二个参数必须是字符串"}
		}
		path, ok := args[2].(*object.String)
		if !ok {
			return &object.String{Value: "add_route第三个参数必须是字符串"}
		}
		contentType, ok := args[3].(*object.String)
		if !ok {
			return &object.String{Value: "add_route第四个参数必须是字符串"}
		}
		content, ok := args[4].(*object.String)
		if !ok {
			return &object.String{Value: "add_route第五个参数必须是字符串"}
		}
		
		// 根据内容类型设置相应的Content-Type
		var mimeType string
		switch contentType.Value {
		case "text":
			mimeType = "text/plain; charset=utf-8"
		case "html":
			mimeType = "text/html; charset=utf-8"
		case "json":
			mimeType = "application/json; charset=utf-8"
		case "xml":
			mimeType = "application/xml; charset=utf-8"
		default:
			mimeType = "text/plain; charset=utf-8"
		}
		
		// 创建处理器函数
		handlerFunc := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", mimeType)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(content.Value))
		}
		
		result := np.AddRoute(serverID.Value, method.Value, path.Value, handlerFunc)
		logger.Info("添加路由: %s %s (Content-Type: %s)", method.Value, path.Value, mimeType)
		return result

	case "start_server":
		if len(args) != 1 {
			return &object.String{Value: "start_server需要1个参数: server_id"}
		}
		serverID, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "start_server参数必须是整数"}
		}
		return np.StartServer(serverID.Value)

	case "start_server_and_wait":
		if len(args) != 1 {
			return &object.String{Value: "start_server_and_wait需要1个参数: server_id"}
		}
		serverID, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "start_server_and_wait参数必须是整数"}
		}
		result := np.StartServerAndWait(serverID.Value)
		
		// 添加专门的日志输出
		if result.Value != "服务器已停止" {
			logger.Error("服务器启动失败: %s", result.Value)
		} else {
			logger.Info("服务器已停止")
		}
		
		return result

	case "stop_server":
		if len(args) != 1 {
			return &object.String{Value: "stop_server需要1个参数: server_id"}
		}
		serverID, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "stop_server参数必须是整数"}
		}
		return np.StopServer(serverID.Value)

	case "is_server_running":
		if len(args) != 1 {
			return &object.String{Value: "is_server_running需要1个参数: server_id"}
		}
		serverID, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "is_server_running参数必须是整数"}
		}
		return np.IsServerRunning(serverID.Value)

	case "get_server_port":
		if len(args) != 1 {
			return &object.String{Value: "get_server_port需要1个参数: server_id"}
		}
		serverID, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "get_server_port参数必须是整数"}
		}
		return np.GetServerPort(serverID.Value)

	default:
		return &object.String{Value: fmt.Sprintf("未知的网络方法: %s", exp.Member.Value)}
	}
}
