package client

import (
    "bytes"
    "encoding/json"
    "errors"
    "github.com/tanis2000/comfy-client/objectinfo"
    "github.com/tanis2000/comfy-client/workflow"
    "io"
    "net/http"
    "net/url"
    "strconv"
    "strings"
)

func (c *Client) GetExtensions() ([]string, error) {
    res, err := http.Get(c.buildUrl("/extensions"))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetEmbeddings() ([]string, error) {
    res, err := http.Get(c.buildUrl("/embeddings"))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetObjectInfo() (*objectinfo.ObjectInfo, error) {
    res, err := http.Get(c.buildUrl("/object_info"))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := &objectinfo.ObjectInfo{}
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) QueuePrompt(number int, workflow *workflow.Workflow) (*QueuePromptResponse, error) {
    req := &QueuePromptRequest{
        ClientId: c.clientId,
        Prompt:   &workflow.Map,
        //ExtraData: QueuePromptExtraData{ExtraPngInfo: QueuePromptWorkflow{Workflow: &workflow.Map}}
    }
    if number == -1 {
        req.Front = true
    } else if number != 0 {
        req.Number = number
    }

    reqBody, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    res, err := http.Post(c.buildUrl("/prompt"), "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return nil, err
    }

    if res.StatusCode != 200 {
        body, _ := io.ReadAll(res.Body)
        return nil, errors.New(string(body))
    }

    body, _ := io.ReadAll(res.Body)
    queueItem := &QueuePromptResponse{}
    err = json.Unmarshal(body, &queueItem)
    if err != nil {
        return nil, err
    }
    queueItem.Messages = make(chan string)
    c.queuedItems[queueItem.PromptID] = queueItem

    return queueItem, nil
}

func (c *Client) PollPrompt() ([]string, error) {
    res, err := http.Get(c.buildUrl("/prompt"))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetQueue() ([]string, error) {
    res, err := http.Get(c.buildUrl("/queue"))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetHistory(maxItems int) ([]string, error) {
    if maxItems == 0 {
        maxItems = 200
    }
    res, err := http.Get(c.buildUrl("/history?max_items=" + strconv.Itoa(maxItems)))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetHistoryByPromptID(promptID string, maxItems int) (*HistoryPrompt, error) {
    if maxItems == 0 {
        maxItems = 200
    }
    res, err := http.Get(c.buildUrl("/history/" + promptID + "?max_items=" + strconv.Itoa(maxItems)))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := &HistoryPrompt{}
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetSystemStats() (*SystemStatsResponse, error) {
    res, err := http.Get(c.buildUrl("/system_stats"))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := &SystemStatsResponse{}
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) DeleteQueueItem(id int) ([]string, error) {
    req := &DeleteRequest{
        Delete: []int{id},
    }

    reqBody, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    res, err := http.Post(c.buildUrl("/queue"), "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) ClearQueue() ([]string, error) {
    req := &ClearRequest{
        Clear: true,
    }

    reqBody, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    res, err := http.Post(c.buildUrl("/queue"), "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) InterruptQueue() ([]string, error) {
    res, err := http.Post(c.buildUrl("/queue"), "application/json", nil)
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) DeleteHistoryItem(id int) ([]string, error) {
    req := &DeleteRequest{
        Delete: []int{id},
    }

    reqBody, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    res, err := http.Post(c.buildUrl("/history"), "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) ClearHistory() ([]string, error) {
    req := &ClearRequest{
        Clear: true,
    }

    reqBody, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    res, err := http.Post(c.buildUrl("/history"), "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) InterruptHistory() ([]string, error) {
    res, err := http.Post(c.buildUrl("/history"), "application/json", nil)
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetUserConfig() ([]string, error) {
    res, err := http.Get(c.buildUrl("/users"))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) CreateUser(username string) ([]string, error) {
    res, err := http.Post(c.buildUrl("/users"), "application/json", strings.NewReader(username))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetSettings() ([]string, error) {
    res, err := http.Get(c.buildUrl("/settings"))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetSetting(id string) ([]string, error) {
    res, err := http.Get(c.buildUrl("/settings/" + url.QueryEscape(id)))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) StoreSettings(settings map[string]interface{}) ([]string, error) {
    reqBody, err := json.Marshal(settings)
    if err != nil {
        return nil, err
    }

    res, err := http.Post(c.buildUrl("/settings"), "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) StoreSetting(id string, value interface{}) ([]string, error) {
    reqBody, err := json.Marshal(value)
    if err != nil {
        return nil, err
    }

    res, err := http.Post(c.buildUrl("/settings/"+url.QueryEscape(id)), "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetUserData(file string, data interface{}) ([]string, error) {
    reqBody, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    res, err := http.Post(c.buildUrl("/userdata/"+url.QueryEscape(file)), "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)
    jsonBody := make([]string, 0)
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return nil, err
    }

    return jsonBody, nil
}

func (c *Client) GetView(filename string, subfolder string, imageType string) ([]byte, error) {
    res, err := http.Get(c.buildUrl("/view?filename=" + url.QueryEscape(filename) + "&subfolder=" + url.QueryEscape(subfolder) + "&type=" + url.QueryEscape(imageType)))
    if err != nil {
        return nil, err
    }

    body, _ := io.ReadAll(res.Body)

    return body, nil
}
