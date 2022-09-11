/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/bootstrap"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "smss",
	Short: "Serves sessions for Squadmortar.xyz™",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrapper, err := bootstrap.New(
			host,
			port,
			privateKeyFilepath,
			websocketKeepAlive,
			allowedOrigins,
			enablePlayground,
			enableIntrospection)
		if err != nil {
			panic(err)
		}

		if err := bootstrapper.Listen(); err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var host string
var port uint16
var privateKeyFilepath string
var websocketKeepAlive time.Duration
var allowedOrigins []string
var enablePlayground bool
var enableIntrospection bool

func init() {
	rootCmd.Flags().StringVar(&host, "host", "localhost", "Listener host for the GraphQL server.")
	rootCmd.Flags().Uint16Var(&port, "port", 8080, "Listener port for the GraphQL server.")
	rootCmd.Flags().StringVarP(&privateKeyFilepath, "private-key-file", "k", "./.secrets/token.key.pem", "Private key file to create authentication tokens with.")
	rootCmd.Flags().DurationVar(&websocketKeepAlive, "websocket-keep-alive", time.Second*10, "Websocket keep-alive interval")
	rootCmd.Flags().StringSliceVar(&allowedOrigins, "allowed-origins", []string{"*"}, "Allowed origin for CORS")
	rootCmd.Flags().BoolVar(&enablePlayground, "enable-playground", false, "Enables the GraphiQL playground under /graphql/playground. Only available for HTTP")
	rootCmd.Flags().BoolVar(&enableIntrospection, "enable-introspection", false, "Enables introspection for GraphQL responses. Disable it in production.")
}
