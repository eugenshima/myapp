// Package handlers provides HTTP/2 request handler functions for a web service written in Go using gRPC (Remote Procedure Call)
package handlers

import (
	"bufio"
	shadowByte "bytes"
	"context"
	"image"
	"image/jpeg"
	"os"
	"sync"

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

// UploadImage uploads image to server
func (s *GRPCPersonHandler) UploadImage(src protos.PersonHandler_UploadImageServer) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		logrus.Info("Stream started")
		for {
			rr, err := src.Recv()
			if err != nil {
				break
			}
			logrus.Infof("Received message: %v", rr)
			img, err := ImageToBytes(rr.Name)
			if err != nil {
				logrus.WithFields(logrus.Fields{"image": img}).Errorf("saveFrames: %v", err)
				src.Context().Done()
			}
			finalImg, err := bytesToImg(img)
			if err != nil {
				logrus.WithFields(logrus.Fields{"image": img}).Errorf("saveFrames: %v", err)
				src.Context().Done()
			}
			err = src.Send(&protos.UploadImageResponse{
				Image: finalImg,
			})
			if err != nil {
				logrus.WithFields(logrus.Fields{"image": img}).Errorf("Send: %v", err)
				src.Context().Done()
			}
		}

		wg.Done()
		logrus.Info("Stream finished")
	}()
	wg.Wait()

	return nil
}

// DownloadImage downloads the image from the server
func (s *GRPCPersonHandler) DownloadImage(src protos.PersonHandler_DownloadImageServer) error {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		logrus.Info("Stream started")
		for {
			rr, err := src.Recv()
			if err != nil {
				break
			}
			logrus.Infof("Received message: %v", rr.Name)
			img, err := ImageToBytes(rr.Name)
			if err != nil {
				logrus.WithFields(logrus.Fields{"image": img}).Errorf("saveFrames: %v", err)
			}
			wg.Done()
			logrus.Info("Stream finished")
		}
	}()
	wg.Wait()

	return nil
}

// ImageToBytes converts image from jpeg to []byte
//
//nolint:gosec // needed
func ImageToBytes(receivedPath string) ([]byte, error) {
	fileToBeUploaded := receivedPath
	file, err := os.Open(fileToBeUploaded)
	if err != nil {
		logrus.WithFields(logrus.Fields{"file": file}).Errorf("Open: %v", err)
		return nil, err
	}

	defer func() {
		if err = file.Close(); err != nil {
			logrus.Errorf("Close: %v", err)
		}
	}()

	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	bytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		logrus.WithFields(logrus.Fields{"bytes": bytes}).Errorf("Read: %v", err)
		return nil, err
	}
	return bytes, nil
}

// bytesToImg converts image from []byte to image
func bytesToImg(imgByte []byte) ([]byte, error) {
	img, _, err := image.Decode(shadowByte.NewReader(imgByte))
	if err != nil {
		logrus.Errorf("Decode: %v", err)
		return nil, err
	}

	out, _ := os.Create("./internal/images/img.jpeg")
	defer func() {
		if err = out.Close(); err != nil {
			logrus.Errorf("Close: %v", err)
		}
	}()

	var opts jpeg.Options
	opts.Quality = 1

	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		logrus.Errorf("Encode: %v", err)
		return nil, err
	}
	return []byte(out.Name()), nil
}
