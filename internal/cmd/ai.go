package cmd

import (
	"fmt"
	"github.com/coldter/kommit/internal/utils"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var flagPlainPrint bool

func init() {
	rootCmd.AddCommand(aiCmd)

	aiCmd.Flags().BoolVarP(&flagPlainPrint, "plain", "p", false, "Just print the commit message without commiting")
	aiCmd.Flags().BoolVarP(&flagStageAll, "stage-all", "s", false, "Stage all changes")
}

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "Non-interactive cmd to generate commit msg from diff that doesn't prompt for everything",
	RunE:  aiCmdRunE,
}

func aiCmdRunE(_ *cobra.Command, _ []string) error {
	pterm.Info.Println("Non-interactive [Ai] 🤖 mode")

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	gitRepo, err := utils.OpenGitRepo(dir)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return errors.Errorf("%s is not a git repository", dir)
		}
		return err
	}

	if flagStageAll {
		if err = gitRepo.AddAll(); err != nil {
			return errors.Wrap(err, "failed to stage changes")
		}
		pterm.Success.Printfln("Changes staged successfully in %s", pterm.Underscore.Sprint(gitRepo.RootDir))
	}
	hasChanges, err := gitRepo.HasStagedChanges()
	if err != nil {
		return errors.Wrap(err, "failed to check for staged changes")
	}

	if !hasChanges {
		pterm.Info.Println("No staged changes found")
		return nil
	}

	diff, err := gitRepo.GitDiff()
	if err != nil {
		return errors.Wrap(err, "failed to get git diff")
	}

	client, err := utils.GetClient()
	if err != nil {
		return err
	}

	// strip large diff to avoid model context window exertion
	var diffToSend string
	if len(diff) > 25000 {
		diffToSend = diff[:25000] + "..."
	} else {
		diffToSend = diff
	}

	spinnerInfo, _ := pterm.DefaultSpinner.Start("Requesting [Ai] 🤖 for commit generation...")

	generated, err := client.Ai.GenerateGitCommitMsgFromDiff(strings.Join(
		[]string{
			// models dont like this
			//"Last 10 commits messages",
			//previousCommitMsgs,
			//"------------------",
			diffToSend,
		}, "\n"))
	if err != nil {
		return err
	}

	if generated.Message != "" {
		spinnerInfo.Success("Generated commit message successfully ")
	} else {
		spinnerInfo.Warning("Couldn't generate commit message correctly")
		return nil
	}

	if flagPlainPrint {
		pterm.DefaultSection.Println("Generated Git Diffs::")
		pterm.Info.Println("Commit Message:")
		pterm.Bold.Printfln(generated.Message)
		pterm.Info.Println("Commit Description:")
		pterm.Bold.Printfln(generated.Body)

		pterm.Println()
		return nil
	}

	pterm.Info.Printfln("Commit msg:: %s", pterm.Underscore.Sprint(generated.Message))
	if generated.Body != "" {
		pterm.Info.Printfln("Commit Description:: \n%s", pterm.Underscore.Sprint(generated.Body))
	}

	cnfmCommit, _ := pterm.DefaultInteractiveConfirm.Show(
		fmt.Sprintf("Confirm git commit"),
	)
	if !cnfmCommit {
		pterm.Warning.Println("Aborting git commit")
		return nil
	}
	hash, err := gitRepo.Commit(generated.Message, generated.Body)
	if err != nil {
		return errors.Wrap(err, "failed to commit")
	}
	pterm.Success.Printfln("Changes committed successfully commit:: %s", pterm.Underscore.Sprint(hash))

	return nil
}
