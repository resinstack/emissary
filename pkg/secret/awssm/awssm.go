package awssm

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/resinstack/emissary/pkg/secret"
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

	log.Println("AWS SM engine is initialized and available")
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

	scrt := strings.TrimSpace(*output.SecretString)
	if scrt == "" || scrt == "intentionally-empty" {
		return "", secret.ErrNotFound
	}

	log.Printf("Secret %s aquired", id)
	return scrt, nil
}
