package memrecords

import (
	"errors"
	apierrors "github.com/fufu-yedek/getir-challange/apihelper/errors"
	"github.com/fufu-yedek/getir-challange/apihelper/response"
	"github.com/fufu-yedek/getir-challange/gerrors"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	CreateOrUpdate(params CreateOrUpdateParams) (response.Responder, error)
	Retrieve(params RetrieveParams) (response.Responder, error)
}

type controller struct {
	repository Repository
}

func (c controller) CreateOrUpdate(params CreateOrUpdateParams) (response.Responder, error) {
	// swagger:route POST /in-memory InMemory CreateOrUpdate
	//
	// create or update key-value
	//
	//     Produces:
	//     - application/json
	//
	//     Consumes:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       200: inMemoryResponse
	//       400: inMemoryErrorResponse
	//       500: inMemoryErrorResponse

	logger := logrus.WithFields(logrus.Fields{
		"location": "Controller - CreateOrUpdate",
		"params":   params.Body,
	})

	record, err := c.repository.CreateOrUpdate(Record{
		Key:   params.Body.Key,
		Value: params.Body.Value,
	})

	if err != nil {
		logger.WithError(err).Error("error while creating new record")
		return nil, apierrors.ErrInternalServer
	}

	return RecordSerializer{Record: record}, nil
}

func (c controller) Retrieve(params RetrieveParams) (response.Responder, error) {
	// swagger:route GET /in-memory InMemory Retrieve
	//
	// retrieve key-value pair
	//
	//     Produces:
	//     - application/json
	//
	//     Consumes:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       200: inMemoryResponse
	//       400: inMemoryErrorResponse
	//       500: inMemoryErrorResponse

	logger := logrus.WithFields(logrus.Fields{
		"location": "Controller - Retrieve",
		"params":   params,
	})

	record, err := c.repository.FindOne(Filter{
		Key: params.Key,
	})

	if err != nil {
		if errors.Is(err, gerrors.ErrRecordNotFound) {
			logrus.WithError(err).Warn("could not find record")
			return nil, apierrors.NewErrUserReadable("could not find record")
		}

		logger.WithError(err).Error("error while finding record")
		return nil, apierrors.ErrInternalServer
	}

	return RecordSerializer{Record: record}, nil
}

func NewController(repository Repository) Controller {
	return controller{repository: repository}
}

func NewDefaultController() Controller {
	return NewController(NewDefaultInMemRepository())
}
