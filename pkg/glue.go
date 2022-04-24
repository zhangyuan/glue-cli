package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/pkg/errors"
)

func NewGlue() (*glue.Glue, error) {
	sess, err := session.NewSession(&aws.Config{
		MaxRetries: aws.Int(3),
	})
	if err != nil {
		return nil, errors.Wrap(err, "fail to create new session")
	}

	svc := glue.New(sess)

	return svc, nil
}
