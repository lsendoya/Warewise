package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lsendoya/Warewise/internal/product/infrastructure"
	"github.com/lsendoya/Warewise/pkg/aws"
	"github.com/lsendoya/Warewise/pkg/db"
	"github.com/lsendoya/Warewise/pkg/logger"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
	"time"
)

func StartApp(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {

	cfg, errCfg := aws.Config()
	if errCfg != nil {
		logger.Fatalf(errCfg.Error())
	}

	gormDB, errConn := db.NewDataBase(cfg)
	if errConn != nil {
		logger.Fatalf(errConn.Error())
	}

	errMigrate := db.MigrateDB(gormDB)
	if errMigrate != nil {
		logger.Fatalf(errMigrate.Error())
	}

	s3Client := s3.NewFromConfig(cfg)

	return initApiGateway(gormDB, s3Client, ctx, request)
}

func initApiGateway(db *gorm.DB, s3 *s3.Client, ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	r := infrastructure.NewRouter(&request, db, s3)
	prefix := os.Getenv("URL_PREFIX")
	path := strings.Replace(request.RawPath, prefix, "", -1)
	method := request.RequestContext.HTTP.Method

	if strings.HasPrefix(path, "/health") && method == http.MethodGet {
		return health()
	}

	return r.Routes(ctx, request, path, method)
}

func health() (*events.APIGatewayProxyResponse, error) {

	body := map[string]interface{}{
		"time":         time.Now().Format(time.RFC3339),
		"message":      "Hello World!",
		"service_name": "MyService",
	}

	res, _ := json.Marshal(body)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(res),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
