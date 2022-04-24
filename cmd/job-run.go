package cmd

import (
	"encoding/json"
	"glue/pkg"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var jobRunCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		jobName := args[0]
		arguments := args[1]
		if err := runJob(jobName, arguments); err != nil {
			log.Default().Fatal(err)
		}
	},
}

func runJob(jobName string, jsonArguments string) error {
	svc, err := pkg.NewGlue()
	if err != nil {
		return errors.Wrap(err, "fail to new glue")
	}

	var arguments map[string]string
	if err := json.Unmarshal([]byte(jsonArguments), &arguments); err != nil {
		return errors.Wrapf(err, "fail to parse arguments %s", jsonArguments)
	}

	startJobRun, err := svc.StartJobRun(&glue.StartJobRunInput{
		JobName:   aws.String(jobName),
		Arguments: aws.StringMap(arguments),
	})
	if err != nil {
		return errors.Wrap(err, "fail to start job run")
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRow([]interface{}{"Job", jobName})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"JobRunId", *startJobRun.JobRunId})
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}

func init() {
	jobCmd.AddCommand(jobRunCmd)
}
