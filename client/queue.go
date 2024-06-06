package client

func (c *Client) GetQueuedItem(ID string) *QueuePromptResponse {
    if val, ok := c.queuedItems[ID]; ok {
        return val
    }
    return nil
}
