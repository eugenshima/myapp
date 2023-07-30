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

func (s *GRPCPersonHandler) GetByID(ctx context.Context, req *protos.GetByIDRequest) (*protos.GetByIDResponse, error) {
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": req.Id}).Errorf("GetByID: %v", err)
		return nil, err
	}
	result, err := s.srv.GetByID(ctx, ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": ID}).Errorf("GetByID: %v", err)
		return nil, err
	}
	person := &protos.Person{
		Id:        ID.String(),
		Name:      result.Name,
		Age:       int64(result.Age),
		IsHealthy: result.IsHealthy,
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

func (s *GRPCPersonHandler) Delete(ctx context.Context, req *protos.DeleteRequest) (*protos.DeleteResponse, error) {
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": req.Id}).Errorf("GetByID: %v", err)
		return nil, err
	}
	deletedId, err := s.srv.Delete(ctx, ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": &deletedId}).Errorf("Delete: %v", err)
		return nil, err
	}
	return &protos.DeleteResponse{
		Id: deletedId.String(),
	}, nil
}

func (s *GRPCPersonHandler) Create(ctx context.Context, req *protos.CreateRequest) (*protos.CreateResponse, error) {
	newPerson := model.Person{
		ID:        uuid.New(),
		Name:      req.Person.Name,
		Age:       req.Person.Age,
		IsHealthy: req.Person.IsHealthy,
	}

	id, err := s.srv.Create(ctx, &newPerson)
	if err != nil {
		logrus.WithFields(logrus.Fields{"person": newPerson}).Errorf("Create: %v", err)
		return nil, err
	}

	return &protos.CreateResponse{Id: id.String()}, nil
}

func (s *GRPCPersonHandler) Update(ctx context.Context, req *protos.UpdateRequest) (*protos.UpdateResponse, error) {
	updatePerson := model.Person{
		ID:        uuid.MustParse(req.Person.Id),
		Name:      req.Person.Name,
		Age:       req.Person.Age,
		IsHealthy: req.Person.IsHealthy,
	}
	id, err := s.srv.Update(ctx, updatePerson.ID, &updatePerson)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": id, "person": updatePerson}).Errorf("Update: %v", err)
		return nil, err
	}
	return &protos.UpdateResponse{Id: id.String()}, nil
}
