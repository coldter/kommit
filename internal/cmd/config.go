package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/coldter/kommit/settings"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}

var configCmd = &cobra.Command{Use: "config", Short: "Manage Kommit configuration"}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Kommit configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := settings.ReadSettings()
		if err != nil {
			return err
		}

		configPath := s.GetConfigPath()
		pterm.Info.Printfln("Config path: %s", pterm.Underscore.Sprint(configPath))

		pterm.Println()

		allConfigFromFile := s.GetAlLSettings()
		jsonString, err := json.MarshalIndent(allConfigFromFile, "", "  ")
		fmt.Println(pterm.Bold.Sprint(string(jsonString)))

		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set Kommit configuration",
	Args:  cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return []string{"baseUrl", "token"}, cobra.ShellCompDirectiveNoFileComp
		}
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		validArgs := []string{"baseUrl", "token"}

		s, err := settings.ReadSettings()
		if err != nil {
			return err
		}
		switch args[0] {
		case "baseUrl":
			argUrl := args[1]
			parsedUrl, err := url.Parse(argUrl)
			if err != nil || argUrl == "" {
				return fmt.Errorf("invalid URL: %s", argUrl)
			}
			s.SetBaseUrl(parsedUrl.String())
		case "token":
			argToken := args[1]
			s.SetToken(argToken)
		default:
			return fmt.Errorf("invalid argument: %s must be one of %v", args[0], validArgs)
		}

		err = settings.TryToPersistChanges()
		if err != nil {
			return err
		}

		return nil
	},
}
