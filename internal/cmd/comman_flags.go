package cmd

import "github.com/spf13/cobra"

var flagCommitMsg string
var flagCommitDescription string
var flagStageAll bool
var flagUseAi bool

func addCommitMsgFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&flagCommitMsg, "message", "m", "", "Commit message")
	cmd.Flags().StringVarP(&flagCommitDescription, "description", "d", "", "Commit description")
	cmd.Flags().BoolVarP(&flagStageAll, "stage-all", "s", false, "Stage all changes")
}

func addUseAiFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&flagUseAi, "ai", "", false, "Use AI to generate git commit messages")
}
