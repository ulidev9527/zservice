package nbioservice

import (
	"fmt"
	"net/http"
	"zservice/zservice"

	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

type NbioService_WebSocket struct {
	*zservice.ZService
	upgrader *websocket.Upgrader
}

type NbioServiceOption_WebScoket struct {
	ListenPort string                                                                  // 监听的端口
	MaxLoad    uint                                                                    // 最大连接数 默认:10K
	OnRequest  func(w http.ResponseWriter, r *http.Request) bool                       // 响应请求, 在 Upgrader 之前调用, 返回 false 表示拒绝进行 Upgrader 操作
	OnOpen     func(c *websocket.Conn)                                                 // 响应链接
	OnMessage  func(c *websocket.Conn, messageType websocket.MessageType, data []byte) // 响应消息
	OnClose    func(c *websocket.Conn, err error)                                      // 响应关闭
}

// websocket 连接
func NewNbioService_WebSocket(opt NbioServiceOption_WebScoket) *NbioService_WebSocket {
	name := fmt.Sprint("NbioService_WebSocket-", opt.ListenPort)
	ns := &NbioService_WebSocket{}

	ns.ZService = zservice.NewService(name, func(z *zservice.ZService) {

		// upgrader 创建
		ns.upgrader = websocket.NewUpgrader()
		ns.upgrader.BlockingModAsyncWrite = true
		ns.upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		if opt.OnOpen != nil {
			ns.upgrader.OnOpen(opt.OnOpen)
		}
		if opt.OnMessage != nil {
			ns.upgrader.OnMessage(opt.OnMessage)
		}
		if opt.OnClose != nil {
			ns.upgrader.OnClose(opt.OnClose)
		}

		// listener 创建
		mux := &http.ServeMux{}
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
				w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
				w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if opt.OnRequest != nil && !opt.OnRequest(w, r) {
				return
			}

			conn, err := ns.upgrader.Upgrade(w, r, nil)
			if err != nil {
				z.LogPanic(err)
			}
			z.LogInfo("Upgraded:", conn.RemoteAddr().String())
		})

		engine := nbhttp.NewEngine(nbhttp.Config{
			Network:                 "tcp",
			Addrs:                   []string{fmt.Sprint("0.0.0.0:", opt.ListenPort)},
			MaxLoad:                 zservice.MaxInt(int(opt.MaxLoad), 10000),
			ReleaseWebsocketPayload: true,
			Handler:                 mux,
		})

		zservice.Go(func() {
			e := engine.Start()
			if e != nil {
				z.LogPanic(e)
			}
		})
		z.LogInfof("nbioService Listen on :%v", opt.ListenPort)
		z.StartDone()
	})
	return ns
}
