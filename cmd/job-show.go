package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"glue/pkg"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var jobShowCmd = &cobra.Command{
	Use: "show",
	Run: func(cmd *cobra.Command, args []string) {
		jobName := args[0]
		if err := showJob(jobName); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	jobCmd.AddCommand(jobShowCmd)
}

func showJob(jobName string) error {
	svc, err := pkg.NewGlue()
	if err != nil {
		return err
	}

	jobOutput, err := svc.GetJob(&glue.GetJobInput{
		JobName: aws.String(jobName),
	})
	if err != nil {
		return errors.Wrap(err, "fail to get job")
	}

	renderJob(jobOutput)

	runsOutput, err := svc.GetJobRuns(&glue.GetJobRunsInput{
		JobName:    aws.String(jobName),
		MaxResults: aws.Int64(10),
	})

	renderJobRuns(*runsOutput)

	if err != nil {
		return errors.Wrap(err, "fail to batch get jobs")
	}
	return nil
}

func renderJob(jobOutput *glue.GetJobOutput) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRow([]interface{}{"Name", *jobOutput.Job.Name})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"Description", *jobOutput.Job.Description})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"CreatedOn", jobOutput.Job.CreatedOn.Local()})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"LastModifiedOn", jobOutput.Job.LastModifiedOn.Local()})
	t.AppendSeparator()

	var buffer bytes.Buffer
	for key, element := range jobOutput.Job.DefaultArguments {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", key, *element))
	}

	t.AppendRow([]interface{}{"DefaultArguments", buffer.String()})

	t.SetStyle(table.StyleLight)
	t.Render()
}

func renderJobRuns(runsOutput glue.GetJobRunsOutput) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Run Id", "State", "Started On", "Completed On"})

	for index, jr := range runsOutput.JobRuns {
		if jr.CompletedOn != nil {
			t.AppendRow([]interface{}{index, *jr.Id, *jr.JobRunState, jr.StartedOn.Local(), jr.CompletedOn.Local()})
		} else {
			t.AppendRow([]interface{}{index, *jr.Id, *jr.JobRunState, jr.StartedOn.Local(), ""})
		}

		t.AppendSeparator()
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}
