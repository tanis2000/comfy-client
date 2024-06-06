package graph

import (
    "encoding/json"
    "strconv"
)

type Graph struct {
    Nodes      []*Node       `json:"nodes"`
    Links      []*Link       `json:"links"`
    Groups     []*Group      `json:"groups"`
    LastNodeID int           `json:"last_node_id"`
    LastLinkID int           `json:"last_link_id"`
    Version    float32       `json:"version"`
    NodesByID  map[int]*Node `json:"-"`
    LinksByID  map[int]*Link `json:"-"`
}

type QueuePromptNode struct {
    Inputs    map[string]interface{} `json:"inputs"`
    ClassType string                 `json:"class_type"`
}

func NewGraph(jsonData string) (*Graph, error) {
    graph := &Graph{
        NodesByID: make(map[int]*Node),
        LinksByID: make(map[int]*Link),
    }
    err := json.Unmarshal([]byte(jsonData), &graph)
    if err != nil {
        return nil, err
    }
    for _, node := range graph.Nodes {
        node.Graph = graph
        graph.NodesByID[node.ID] = node
    }

    for _, link := range graph.Links {
        graph.LinksByID[link.ID] = link
    }

    return graph, nil
}

func (g *Graph) GetNodeByID(id int) *Node {
    return g.NodesByID[id]
}

func (g *Graph) GetLinkByID(id int) *Link {
    return g.LinksByID[id]
}

func (g *Graph) GraphToPromptNodes() (map[int]QueuePromptNode, error) {
    nodes := make(map[int]QueuePromptNode)
    for _, node := range g.Nodes {
        if node.IsVirtual() {
            // Don't serialize frontend only nodes but let them make changes
            continue
        }

        pn := QueuePromptNode{
            ClassType: node.Type,
            Inputs:    make(map[string]interface{}),
        }

        for k, prop := range *node.Properties {
            pn.Inputs[k] = prop
        }

        for i, slot := range node.Inputs {
            parent := node.GetNodeForInput(i)
            if parent != nil {
                link := g.GetLinkByID(slot.Link)
                for parent != nil && parent.IsVirtual() {
                    link = parent.GetInputLink(link.OriginSlot)
                    if link != nil {
                        parent = parent.GetNodeForInput(link.OriginSlot)
                    } else {
                        break
                    }
                }

                if link != nil {
                    linfo := make([]interface{}, 2)
                    linfo[0] = strconv.Itoa(link.OriginID)
                    linfo[1] = link.OriginSlot
                    pn.Inputs[node.Inputs[i].Name] = linfo
                }
            }
        }
        nodes[node.ID] = pn
    }
    return nodes, nil
}
