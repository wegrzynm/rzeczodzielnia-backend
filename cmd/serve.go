/*
Copyright © 2024 Mateusz Węgrzyn matzyn.yt@gmail.com
*/
package cmd

import (
	"Rzeczodzielnia/internal/server"
	"fmt"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server",
	Long:  `Start HTTP server`,
	Run:   startServer,
}

func startServer(cmd *cobra.Command, args []string) {
	s := server.NewServer()

	err := s.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
