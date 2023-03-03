package commands

import (
	"fmt"
	"os"

	"github.com/ohmpatel1997/CommitGPT/pkg/conversation"
	"github.com/ohmpatel1997/CommitGPT/pkg/openai"
	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "The generate command to generate AI commit message",
		Long:  "The generate command to generate AI commit message",
		Run:   generate,
	}
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().Bool("stage", false, "use it to stage the changes as well")
}

func generate(cmd *cobra.Command, args []string) {
	stage, _ := cmd.Flags().GetBool("stage")
	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) == 0 {
		fmt.Println("OPENAI_API_KEY is not set")
		os.Exit(1)
	}

	fmt.Println("testing the ai generated commit message")
	conversation.NewConversation(openai.NewGptClient(apiKey)).StartConversation(stage)
}
