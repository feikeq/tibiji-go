package controllers

import (
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
)

type WebsocketController struct {
	Conn *neffos.Server
}

func NewWebsocketController() *WebsocketController {
	// 当应用程序只接受和发送原始的websocket原生消息时，neffos的几乎所有特性都被禁用，因为没有自定义消息可以传递。
	// 什么情况下只允许原生消息呢？ 当注册的命名空间只有一个，并且是空的， 并且只包含一个已注册的事件，即OnNativeMessage。
	// 当使用Events{...}代替Namespaces{ "namespaceName": Events{...}}时， 那么命名空间就是空的""。
	ws := websocket.New(websocket.DefaultGobwasUpgrader, websocket.Events{
		//用websocket.DefaultGobwasUpgrader来启用CORS。
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			println("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

			pong := "server say " + string(msg.Body)

			mg := websocket.Message{
				Body:     []byte(pong),
				IsNative: true,
			}

			nsConn.Conn.Write(mg)

			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},
	})

	ws.OnConnect = func(c *websocket.Conn) error {
		println("[%s] Connected to server!", c.ID())
		return nil
	}

	ws.OnDisconnect = func(c *websocket.Conn) {
		println("[%s] Disconnected from server", c.ID())
	}
	ws.OnUpgradeError = func(err error) {
		println("Upgrade Error: %v", err)
	}

	return &WebsocketController{
		Conn: ws,
	}
}

/*
func (c *WebsocketController) Get(ctx iris.Context) {
	websocket.Serve(ctx, c.Conn)
}

func (c *WebsocketController) OnPing(msg websocket.Message) {
	c.Conn.To(websocket.Broadcast).Emit("pong", msg)
}

func (c *WebsocketController) OnNamespaceDisconnect(msg websocket.Message) error {
	fmt.Printf("Namespace disconnected: %s\n", msg.Body)
	return nil
}

func (c *WebsocketController) OnNamespaceConnected(msg websocket.Message) error {
	fmt.Printf("Namespace connected: %s\n", msg.Body)
	return nil
}

func (c *WebsocketController) OnChat(msg websocket.Message) error {
	fmt.Printf("Chat message: %s\n", msg.Body)
	return nil
}

*/
