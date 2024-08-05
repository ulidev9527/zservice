package main

import (
	"net/http"

	"github.com/ulidev9527/zservice/zservice"
	"github.com/ulidev9527/zservice/zserviceex/nbioservice"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

func init() {
	zservice.Init("nbioservice.test", "1.0.0")
}
func main() {
	n := 0
	nbioS := nbioservice.NewNbioService_WebSocket(nbioservice.NbioServiceOption_WebScoket{
		ListenPort: "8801",
		OnRequest: func(w http.ResponseWriter, r *http.Request) bool {
			n++
			allow := n%2 == 0

			if !allow {
				w.Write([]byte("close this"))
			}

			return allow
		},
		OnOpen: func(c *websocket.Conn) {
			zservice.LogInfo("OnOpen:", c.RemoteAddr())
		},
		OnMessage: func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
			zservice.LogInfo("OnMessage", messageType, data)
			c.WriteMessage(messageType, data)
		},
		OnClose: func(c *websocket.Conn, err error) {
			zservice.LogInfo("OnClose", c.RemoteAddr())
		},
	})

	zservice.AddDependService(nbioS.ZService)

	zservice.Start().WaitStop()
}
