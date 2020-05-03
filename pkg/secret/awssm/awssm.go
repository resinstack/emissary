package awssm

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/the-maldridge/emissary/pkg/secret"
)

type awssm struct {
	sm *secretsmanager.SecretsManager
}

func init() {
	secret.RegisterProvider("awssm", initialize)
}

func initialize() (secret.Provider, error) {
	s, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	ec2m := ec2metadata.New(s)
	r, err := ec2m.Region()
	if err == nil {
		s.Config.Region = &r
	}

	x := new(awssm)
	x.sm = secretsmanager.New(s)

	return x, nil
}

func (a *awssm) FetchSecret(id string) (string, error) {
	output, err := a.sm.GetSecretValue(
		&secretsmanager.GetSecretValueInput{SecretId: &id},
	)

	if err != nil {
		log.Println(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				return "", secret.ErrNotFound
			default:
				return "", secret.ErrTerminal
			}
		}
		return "", secret.ErrTerminal
	}

	if *output.SecretString == "" {
		return "", secret.ErrNotFound
	}

	return *output.SecretString, nil
}
