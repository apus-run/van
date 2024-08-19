package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/apus-run/van/cmd/van/config"
	"github.com/apus-run/van/cmd/van/internal/new"
)

var RootCmd = &cobra.Command{
	Use:          "van",
	Example:      "van <command> <subcommand> [flags]",
	Short:        "Van CLI",
	Long:         "+-------------------------------------------+\n|     █████   █████                         |\n|    ░░███   ░░███                          |\n|     ░███    ░███   ██████   ████████      |\n|     ░███    ░███  ░░░░░███ ░░███░░███     |\n|     ░░███   ███    ███████  ░███ ░███     |\n|      ░░░█████░    ███░░███  ░███ ░███     |\n|        ░░███     ░░████████ ████ █████    |\n|         ░░░       ░░░░░░░░ ░░░░ ░░░░░     |\n+-------------------------------------------+\nVan: 一个轻量级的Golang应用搭建脚手架",
	Version:      fmt.Sprintf("\n__     __              \n\\ \\   / /  __ _  _ __  \n \\ \\ / /  / _` || '_ \\ \n  \\ V /  | (_| || | | |\n   \\_/    \\__,_||_| |_|\n \nVan %s - Copyright (c) 2024-2026 Van\nReleased under the MIT License.\n\n", config.Release),
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

func init() {
	// RootCmd.AddCommand(rpc.Cmd)
	RootCmd.AddCommand(new.Cmd)
	// RootCmd.AddCommand(run.Cmd)
	// RootCmd.AddCommand(upgrade.Cmd)
}
func main() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
