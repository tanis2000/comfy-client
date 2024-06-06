package workflow

import "encoding/json"

type WorkflowMeta struct {
    Title string `json:"title"`
}

type WorkflowNode struct {
    Inputs    WorkflowInput `json:"inputs"`
    ClassType string        `json:"class_type"`
    Meta      WorkflowMeta  `json:"_meta"`
}
type WorkflowMap map[string]WorkflowNode

type Workflow struct {
    Map WorkflowMap
}

func NewWorkflow(jsonData string) (*Workflow, error) {
    workflow := &Workflow{
        Map: make(WorkflowMap),
    }
    err := json.Unmarshal([]byte(jsonData), &workflow.Map)
    if err != nil {
        return nil, err
    }

    return workflow, nil
}

func (workflow *Workflow) NodeByID(id string) *WorkflowNode {
    if val, ok := workflow.Map[id]; ok {
        return &val
    }
    return nil
}
