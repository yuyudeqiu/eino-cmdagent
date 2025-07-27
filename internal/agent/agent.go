package agent

import (
	"context"
	"errors"
	"fmt"
	"io"

	"cmd-agent/internal/model"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	mytool "cmd-agent/internal/tool"
)

func NewCmdAgent() *CmdAgent {
	return &CmdAgent{}
}

type CmdAgent struct{}

func (ca *CmdAgent) Process(ctx context.Context, text string) error {
	cm, err := model.NewArkModel(ctx)
	if err != nil {
		return err
	}
	terminalTool := mytool.NewTerminalTool()
	info, _ := terminalTool.Info(ctx)
	if err = cm.BindTools([]*schema.ToolInfo{info}); err != nil {
		return err
	}

	ragent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: cm,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{terminalTool},
		},
		MaxStep: 10,
	})
	if err != nil {
		return err
	}
	stream, err := ragent.Stream(ctx, []*schema.Message{
		schema.SystemMessage(`
你是一个终端命令小助手，请根据用户的需求，选择合适的命令并且调用终端工具帮助用户执行命令，
你不需要询问用户是否需要完成确定什么，直接开始就好了。另外，如果需求比较复杂，需要执行多条命令，你可以分成多个步骤一次一次的调用工具来完成`),
		{
			Role:    schema.User,
			Content: text,
		},
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Print(msg.Content)
	}
	return nil
}
