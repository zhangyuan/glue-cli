package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("job name is not provided")
	}

	jobName := os.Args[1]

	if err := run(jobName); err != nil {
		log.Fatalln(err.Error())
	}
}

func run(jobName string) error {
	sess, err := session.NewSession(&aws.Config{
		MaxRetries: aws.Int(3),
	})
	if err != nil {
		return err
	}

	svc := glue.New(sess)

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
	t.AppendHeader(table.Row{"#", "Run Id", "State", "Completed On"})

	for index, jr := range runsOutput.JobRuns {
		t.AppendRow([]interface{}{index, *jr.Id, *jr.JobRunState, jr.CompletedOn.Local()})
		t.AppendSeparator()
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}
