package main

import (
    "encoding/json"
    "fmt"
    "github.com/tanis2000/comfy-client/client"
    "github.com/tanis2000/comfy-client/graph"
    "github.com/tanis2000/comfy-client/workflow"
    "io"
    "log"
    "os"
)

func main() {
    callbacks := &client.Callbacks{
        OnStatus: func(c *client.Client, queuedItems int) {
            log.Printf("Queue size: %d", queuedItems)
        },
    }
    c := client.NewClient("localhost", 8188, callbacks)

    println("Getting System Stats")
    stats, err := c.GetSystemStats()
    if err != nil {
        panic(err)
    }
    fmt.Printf("%v\n", stats)

    info, err := c.GetObjectInfo()
    if err != nil {
        panic(err)
    }
    fmt.Printf("%v\n", info)

    println("Loading JSON workflow")
    f, err := os.Open("examples/txt2img/txt2img.json")
    if err != nil {
        panic(err)
    }
    content, err := io.ReadAll(f)
    if err != nil {
        panic(err)
    }
    err = f.Close()
    if err != nil {
        panic(err)
    }

    println("Building graph from workflow")
    g, err := graph.NewGraph(string(content))
    if err != nil {
        panic(err)
    }
    j, err := json.MarshalIndent(g, "", "  ")
    if err != nil {
        panic(err)
    }
    out, err := os.Create("examples/txt2img/txt2img2.json")
    if err != nil {
        panic(err)
    }
    _, err = out.WriteString(string(j))
    if err != nil {
        panic(err)
    }
    err = out.Close()
    if err != nil {
        panic(err)
    }

    nodes, err := g.GraphToPromptNodes()
    if err != nil {
        panic(err)
    }
    j, err = json.MarshalIndent(nodes, "", "  ")
    if err != nil {
        panic(err)
    }
    out, err = os.Create("examples/txt2img/prompt.json")
    if err != nil {
        panic(err)
    }
    _, err = out.WriteString(string(j))
    if err != nil {
        panic(err)
    }
    err = out.Close()
    if err != nil {
        panic(err)
    }

    println("Loading JSON workflow API")
    f, err = os.Open("examples/txt2img/txt2img_api.json")
    if err != nil {
        panic(err)
    }
    content, err = io.ReadAll(f)
    if err != nil {
        panic(err)
    }
    err = f.Close()
    if err != nil {
        panic(err)
    }
    println("Building workflow from workflow api")
    w, err := workflow.NewWorkflow(string(content))
    if err != nil {
        panic(err)
    }
    j, err = json.MarshalIndent(w, "", "  ")
    if err != nil {
        panic(err)
    }
    out, err = os.Create("examples/txt2img/txt2img2_api.json")
    if err != nil {
        panic(err)
    }
    _, err = out.WriteString(string(j))
    if err != nil {
        panic(err)
    }
    err = out.Close()
    if err != nil {
        panic(err)
    }

    println("Enqueueing a prompt")
    res, err := c.QueuePrompt(-1, w)
    if err != nil {
        panic(err)
    }
    println(res)

    println("Starting the websocket")
    wsc := client.NewWebSocketClient(c)
    err = wsc.Connect("localhost", 8188, c.ClientId())
    if err != nil {
        panic(err)
    }
    println("Pinging the websocket")
    err = wsc.Ping()
    if err != nil {
        panic(err)
    }
    go func() {
        println("Handling messages")
        wsc.HandleMessages()
    }()

    println("Starting the loop")
    for continueLoop := true; continueLoop; {
        msg := <-res.Messages
        println(msg)
    }
    err = wsc.Close()
    if err != nil {
        panic(err)
    }
}
