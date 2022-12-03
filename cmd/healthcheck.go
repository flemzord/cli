package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/appclacks/cli/client"
	apitypes "github.com/appclacks/go-types"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

func toMap(slice []string) (map[string]string, error) {
	sliceMap := make(map[string]string)
	for _, l := range slice {
		splitted := strings.Split(l, "=")
		if len(splitted) != 2 {
			return nil, fmt.Errorf("Invalid label %s", l)
		}
		sliceMap[splitted[0]] = splitted[1]
	}
	return sliceMap, nil

}

func getHealthcheckCmd(client *client.Client) *cobra.Command {
	var healthcheckID string
	var getHealthcheck = &cobra.Command{
		Use:   "get",
		Short: "Get an API healthcheck",
		Run: func(cmd *cobra.Command, args []string) {
			input := apitypes.GetHealthcheckInput{
				ID: healthcheckID,
			}
			healthcheck, err := client.GetHealthcheck(input)
			exitIfError(err)
			if outputFormat == "json" {
				json, err := json.Marshal(healthcheck)
				exitIfError(err)
				fmt.Println(string(json))
				os.Exit(0)
			}
			t := tabby.New()

			t.AddHeader("ID", "Name", "Description", "Interval", "Timeout", "Labels", "Enabled", "Definition")
			jsonLabels, err := json.Marshal(healthcheck.Labels)
			exitIfError(err)
			jsonDef, err := json.Marshal(healthcheck.Definition)
			exitIfError(err)
			t.AddLine(healthcheck.ID, healthcheck.Name, healthcheck.Description, healthcheck.Interval, healthcheck.Timeout, string(jsonLabels), healthcheck.Enabled, string(jsonDef))

			t.Print()
			os.Exit(0)
		},
	}
	getHealthcheck.PersistentFlags().StringVar(&healthcheckID, "id", "", "Healthcheck ID")
	err := getHealthcheck.MarkPersistentFlagRequired("id")
	exitIfError(err)

	return getHealthcheck
}

func deleteHealthcheckCmd(client *client.Client) *cobra.Command {
	var tokenID string
	var deleteHealthcheck = &cobra.Command{
		Use:   "delete",
		Short: "Delete an healthcheck",
		Run: func(cmd *cobra.Command, args []string) {
			input := apitypes.DeleteHealthcheckInput{
				ID: tokenID,
			}
			result, err := client.DeleteHealthcheck(input)
			exitIfError(err)
			if outputFormat == "json" {
				json, err := json.Marshal(result)
				exitIfError(err)
				fmt.Println(string(json))
				os.Exit(0)
			}
			t := tabby.New()
			t.AddHeader("Messages")
			for _, message := range result.Messages {
				t.AddLine(message)
			}
			t.Print()
			os.Exit(0)
		},
	}
	deleteHealthcheck.PersistentFlags().StringVar(&tokenID, "id", "", "Token ID")
	err := deleteHealthcheck.MarkPersistentFlagRequired("id")
	exitIfError(err)

	return deleteHealthcheck
}

func listHealthchecksCmd(client *client.Client) *cobra.Command {
	var listHealthchecks = &cobra.Command{
		Use:   "list",
		Short: "List API healthcheck",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := client.ListHealthchecks()
			exitIfError(err)
			if outputFormat == "json" {
				json, err := json.Marshal(result)
				exitIfError(err)
				fmt.Println(string(json))
				os.Exit(0)
			}
			t := tabby.New()

			t.AddHeader("ID", "Name", "Description", "Interval", "Timeout", "Labels", "Enabled", "Definition")
			for _, healthcheck := range result.Result {
				jsonLabels, err := json.Marshal(healthcheck.Labels)
				exitIfError(err)
				jsonDef, err := json.Marshal(healthcheck.Definition)
				exitIfError(err)
				t.AddLine(healthcheck.ID, healthcheck.Name, healthcheck.Description, healthcheck.Interval, healthcheck.Timeout, string(jsonLabels), healthcheck.Enabled, string(jsonDef))
			}

			t.Print()
			os.Exit(0)
		},
	}

	return listHealthchecks
}