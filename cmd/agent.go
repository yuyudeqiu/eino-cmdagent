/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"cmd-agent/internal/agent"

	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取命令行参数
		prompt, _ := cmd.Flags().GetString("prompt")

		cmdAgent := agent.NewCmdAgent()
		err := cmdAgent.Process(context.Background(), prompt)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().StringP("prompt", "p", "", "prompt")
}
