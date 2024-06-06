package client

import (
    "errors"
    "fmt"
    "github.com/gorilla/websocket"
    "log"
    "log/slog"
    "strconv"
)

type WebSocketCallback interface {
    OnMessage(message string)
}
type WebSocketClient struct {
    connectedChannel chan bool
    connection       *websocket.Conn
    callback         WebSocketCallback
}

func NewWebSocketClient(callback WebSocketCallback) *WebSocketClient {
    connectedChannel := make(chan bool, 1)
    return &WebSocketClient{
        connectedChannel: connectedChannel,
        callback:         callback,
    }
}

func (ws *WebSocketClient) Connect(serverAddress string, serverPort int, clientId string) error {
    conn, _, err := websocket.DefaultDialer.Dial("ws://"+serverAddress+":"+strconv.Itoa(serverPort)+"/ws?clientId="+clientId, nil)
    if err != nil {
        return err
    }
    conn.SetReadLimit(2000000)
    conn.SetCloseHandler(func(code int, text string) error {
        log.Println("Close handler: " + text)
        return errors.New("close error")
    })
    ws.connection = conn
    return nil
}

func (ws *WebSocketClient) Close() error {
    if ws.connection != nil {
        err := ws.connection.Close()
        if err != nil {
            return err
        }
    }
    return nil
}

func (ws *WebSocketClient) Ping() error {
    return ws.connection.WriteMessage(websocket.PingMessage, []byte("ping"))
}

func (ws *WebSocketClient) HandleMessages() {
    for {
        _, message, err := ws.connection.ReadMessage()
        if err != nil {
            slog.Warn(fmt.Sprintf("Read error: %v", err))
            break
        }
        if ws.callback != nil {
            ws.callback.OnMessage(string(message))
        }
    }
    slog.Info("Finished HandleMessages")
}
