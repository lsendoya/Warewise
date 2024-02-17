package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/lsendoya/Warewise/pkg/logger"
	"net/http"
)

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type MessageResponse struct {
	Data     interface{} `json:"data"`
	Errors   Response    `json:"errors"`
	Messages Response    `json:"messages"`
}

func headers() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}

func (r Response) OK(data interface{}) (*events.APIGatewayProxyResponse, error) {
	logger.Info("status: %v", MessageOK)
	response := returnResponse(MessageOK, nil, data)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
		Headers:    headers(),
	}, nil

}

func (r Response) Updated(data interface{}) (*events.APIGatewayProxyResponse, error) {
	logger.Info("status: %v", MessageOK)
	response := returnResponse("Entity updated successfully", nil, data)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
		Headers:    headers(),
	}, nil
}

func (r Response) Deleted(id uuid.UUID) (*events.APIGatewayProxyResponse, error) {
	logger.Info("status: %v", MessageOK)
	msg := fmt.Sprintf("Entity with id %v deleted successfully", id)
	response := returnResponse(msg, nil, nil)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
		Headers:    headers(),
	}, nil
}

func (r Response) Forbidden() (*events.APIGatewayProxyResponse, error) {
	logger.Errorf("Access denied. Admin permissions required  status:%v", MessageForbidden)
	err := errors.New("access denied. Admin permissions required")

	response := returnResponse(MessageUnauthorized, err, nil)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusForbidden,
		Body:       string(response),
		Headers:    headers(),
	}, nil
}

func (r Response) Unauthorized(err error) (*events.APIGatewayProxyResponse, error) {
	logger.Errorf("error: %v, status:%v", err.Error(), MessageUnauthorized)
	response := returnResponse(MessageUnauthorized, err, nil)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       string(response),
		Headers:    headers(),
	}, nil
}

func (r Response) BadRequest(err error) (*events.APIGatewayProxyResponse, error) {
	logger.Errorf("error: %v, status:%v", err.Error(), MessageBadRequest)
	response := returnResponse(MessageBadRequest, err, nil)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       string(response),
		Headers:    headers(),
	}, nil
}

func (r Response) Created(data interface{}) (*events.APIGatewayProxyResponse, error) {
	logger.Info("status: %v", MessageCreated)
	response := returnResponse(MessageCreated, nil, data)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(response),
		Headers:    headers(),
	}, nil
}

func (r Response) NotFound(err error) (*events.APIGatewayProxyResponse, error) {
	logger.Errorf("error: %v, status:%v", err.Error(), MessageNotFound)

	response := returnResponse(MessageNotFound, err, nil)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       string(response),
		Headers:    headers(),
	}, nil
}

func (r Response) InternalServerError(who string, err error) (*events.APIGatewayProxyResponse, error) {
	e := NewError()
	e.Err = err
	e.APIMessage = "internal error"
	e.Code = MessageInternalServerError
	e.StatusHTTP = http.StatusInternalServerError
	e.Who = who

	logger.Errorf("error: %v, status:%v", e.Error(), MessageInternalServerError)

	response := returnResponse(MessageInternalServerError, e.Err, nil)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       string(response),
		Headers:    headers(),
	}, nil
}

func returnResponse(msg string, err error, data interface{}) []byte {
	var message MessageResponse
	if err != nil {
		message = MessageResponse{
			Data: nil,
			Errors: Response{
				Code: msg, Message: err.Error(),
			},
			Messages: Response{},
		}
	} else {
		message = MessageResponse{
			Data:   data,
			Errors: Response{},
			Messages: Response{
				Code: msg, Message: msg,
			},
		}
	}

	response, _ := json.Marshal(message)
	return response
}
