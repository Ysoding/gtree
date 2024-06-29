package cmd

import (
	"fmt"
	"os"

	gtree "github.com/Ysoding/gtree/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gtree",
	Short: "a tree clone",
	Long:  `a tree clone`,
	Run:   gtree.Run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
