# ComfyUI Go client library

comfy-client is a REST/WS client library for [ComfyUI](https://github.com/comfyanonymous/ComfyUI) written in Go.

The aim is to provide an easy way to interface with the REST and WebSocket based API offered by ComfyUI.

This is still in early development. Contributions are welcome.

# Basic usage

```go
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
    c, err := client.NewClient("localhost", 8188, callbacks)
    if err != nil {
        panic(err)
    }

    println("Loading JSON workflow API")
    f, err := os.Open("examples/txt2img/txt2img_api.json")
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
    println("Building workflow from workflow api")
    w, err := workflow.NewWorkflow(string(content))
    if err != nil {
        panic(err)
    }

    println("Enqueueing a prompt")
    res, err := c.QueuePrompt(-1, w)
    if err != nil {
        panic(err)
    }
    println(res)

    println("Starting the message loop")
    for continueLoop := true; continueLoop; {
        msg := <-res.Messages
        println(msg)
    }

}
```

Examples of how to use this library can be found in [examples](examples)

