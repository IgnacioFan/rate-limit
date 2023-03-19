/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-rate-limiter/internal/delivery"

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
	Run: func(cmd *cobra.Command, args []string) {
		server := delivery.NewHttpServer()
		if err := server.Run(fmt.Sprintf(":%d", port)); err != nil {
			logrus.Panicf("server.Run failed, err: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().IntVarP(&port, "port", "", 3000, "give port number to server")
	serverCmd.Flags().StringVarP(&redisAddr, "redis_addr", "", "localhost:6379", "set Redis address with host:port")
}
