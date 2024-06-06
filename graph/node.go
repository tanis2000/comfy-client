package graph

type Node struct {
    ID         int                     `json:"id"`
    Type       string                  `json:"type"`
    Position   interface{}             `json:"pos"`
    Size       Size                    `json:"size"`
    Flags      *interface{}            `json:"flags"`
    Order      int                     `json:"order"`
    Mode       int                     `json:"mode"`
    Title      string                  `json:"title,omitempty"`
    Properties *map[string]interface{} `json:"properties"`
    // widgets_values can be an array of values, or a map of values
    // maps of values can represent cascading style properties in which the setting
    // of one property makes certain other properties available
    WidgetValues interface{} `json:"widgets_values,omitempty"`
    Color        string      `json:"color"`
    BGColor      string      `json:"bgcolor"`
    Inputs       []Slot      `json:"inputs,omitempty"`
    Outputs      []Slot      `json:"outputs,omitempty"`
    Graph        *Graph      `json:"-"`
}

func (n *Node) GetNodeForInput(slotIndex int) *Node {
    if slotIndex >= len(n.Inputs) {
        return nil
    }

    slot := n.Inputs[slotIndex]
    l := n.Graph.GetLinkByID(slot.Link)
    if l == nil {
        return nil
    }
    return n.Graph.GetNodeByID(l.OriginID)
}

func (n *Node) IsVirtual() bool {
    // current nodes that are 'virtual':
    switch n.Type {
    case "PrimitiveNode":
        return true
    case "Reroute":
        return true
    case "Note":
        return true
    }
    return false
}

func (n *Node) GetInputLink(slotIndex int) *Link {
    count := len(n.Inputs)
    if count == 0 || slotIndex >= count {
        return nil
    }

    slot := n.Inputs[slotIndex]
    return n.Graph.GetLinkByID(slot.Link)
}
