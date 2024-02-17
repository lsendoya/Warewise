package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lsendoya/Warewise/pkg/logger"
)

func main() {
	defer logger.SyncLogger()
	lambda.Start(StartApp)

}
