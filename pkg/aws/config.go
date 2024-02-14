package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/labstack/gommon/log"
)

func Config() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithDefaultRegion("us-east-1"),
	)
	if err != nil {
		log.Panic("error loading config .aws/config ", err.Error())
	}
	return cfg, nil
}
