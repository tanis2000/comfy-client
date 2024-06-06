package workflow

type WorkflowInput map[string]interface{}

func (input *WorkflowInput) Set(key string, value interface{}) {
    (*input)[key] = value
}
