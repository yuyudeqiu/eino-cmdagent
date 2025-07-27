package agent

import (
	"context"
	"errors"
	"fmt"
	"io"

	chatModel "cmd-agent/internal/model"
	"cmd-agent/internal/template"
	agentTool "cmd-agent/internal/tool"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	callbackHelpers "github.com/cloudwego/eino/utils/callbacks"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type CmdChain struct {
}

func NewCmdChain(ctx context.Context) (*CmdChain, error) {
	return &CmdChain{}, nil
}

type CmdParam struct {
	OS     string // 用户使用的操作系统
	Prompt string // 用户的输入提示词
}

func (cmd *CmdChain) Process(ctx context.Context, param CmdParam) (string, error) {
	// 大模型回调函数
	modelHandler := &callbackHelpers.ModelCallbackHandler{
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *model.CallbackOutput) context.Context {
			// 1. output.Result 类型是 string
			fmt.Println("模型思考过程为：")
			fmt.Println(output.Message.Content)
			return ctx
		},
	}
	// 工具回调函数
	toolHandler := &callbackHelpers.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			fmt.Printf("开始执行工具，参数: %s\n", input.ArgumentsInJSON)
			return ctx
		},
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *tool.CallbackOutput) context.Context {
			fmt.Printf("工具执行完成，结果: %s\n", output.Response)
			return ctx
		},
	}
	// 构建实际回调函数Handler
	handler := callbackHelpers.NewHandlerHelper().
		ChatModel(modelHandler).
		Tool(toolHandler).
		Handler()

	cm, err := chatModel.NewChatModel(ctx)
	if err != nil {
		return "", err
	}

	// 模型绑定工具
	terminalTool := agentTool.TerminalTool{}
	info, _ := terminalTool.Info(ctx)
	if err := cm.BindTools([]*schema.ToolInfo{
		info,
	}); err != nil {
		return "", err
	}

	chain := compose.NewChain[map[string]any, []*schema.Message]()
	// 创建模板，加入chain
	ctp, err := template.NewGetCommandTemplate(ctx)
	if err != nil {
		return "", err
	}
	executeCommandToolNode, err := agentTool.NewExecuteCommandToolNode(ctx)
	if err != nil {
		return "", err
	}
	chain.
		AppendChatTemplate(ctp).
		AppendChatModel(cm, compose.WithNodeName("chat_model")).
		AppendToolsNode(executeCommandToolNode, compose.WithNodeName("execute_command"))
	compile, err := chain.Compile(ctx)
	if err != nil {
		return "", err
	}

	stream, err := compile.Stream(ctx, map[string]any{
		"os":             param.OS,
		"command_prompt": param.Prompt,
	}, compose.WithCallbacks(handler))
	if err != nil {
		return "", err
	}

	defer stream.Close()
	for {
		recv, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		fmt.Print(recv)
	}

	return "", nil
}
