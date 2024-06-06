package objectinfo

type ObjectInfo map[string]InfoContent

type InfoContent struct {
    Input        Input       `json:"input"`
    Output       interface{} `json:"output"`
    OutputIsList []bool      `json:"output_is_list"`
    OutputName   interface{} `json:"output_name"`
    Name         string      `json:"name"`
    DisplayName  string      `json:"display_name"`
    Description  string      `json:"description"`
    Category     string      `json:"category"`
    OutputNode   bool        `json:"output_node"`
}

type Input struct {
    Required map[string][]interface{} `json:"required"`
}

func (o *ObjectInfo) GetByType(typeName string) *InfoContent {
    if val, ok := (*o)[typeName]; ok {
        return &val
    }
    return nil
}

func (o *ObjectInfo) GetListOfCheckpoints() ([]string, error) {
    res := make([]string, 0)
    ckpt := o.GetByType("CheckpointLoader")
    if ckpt == nil {
        return res, nil
    }
    ckptName := ckpt.Input.Required["ckpt_name"]
    if ckptName == nil {
        return res, nil
    }
    list := ckptName[0]
    for _, l := range list.([]interface{}) {
        res = append(res, l.(string))
    }
    return res, nil
}

func (o *ObjectInfo) GetListOfImages() ([]string, error) {
    res := make([]string, 0)
    ckpt := o.GetByType("LoadImage")
    if ckpt == nil {
        return res, nil
    }
    ckptName := ckpt.Input.Required["image"]
    if ckptName == nil {
        return res, nil
    }
    list := ckptName[0]
    for _, l := range list.([]interface{}) {
        res = append(res, l.(string))
    }
    return res, nil
}
