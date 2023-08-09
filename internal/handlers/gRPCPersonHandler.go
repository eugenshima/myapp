// Package handlers provides HTTP/2 request handler functions for a web service written in Go using gRPC (Remote Procedure Call)
package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/eugenshima/myapp/internal/model"
	protos "github.com/eugenshima/myapp/proto_services"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GRPCPersonHandler struct contains service service.PersonService
type GRPCPersonHandler struct {
	srv GRPCPersonService
	protos.UnimplementedPersonHandlerServer
}

// NewGRPCPersonHandler creates a new GRPCUserHandler object
func NewGRPCPersonHandler(srv GRPCPersonService) *GRPCPersonHandler {
	return &GRPCPersonHandler{
		srv: srv,
	}
}

// GRPCPersonService represents the service implementation methods
type GRPCPersonService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]*model.Person, error)
	Delete(ctx context.Context, uuidString uuid.UUID) (uuid.UUID, error)
	Create(ctx context.Context, entity *model.Person) (uuid.UUID, error)
	Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person) (uuid.UUID, error)
}

// GetByID function receives request to Get user from database by ID
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
		Age:       result.Age,
		IsHealthy: result.IsHealthy,
	}
	return &protos.GetByIDResponse{Person: person}, nil
}

// GetAll function receives request to to Get All users from database
func (s *GRPCPersonHandler) GetAll(ctx context.Context, _ *protos.GetAllRequest) (*protos.GetAllResponse, error) {
	results, err := s.srv.GetAll(ctx)
	if err != nil {
		logrus.Errorf("GetAll: %v", err)
		return nil, err
	}
	res := []*protos.Person{}
	for _, result := range results {
		person := &protos.Person{
			Id:        result.ID.String(),
			Name:      result.Name,
			Age:       result.Age,
			IsHealthy: result.IsHealthy,
		}
		res = append(res, person)
	}
	return &protos.GetAllResponse{Person: res}, nil
}

// Delete function receives request to Delete concrete person from database
func (s *GRPCPersonHandler) Delete(ctx context.Context, req *protos.DeleteRequest) (*protos.DeleteResponse, error) {
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": req.Id}).Errorf("GetByID: %v", err)
		return nil, err
	}
	deletedID, err := s.srv.Delete(ctx, ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": &deletedID}).Errorf("Delete: %v", err)
		return nil, err
	}
	return &protos.DeleteResponse{
		Id: deletedID.String(),
	}, nil
}

// Create function receives request to Create person in database
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

// Update function receives request to Update person in database
func (s *GRPCPersonHandler) Update(ctx context.Context, req *protos.UpdateRequest) (*protos.UpdateResponse, error) {
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": req.Id}).Errorf("GetByID: %v", err)
		return nil, err
	}
	updatePerson := model.Person{
		ID:        ID,
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

// DownloadImage downloads image from given path
func (h *GRPCPersonHandler) DownloadImage(req *protos.DownloadImageRequest, stream protos.PersonHandler_DownloadImageServer) error {
	fmt.Println(req.Name)
	imgname := req.Name
	imgpath := filepath.Join("internal", "images", imgname)
	cleanPath := filepath.Clean(imgpath)
	file, err := os.Open(cleanPath)
	if err != nil {
		logrus.Errorf("failed to open file error: %v", err)
		return err
	}
	defer func() {
		errClose := file.Close()
		if errClose != nil {
			logrus.Errorf("failed to close file error: %v", errClose)
		}
	}()

	bufferSize := 4096
	buffer := make([]byte, bufferSize)
	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		logrus.Errorf("failed to read file error: %v", err)
		return err
	}
	if bytesRead == 0 {
		return nil
	}
	err = stream.Send(&protos.DownloadImageResponse{Image: buffer[:bytesRead]})
	if err != nil {
		logrus.Errorf("failed to send GRPC response: %v", err)
		return err
	}
	return nil
}

// UploadImage uploads image from given path
func (h *GRPCPersonHandler) UploadImage(stream protos.PersonHandler_UploadImageServer) error {
	dst, err := os.Create(filepath.Join("internal", "images", "test.png"))
	if err != nil {
		logrus.Errorf("failed to create file error: %v", err)
		return err
	}
	defer func() {
		errClose := dst.Close()
		if errClose != nil {
			logrus.Errorf("failed to close error: %v", errClose)
		}
	}()
	for {
		fileChunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Errorf("failed to receive data chunk error: %v", err)
			return err
		}
		_, err = dst.Write(fileChunk.Image)
		if err != nil {
			logrus.Errorf("failed to write data chunk to the dst file error: %v", err)
		}
	}
	return nil
}
