/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-rate-limiter/internal/delivery/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var port int
var redisAddr string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Args:  cobra.NoArgs,
	Run:   runServerCmd,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVarP(&port, "port", "", 3000, "give port number to server")
	serverCmd.Flags().StringVarP(&redisAddr, "redis_addr", "", "localhost:6379", "set Redis address with host:port")
}

func runServerCmd(cmd *cobra.Command, args []string) {
	server := http.NewHttpServer()
	if err := server.Run(fmt.Sprintf(":%d", port)); err != nil {
		logrus.Panicf("server.Run failed, err: %v", err)
	}
}
