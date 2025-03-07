package cmd

import (
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

var cliDescriptionMsg string

func init() {
	if msg, err := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("K", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithStyle("ommit", pterm.FgLightMagenta.ToStyle()),
	).Srender(); err == nil {
		cliDescriptionMsg = msg + pterm.DefaultBasicText.Sprint(pterm.LightMagenta("Kommit")+" is a tool to generate git commit msg.")
	} else {
		cliDescriptionMsg = "Kommit is a tool to generate git commit msg."
	}

	rootCmd.Long = cliDescriptionMsg
}
