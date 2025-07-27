package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/cloudwego/eino/compose"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func NewTerminalTool() *TerminalTool {
	return &TerminalTool{}
}

func NewExecuteCommandToolNode(ctx context.Context) (tsn *compose.ToolsNode, err error) {
	config := &compose.ToolsNodeConfig{}
	config.Tools = []tool.BaseTool{&TerminalTool{}}
	tsn, err = compose.NewToolNode(ctx, config)
	if err != nil {
		return nil, err
	}
	return tsn, nil
}

// TerminalTool 用于执行终端命令的工具
type TerminalTool struct {
	// 可以添加超时设置等配置
	Timeout time.Duration
}

// 命令参数结构体
type CommandParams struct {
	Command     string `json:"command"` // 需要执行的命令
	NeedConfirm bool   `json:"needConfirm"`
}

// Info 返回工具信息，包括名称、描述和参数信息
func (tt *TerminalTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "terminal_tool",
		Desc: "用于执行终端命令并返回结果，仅支持一条命令，支持跨平台（Windows、Linux、macOS）",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"command": {
				Desc:     "需要执行的命令，每次只可以执行一条",
				Type:     "string",
				Required: true,
			},
			"needConfirm": {
				Desc:     "是否需要用户确认是否执行，如果是敏感操作，需要用户确认则为true",
				Type:     "bool",
				Required: true,
			},
		}),
	}, nil
}

// InvokableRun 执行命令并返回结果
func (tt *TerminalTool) InvokableRun(ctx context.Context, argumentsInJson string, opts ...tool.Option) (string, error) {
	fmt.Println(argumentsInJson)
	// 解析输入的JSON参数
	var params CommandParams
	if err := json.Unmarshal([]byte(argumentsInJson), &params); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	// 验证命令参数
	if params.Command == "" {
		return "", fmt.Errorf("命令不能为空")
	}

	// // 如果需要用户确认，则进行交互确认
	// if params.NeedConfirm {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Printf("确认执行敏感命令: \"%s\"? (y/N): ", params.Command)
	// 	input, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		return "", fmt.Errorf("读取用户输入失败: %w", err)
	// 	}
	//
	// 	// 处理输入（去除空白字符并转为小写）
	// 	input = strings.TrimSpace(input)
	// 	if strings.ToLower(input) != "y" && strings.ToLower(input) != "yes" {
	// 		return "", fmt.Errorf("用户已取消命令执行")
	// 	}
	// }

	// 根据操作系统选择合适的shell
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// Windows使用cmd.exe执行命令
		cmd = exec.CommandContext(ctx, "cmd.exe", "/c", params.Command)
	case "linux", "darwin":
		// Linux和macOS使用bash执行命令
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", params.Command)
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput() // 同时捕获stdout和stderr
	if err != nil {
		// 命令执行失败时，也返回输出信息以便调试
		return string(output), nil
	}

	fmt.Println(string(output))

	// 构建返回结果
	result := fmt.Sprintf("命令执行成功:\n命令: %s\n输出:\n%s", params.Command, string(output))
	return result, nil
}
