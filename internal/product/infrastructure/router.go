package infrastructure

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lsendoya/Warewise/internal/product/application"
	"github.com/lsendoya/Warewise/internal/product/infrastructure/http"
	"github.com/lsendoya/Warewise/internal/product/infrastructure/storage"
	"github.com/lsendoya/Warewise/pkg/auth"
	"github.com/lsendoya/Warewise/pkg/aws"
	"github.com/lsendoya/Warewise/pkg/errors"
	"github.com/lsendoya/Warewise/pkg/tools"
	"gorm.io/gorm"
	net "net/http"
	"strings"
)

type Route struct {
	e   *events.APIGatewayV2HTTPRequest
	h   http.Handler
	res errors.Response
}

func NewRouter(e *events.APIGatewayV2HTTPRequest, db *gorm.DB, s3 *s3.Client) Route {
	s := storage.New(db)
	awsScv := aws.NewAWSService(s3)
	uc := application.New(&s, awsScv)
	h := http.NewHandler(&uc)

	return Route{
		e: e,
		h: h,
	}

}

func (r *Route) Routes(ctx context.Context, request events.APIGatewayV2HTTPRequest, path, method string) (*events.APIGatewayProxyResponse, error) {
	token := request.Headers["authorization"]

	if strings.HasPrefix(path, "/api/v1/admin/products") {
		jwtParsed, err := auth.ParseJWT(token)
		if err != nil {
					return r.res.Unauthorized(err)
		}


		if !tools.Contains(jwtParsed.CognitoGroups, "admin") {
			return r.res.Forbidden()
		}

		switch method {
		case net.MethodPost:
			return r.h.Add(ctx, request)
		case net.MethodPut:
			return r.h.Update(ctx, request)
		case net.MethodDelete:
			return r.h.Delete(ctx, request)
		default:
			return &events.APIGatewayProxyResponse{StatusCode: net.StatusMethodNotAllowed, Body: "Method not allowed"}, nil
		}
	}

	id := request.PathParameters["id"]

	if id != "" && method == net.MethodGet {
		return r.h.Get(ctx, request)
	}

	if path == "/api/v1/public/products" {
		return r.h.List(ctx, request)
	}

	return &events.APIGatewayProxyResponse{StatusCode: net.StatusNotFound, Body: "Not Found"}, nil
}
