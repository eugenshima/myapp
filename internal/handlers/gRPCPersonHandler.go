package handlers

import (
	"context"

	"github.com/eugenshima/myapp/internal/model"
	protos "github.com/eugenshima/myapp/proto_services"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PersonHandler struct contains service service.PersonService
type GRPCPersonHandler struct {
	srv GRPCPersonService
	protos.UnimplementedPersonHandlerServer
}

// NewPersonHandler is a constructor
func NewGRPCPersonHandler(srv GRPCPersonService) *GRPCPersonHandler {
	return &GRPCPersonHandler{
		srv: srv,
	}
}

// PersonService interface, which contains Service methods
type GRPCPersonService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]*model.Person, error)
	Delete(ctx context.Context, uuidString uuid.UUID) (uuid.UUID, error)
	Create(ctx context.Context, entity *model.Person) (uuid.UUID, error)
	Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person) (uuid.UUID, error)
}

type GRPCServer struct {
}

func (s *GRPCPersonHandler) GetByID(ctx context.Context, req *protos.GetByIDRequest) (*protos.GetByIDResponse, error) {
	person := &protos.Person{
		Id:        req.Id,
		Name:      "eugen",
		Age:       21,
		IsHealthy: false,
	}
	return &protos.GetByIDResponse{Person: person}, nil
}

func (s *GRPCPersonHandler) GetAll(ctx context.Context, req *protos.GetAllRequest) (*protos.GetAllResponse, error) {
	results, err := s.srv.GetAll(ctx)
	if err != nil {
		logrus.Errorf("GetAll: %v", err)
		return nil, err
	}
	var res []*protos.Person
	for _, result := range results {
		person := &protos.Person{
			Id:        result.ID.String(),
			Name:      result.Name,
			Age:       int64(result.Age),
			IsHealthy: result.IsHealthy,
		}
		res = append(res, person)
	}
	return &protos.GetAllResponse{Person: res}, nil
}

func (s *GRPCPersonHandler) Delete(ctx context.Context, in *protos.DeleteRequest) (*protos.DeleteResponse, error) {
	return &protos.DeleteResponse{
		Id: "delete working",
	}, nil
}

func (s *GRPCPersonHandler) Create(ctx context.Context, req *protos.CreateRequest) (*protos.CreateResponse, error) {
	var person *model.Person
	var err error
	person.ID, err = uuid.Parse(req.Person.Id) // TODO: fix gub w/ nil pointer dereference
	if err != nil {
		logrus.WithFields(logrus.Fields{"person": &person}).Errorf("Create: %v", err)
		return nil, err
	}
	person.Name = req.Person.Name
	person.Age = int(req.Person.Age)
	person.IsHealthy = req.Person.IsHealthy

	id, err := s.srv.Create(ctx, person)
	if err != nil {
		logrus.WithFields(logrus.Fields{"person": &person}).Errorf("Create: %v", err)
		return nil, err
	}
	return &protos.CreateResponse{Id: id.String()}, nil
}

func (s *GRPCPersonHandler) Update(context.Context, *protos.UpdateRequest) (*protos.UpdateResponse, error) {
	return &protos.UpdateResponse{
		Id: "update working",
	}, nil
}
