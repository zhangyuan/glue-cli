package main

import (
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
	sess := session.Must(session.NewSession(&aws.Config{
		MaxRetries: aws.Int(3),
	}))

	svc := glue.New(sess)

	runsOutput, err := svc.GetJobRuns(&glue.GetJobRunsInput{
		JobName:    aws.String(jobName),
		MaxResults: aws.Int64(10),
	})

	if err != nil {
		return errors.Wrap(err, "fail to batch get jobs")
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Run Id", "State", "Completed On"})

	for index, jr := range runsOutput.JobRuns {
		t.AppendRow([]interface{}{index, *jr.Id, *jr.JobRunState, jr.CompletedOn})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}
