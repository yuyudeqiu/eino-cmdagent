package template

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

const (
	systemPromptTemplate = `
你是一个命令行助手，用户的操作系统为{os}，请帮助用户写出能在其操作系统终端上运行的命令，只需要给出一个最准确的命令，然后查看是否有工具可以执行该命令。
如果是Windows请使用cmd命令，不要使用powershell`
	userPromptTemplate = `请帮助我生成一个命令实现：{command_prompt}`
)

type ChatTemplateConfig struct {
	FormatType schema.FormatType
	Templates  []schema.MessagesTemplate
}

// NewGetCommandTemplate component initialization function of node 'ParsePromptToCommandTemplate' in graph 'cmdAgent'
func NewGetCommandTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	config := &ChatTemplateConfig{
		FormatType: schema.FString,
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage(systemPromptTemplate),
			&schema.Message{
				Role:    schema.User,
				Content: userPromptTemplate,
			},
		},
	}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}
