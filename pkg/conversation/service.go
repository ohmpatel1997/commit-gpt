package conversation

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ohmpatel1997/CommitGPT/pkg/openai"
	"github.com/ohmpatel1997/CommitGPT/pkg/utils"
	"github.com/pkg/errors"
)

type Conversation struct {
	Client *openai.Client
}

func NewConversation(client *openai.Client) *Conversation {
	return &Conversation{
		Client: client,
	}
}

func (c *Conversation) StartConversation(stage bool) {
	if stage {
		err := stageAllFiles()
		if err != nil {
			fmt.Println(errors.Wrap(err, "failed to stage files	"))
			os.Exit(1)
		}
	}

	// prepare the diff
	diff, err := getDiff()
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get diff"))
		os.Exit(1)
	}

	if len(diff) == 0 {
		fmt.Println("Nothing to commit")
		os.Exit(0)
	}

	var messages = []*openai.Message{
		{
			Role:    "system",
			Content: `Write a commit message for this git files difference, only response the message, no need prefix`,
		},
	}

	commitMessage := ""
	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

		messages = append(messages, &openai.Message{
			Role:    "user",
			Content: diff,
		})

		commitMessage, err = c.Client.ChatComplete(ctx, messages)
		if err != nil {
			fmt.Println("failed to generate commit message: " + err.Error())
			os.Exit(1)
		}

		if commitMessage == "" {
			fmt.Println("I can not understand the message, please try again")
			os.Exit(1)
		} else {
			fmt.Print("\n")
			color.New(color.FgHiRed).Print("Message: ")
			color.New(color.FgHiGreen).Print(commitMessage + "\n")
		}

		userRequest := ""
		color.New(color.FgHiWhite).Print("\n\nEnter")
		color.New(color.FgHiYellow).Print(" `yes` ")
		color.New(color.FgHiWhite).Print("if you want to use the message or press")
		color.New(color.FgHiYellow).Print(" Ctrl+C ")
		color.New(color.FgHiWhite).Print("to exit.\n")

		fmt.Println("you can also ask to improve the message")
		for {
			color.New(color.FgHiWhite).Print("\nInput: ")
			reader := bufio.NewReader(os.Stdin)
			userRequest, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("failed to read user input: " + err.Error())
				os.Exit(1)
			}

			userRequest = strings.TrimSpace(userRequest)
			if len(userRequest) == 0 {
				continue
			}

			break
		}

		if userRequest == "yes" {
			break
		}

		if commitMessage != "" {
			messages = append(messages, &openai.Message{
				Role:    "assistant",
				Content: commitMessage,
			})
		} else {
			messages[len(messages)-1].Content = userRequest
		}
	}

	if err := commit(commitMessage); err != nil {
		fmt.Println("failed to commit: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("Commit successfully with message: " + commitMessage)
}

func commit(message string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	return utils.Run("git",
		utils.WithDir(workingDir),
		utils.WithArgs("commit", "-m", message),
	)
}

func getDiff() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	out := strings.Builder{}
	utils.Run("git",
		utils.WithDir(workingDir),
		utils.WithArgs("diff", "--cached", "--unified=0"),
		utils.WithStdOut(&out),
	)

	return strings.TrimSpace(out.String()), nil
}

func stageAllFiles() error {
	out := strings.Builder{}
	err := utils.Run("git",
		utils.WithArgs("add", "."),
		utils.WithStdOut(&out),
	)
	return err
}
