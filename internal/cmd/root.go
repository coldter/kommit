package cmd

import (
	"fmt"
	"github.com/coldter/kommit/internal/utils"
	"github.com/coldter/kommit/settings"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:           "kommit",
	Version:       "dev",
	RunE:          runCmdInteractive,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(err)
	}
}

func init() {
	addCommitMsgFlag(rootCmd)
	addUseAiFlag(rootCmd)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		settings.PersistChanges()
	}
}

// Interactive root command
func runCmdInteractive(_ *cobra.Command, _ []string) error {
	_, _ = settings.ReadSettings()
	pterm.Info.Println("Interactive mode")

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
	pterm.Success.Printfln("Git repository found on %s", pterm.Underscore.Sprint(gitRepo.RootDir))

	// ask if user wants to perform git add .
	stageAll := flagStageAll
	if !flagStageAll {
		// only ask if flag is not set or commit msg flag is not set
		if flagCommitMsg == "" {
			stageAll, _ = pterm.DefaultInteractiveConfirm.Show(
				fmt.Sprintf("Would you like to %s ?", pterm.Red("stage[git add] every change in this repository")),
			)
		}
	}

	if stageAll {
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

	pterm.DefaultSection.Println("Git diff")
	pterm.Println(diff)

	var previousCommitMsgs string
	last10Commits, err := gitRepo.GetLastNCommitMsg(10)
	if err != nil {
		if !errors.Is(err, utils.ErrNoCommit) {
			return errors.Wrap(err, "failed to get last 10 commits")
		}
	} else {
		previousCommitMsgs = strings.Join(last10Commits, "\n")
	}

	pterm.DefaultSection.Println("Last 10 commits")
	pterm.Println(previousCommitMsgs)

	var commitMsg string
	var commitDescription string

	// hack if ai flag and commit flag is also set override the msg flag
	if flagUseAi {
		flagCommitMsg = ""
		commitDescription = ""
	}
	if flagCommitMsg != "" {
		commitMsg = flagCommitMsg
	}
	if flagCommitDescription != "" {
		commitDescription = flagCommitDescription
	}

	// manual commit msg
	prefixOptions := []string{
		"none",
		"fix",
		"feat",
		"perf",
		"refactor",
		"revert",
		"style",
		"test",
		"chore",
		"docs",
	}

	// ask for ai use if flag not set
	if !flagUseAi {
		pterm.Warning.Println("Using Ai to generate commit message will send git diff to remote rpc service")
		flagUseAi, _ = pterm.DefaultInteractiveConfirm.Show(
			fmt.Sprintf("Would you like to use %s to generate git commit message from diff?", pterm.Red(" [Ai] 🤖 ")),
		)
	}

	var commitPrefix string
	if !flagUseAi {
		commitPrefix, _ = pterm.DefaultInteractiveSelect.
			WithOptions(prefixOptions).
			Show("Commit Prefix")

		if commitPrefix == "none" {
			commitPrefix = ""
		}
	}

	var generatedCommitMsg string
	var generatedCommitDescription string
	if flagUseAi {
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
			generatedCommitMsg = generated.Message
			generatedCommitDescription = generated.Body
		} else {
			spinnerInfo.Warning("Couldn't generate commit message correctly")
		}
	}

	if flagCommitMsg == "" {
		commitMsgInteractiveTextInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)

		if generatedCommitMsg != "" {
			commitMsgInteractiveTextInput = commitMsgInteractiveTextInput.WithDefaultValue(generatedCommitMsg)
		}

		commitMsg, _ = commitMsgInteractiveTextInput.Show("Commit Message")
	}

	if commitPrefix != "" {
		commitMsg = commitPrefix + ": " + commitMsg
	}

	if flagCommitMsg == "" {
		commitDescriptionInteractiveTextInput := pterm.DefaultInteractiveTextInput.
			WithMultiLine(true)

		if generatedCommitDescription != "" {
			commitDescriptionInteractiveTextInput = commitDescriptionInteractiveTextInput.WithDefaultValue(generatedCommitDescription)
		}

		commitDescription, _ = commitDescriptionInteractiveTextInput.Show("Commit Description")
	}

	pterm.Info.Printfln("Commit msg:: %s", pterm.Underscore.Sprint(commitMsg))
	if commitDescription != "" {
		pterm.Info.Printfln("Commit Description:: \n%s", pterm.Underscore.Sprint(commitDescription))
	}

	cnfmCommit, _ := pterm.DefaultInteractiveConfirm.Show(
		fmt.Sprintf("Confirm git commit"),
	)
	if !cnfmCommit {
		pterm.Warning.Println("Aborting git commit")
		return nil
	}
	hash, err := gitRepo.Commit(commitMsg, commitDescription)
	if err != nil {
		return errors.Wrap(err, "failed to commit")
	}
	pterm.Success.Printfln("Changes committed successfully commit:: %s", pterm.Underscore.Sprint(hash))

	return nil
}
