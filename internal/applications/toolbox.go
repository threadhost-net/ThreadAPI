package applications

import "github.com/threadhost-net/ThreadAPI/platform"

type ToolBox struct {
	*platform.Stack
}

func NewToolBox(stack *platform.Stack) *ToolBox {
	return &ToolBox{Stack: stack}
}
