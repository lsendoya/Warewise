package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/lsendoya/Warewise/internal/product/domain"
	errHandle "github.com/lsendoya/Warewise/pkg/errors"
	"github.com/lsendoya/Warewise/pkg/logger"
	"gorm.io/gorm"
	"mime"
	"mime/multipart"
	"strconv"
)

type Handler struct {
	uc  domain.UseCase
	res errHandle.Response
}

func NewHandler(uc domain.UseCase) Handler {

	return Handler{uc: uc}
}

func (h Handler) Add(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	reader, errRead := readMultipart(request)
	if errRead != nil {
		return h.res.BadRequest(errRead)
	}

	form, errReader := reader.ReadForm(32 << 20)
	if errReader != nil {
		return h.res.InternalServerError("Error reading form", errReader)
	}

	defer func(form *multipart.Form) {
		err := form.RemoveAll()
		if err != nil {

		}
	}(form)

	formData, errExtract := extractForm(form)
	if errRead != nil {
		return h.res.BadRequest(errExtract)
	}

	product, errCreate := h.uc.Add(&formData)
	if errCreate != nil {
		logger.Errorf("h.uc.Add", errCreate)
		return h.res.InternalServerError("Error creating product", errCreate)
	}

	return h.res.Created(product)

}

func (h Handler) Update(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	var payload domain.Product
	if err := json.Unmarshal([]byte(request.Body), &payload); err != nil {
		return h.res.BadRequest(err)
	}
	logger.Info("payload -handler", payload)
	id, err := uuid.Parse(request.PathParameters["id"])
	if err != nil {
		return h.res.BadRequest(err)
	}

	product, errUpdated := h.uc.Update(id, &payload)
	if errUpdated != nil {
		logger.Errorf("h.uc.Update(id, payload)-->", errUpdated)
		if errors.Is(errUpdated, gorm.ErrRecordNotFound) {
			return h.res.NotFound(errUpdated)
		}
		return h.res.InternalServerError("h.uc.Update()", errUpdated)
	}

	return h.res.Updated(*product)

}
func (h Handler) Delete(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	id, err := uuid.Parse(request.PathParameters["id"])
	if err != nil {
		return h.res.BadRequest(err)
	}

	errDeleted := h.uc.Delete(id)
	if errDeleted != nil {
		logger.Errorf("h.uc.Delete(id)-->", errDeleted)
		if errors.Is(errDeleted, gorm.ErrRecordNotFound) {
			return h.res.NotFound(errDeleted)
		}
		return h.res.InternalServerError("h.uc.Delete()", errDeleted)
	}

	return h.res.Deleted(id)
}

func (h Handler) Get(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	id, errParse := uuid.Parse(request.PathParameters["id"])
	if errParse != nil {
		return h.res.BadRequest(errParse)
	}

	product, errGet := h.uc.Get(id)
	if errGet != nil {
		logger.Errorf("h.uc.Get(id)-->", errGet)
		if errors.Is(errGet, gorm.ErrRecordNotFound) {
			return h.res.NotFound(errGet)
		}
		return h.res.InternalServerError("h.uc.Get()", errGet)
	}

	return h.res.OK(*product)
}

func (h Handler) List(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	products, errList := h.uc.List()
	if errList != nil {
		logger.Errorf("h.uc.List()-->", errList)
		return h.res.InternalServerError("h.uc.List()", errList)
	}

	return h.res.OK(products)
}

func readMultipart(request events.APIGatewayV2HTTPRequest) (*multipart.Reader, error) {
	if request.Body == "" {
		return nil, errors.New("no body in the request")
	}

	body, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return nil, err
	}

	contentType := request.Headers["content-type"]

	_, params, errParser := mime.ParseMediaType(contentType)
	if errParser != nil {
		return nil, errParser
	}

	return multipart.NewReader(bytes.NewReader(body), params["boundary"]), nil
}

func extractForm(form *multipart.Form) (domain.FormDataProduct, error) {
	var formData domain.FormDataProduct

	formData.Name = form.Value["name"][0]
	formData.Description = form.Value["description"][0]
	formData.Color = form.Value["color"][0]

	var sizes json.RawMessage
	errMarshall := json.Unmarshal([]byte(form.Value["size"][0]), &sizes)
	if errMarshall != nil {
		return domain.FormDataProduct{}, errMarshall
	}
	formData.Size = sizes

	if price, err := strconv.ParseFloat(form.Value["price"][0], 64); err == nil {
		formData.Price = price
	} else {
		return domain.FormDataProduct{}, err
	}

	if Available, err := strconv.ParseBool(form.Value["available"][0]); err == nil {
		formData.Available = Available
	} else {
		return domain.FormDataProduct{}, err
	}

	if files, ok := form.File["imageFiles"]; ok && len(files) > 0 {
		formData.ImageFiles = files
	}

	return formData, nil
}
