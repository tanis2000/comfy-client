package client

import (
    "encoding/json"
    "github.com/tanis2000/comfy-client/workflow"
)

type QueuePromptWorkflow struct {
    Workflow *workflow.WorkflowMap `json:"workflow"`
}
type QueuePromptExtraData struct {
    ExtraPngInfo QueuePromptWorkflow `json:"extra_pnginfo"`
}

type QueuePromptRequest struct {
    ClientId  string                `json:"client_id"`
    Prompt    *workflow.WorkflowMap `json:"prompt"`
    ExtraData QueuePromptExtraData  `json:"extra_data"`
    Front     bool                  `json:"front,omitempty"`
    Number    int                   `json:"number,omitempty"`
}

type QueuePromptResponse struct {
    PromptID   string                 `json:"prompt_id"`
    Number     int                    `json:"number"`
    NodeErrors map[string]interface{} `json:"node_errors"`
    Messages   chan string            `json:"-"`
}

type DeleteRequest struct {
    Delete []int `json:"delete"`
}

type ClearRequest struct {
    Clear bool `json:"clear"`
}

type SystemStatsResponse struct {
    System  System `json:"system"`
    Devices []GPU  `json:"devices"`
}

type System struct {
    OS             string `json:"os"`
    PythonVersion  string `json:"python_version"`
    EmbeddedPython bool   `json:"embedded_python"`
}

type GPU struct {
    Name           string `json:"name"`
    Type           string `json:"type"`
    Index          int    `json:"index"`
    VRAMTotal      int64  `json:"vram_total"`
    VRAMFree       int64  `json:"vram_free"`
    TorchVRAMTotal int64  `json:"torch_vram_total"`
    TorchVRAMFree  int64  `json:"torch_vram_free"`
}

type DataOutput struct {
    Filename  string `json:"filename"`
    Subfolder string `json:"subfolder"`
    Type      string `json:"type"`
    Text      string `json:"-"` // for "text" type data output
}

type HistoryPrompt map[string]HistoryContent

type HistoryContent struct {
    Prompt  []interface{}            `json:"prompt"`
    Outputs map[string]HistoryOutput `json:"outputs"`
}

type HistoryOutput struct {
    Images []HistoryImage `json:"images"`
}

type HistoryImage struct {
    Filename  string `json:"filename"`
    Subfolder string `json:"subfolder"`
    Type      string `json:"type"`
}

func (ho *HistoryOutput) GetImagesByType(imageType string) []HistoryImage {
    res := make([]HistoryImage, 0)
    for _, img := range ho.Images {
        if img.Type == imageType {
            res = append(res, img)
        }
    }
    return res
}

func (ho *HistoryContent) GetImagesByType(imageType string) []HistoryImage {
    res := make([]HistoryImage, 0)
    for _, output := range ho.Outputs {
        res = append(res, output.GetImagesByType(imageType)...)
    }
    return res
}

type WSStatusMessage struct {
    MessageType string      `json:"type"`
    Data        interface{} `json:"data"`
}

type WSStatusMessageDataStatus struct {
    Status struct {
        ExecInfo struct {
            QueueRemaining int `json:"queue_remaining"`
        } `json:"exec_info"`
    } `json:"status"`
    SID string `json:"sid"`
}

type WSStatusMessageDataExecutionStart struct {
    PromptID string `json:"prompt_id"`
}

type WSStatusMessageDataExecutionCached struct {
}

type WSStatusMessageDataExecuting struct {
    Node     string `json:"node"`
    PromptID string `json:"prompt_id"`
}

type WSStatusMessageDataProgress struct {
    Value int `json:"value"`
    Max   int `json:"max"`
}

type WSStatusMessageDataExecuted struct {
    Node     int                      `json:"node"`
    Output   map[string]*[]DataOutput `json:"output"`
    PromptID string                   `json:"prompt_id"`
}

type WSStatusMessageDataExecutionInterrupted struct {
}

type WSStatusMessageDataExecutionError struct {
}

func (sm *WSStatusMessage) UnmarshalJSON(b []byte) error {
    var temp struct {
        MessageType string          `json:"type"`
        Data        json.RawMessage `json:"data"`
    }
    if err := json.Unmarshal(b, &temp); err != nil {
        return err
    }
    sm.MessageType = temp.MessageType
    switch sm.MessageType {
    case "status":
        sm.Data = &WSStatusMessageDataStatus{}
    case "execution_start":
        sm.Data = &WSStatusMessageDataExecutionStart{}
    case "execution_cached":
        sm.Data = &WSStatusMessageDataExecutionCached{}
    case "executing":
        sm.Data = &WSStatusMessageDataExecuting{}
    case "progress":
        sm.Data = &WSStatusMessageDataProgress{}
    case "executed":
        sm.Data = &WSStatusMessageDataExecuted{}
    case "execution_interrupted":
        sm.Data = &WSStatusMessageDataExecutionInterrupted{}
    case "execution_error":
        sm.Data = &WSStatusMessageDataExecutionError{}
    default:
        sm.Data = nil
    }
    if sm.Data != nil {
        // Unmarshal the data into the selected type
        if err := json.Unmarshal(temp.Data, sm.Data); err != nil {
            return err
        }
    }
    return nil
}
