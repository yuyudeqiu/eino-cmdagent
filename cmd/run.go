/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"runtime"

	"cmd-agent/internal/agent"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "处理自然语言指令并执行",
	Long:  `接收自然语言输入，通过Agent处理为命令行指令并执行，支持确认步骤`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取命令行参数
		prompt, _ := cmd.Flags().GetString("prompt")

		if prompt == "" {
			fmt.Println("请使用 --prompt 或 -p 指定自然语言指令")
			return
		}

		cmdChain, err := agent.NewCmdChain(context.Background())
		if err != nil {
			panic(err)
		}

		process, err := cmdChain.Process(context.Background(), agent.CmdParam{
			OS:     runtime.GOOS,
			Prompt: prompt,
		})
		if err != nil {
			fmt.Println(err)
		}

		println(process)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("prompt", "p", "", "自然语言指令（必填）")
	runCmd.Flags().BoolP("no-confirm", "n", false, "跳过执行前确认步骤")
	runCmd.Flags().BoolP("verbose", "v", false, "显示详细处理过程")

	// 标记prompt为必填参数
	_ = runCmd.MarkFlagRequired("prompt")
}
