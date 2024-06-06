package client

import (
    "encoding/json"
    "fmt"
    "github.com/google/uuid"
    "log/slog"
)

type Client struct {
    serverAddress string
    serverPort    int
    clientId      string
    queuedItems   map[string]*QueuePromptResponse
    queuedCount   int
    callbacks     *Callbacks
    wsc           *WebSocketClient
}

type Callbacks struct {
    OnStatus         func(*Client, int)
    OnExecutionStart func(*Client, *QueuePromptResponse)
    OnExecuted       func(*Client, *QueuePromptResponse)
    OnExecuting      func(*Client, *QueuePromptResponse, string)
    OnProgress       func(*Client, *WSStatusMessageDataProgress)
}

func NewClient(serverAddress string, serverPort int, callbacks *Callbacks) (*Client, error) {
    res := &Client{
        serverAddress: serverAddress,
        serverPort:    serverPort,
        clientId:      uuid.New().String(),
        queuedItems:   make(map[string]*QueuePromptResponse),
        callbacks:     callbacks,
    }
    res.wsc = NewWebSocketClient(res)
    err := res.wsc.Connect(serverAddress, serverPort, res.clientId)
    if err != nil {
        return nil, err
    }
    go func() {
        res.wsc.HandleMessages()
    }()
    return res, nil
}

func (c *Client) buildUrl(path string) string {
    return fmt.Sprintf("http://%s:%d%s", c.serverAddress, c.serverPort, path)
}

func (c *Client) ClientId() string {
    return c.clientId
}

func (c *Client) OnMessage(message string) {
    c.OnWebSocketMessage(message)
}

func (c *Client) OnWebSocketMessage(msg string) {
    slog.Info("msg:", "msg", msg)
    message := &WSStatusMessage{}
    err := json.Unmarshal([]byte(msg), message)
    if err != nil {
        slog.Error("Cannot deserialize status message: ", err)
    }

    slog.Info("comfy ws:", "message_type", message.MessageType, "data", message.Data)

    switch message.MessageType {
    case "status":
        s := message.Data.(*WSStatusMessageDataStatus)
        c.queuedCount = s.Status.ExecInfo.QueueRemaining
        if c.callbacks != nil && c.callbacks.OnStatus != nil {
            c.callbacks.OnStatus(c, c.queuedCount)
        }
    case "execution_start":
        s := message.Data.(*WSStatusMessageDataExecutionStart)
        qi := c.GetQueuedItem(s.PromptID)
        if qi != nil {
            if c.callbacks != nil && c.callbacks.OnExecutionStart != nil {
                c.callbacks.OnExecutionStart(c, qi)
            }
            qi.Messages <- "exec start"
        }
    case "execution_cached":
    case "executing":
        s := message.Data.(*WSStatusMessageDataExecuting)
        qi := c.GetQueuedItem(s.PromptID)
        if qi != nil {
            if c.callbacks != nil && c.callbacks.OnExecuting != nil {
                c.callbacks.OnExecuting(c, qi, s.Node)
            }
            qi.Messages <- "executing " + s.Node
        }
    case "progress":
        s := message.Data.(*WSStatusMessageDataProgress)
        if c.callbacks != nil && c.callbacks.OnProgress != nil {
            c.callbacks.OnProgress(c, s)
        }
    case "executed":
        s := message.Data.(*WSStatusMessageDataExecuted)
        qi := c.GetQueuedItem(s.PromptID)
        if qi != nil {
            // mdata := &PromptMessageData{
            // 	NodeID: s.Node,
            // 	Images: *s.Output["images"],
            // }

            // collect the data from the output
            //mdata := &PromptMessageData{
            //    NodeID: s.Node,
            //    Data:   make(map[string][]DataOutput),
            //}
            //
            //for k, v := range s.Output {
            //    mdata.Data[k] = *v
            //}
            //
            //m := PromptMessage{
            //    Type:    "data",
            //    Message: mdata,
            //}
            if c.callbacks != nil && c.callbacks.OnExecuted != nil {
                c.callbacks.OnExecuted(c, qi)
            }
            qi.Messages <- "executed"
        }
    case "execution_interrupted":
    case "execution_error":
    default:
        slog.Warn("Unhandled message type:", "type", message.MessageType)
    }
}
