package handlers

import (
	"context"

	protos "github.com/eugenshima/myapp/proto_services"
	"github.com/sirupsen/logrus"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
)

type GRPCUserHandler struct {
	srv GRPCUserService
	protos.UnimplementedUserHandlerServer
}

func NewGRPCUserHandler(srv GRPCUserService) *GRPCUserHandler {
	return &GRPCUserHandler{
		srv: srv,
	}
}

type GRPCUserService interface {
	GenerateTokens(ctx context.Context, login, password string) (string, string, error)
	Signup(ctx context.Context, entity *model.User) error
	RefreshTokenPair(ctx context.Context, accessToken string, refreshToken string, id uuid.UUID) (string, string, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func (s *GRPCUserHandler) Login(ctx context.Context, req *protos.LoginRequest) (*protos.LoginResponse, error) {
	accessToken, refreshToken, err := s.srv.GenerateTokens(ctx, req.Login, req.Password)
	if err != nil {
		logrus.Errorf("GenerateTokens %v", err)
		return nil, err
	}
	accessRefresh := protos.AccessRefresh{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return &protos.LoginResponse{AccessRefresh: &accessRefresh}, nil
}

func (s *GRPCUserHandler) SignUp(ctx context.Context, req *protos.SignUpRequest) (*protos.SignUpResponse, error) {
	newUser := model.User{
		ID:       uuid.New(),
		Login:    req.Login,
		Password: []byte(req.Password),
		Role:     req.Role,
	}
	err := s.srv.Signup(ctx, &newUser)
	if err != nil {
		logrus.WithFields(logrus.Fields{"user": newUser}).Errorf("Signup: %v", err)
		return nil, err
	}
	return &protos.SignUpResponse{}, nil
}

func (s *GRPCUserHandler) GetAll(ctx context.Context, req *protos.UserGetAllRequest) (*protos.UserGetAllResponse, error) {
	results, err := s.srv.GetAll(ctx)
	if err != nil {
		logrus.Errorf("GetAll: %v", err)
		return nil, err
	}
	var res []*protos.User
	for _, result := range results {
		user := &protos.User{
			Id:           result.ID.String(),
			Login:        result.Login,
			Password:     string(result.Password),
			Role:         result.Role,
			RefreshToken: string(result.RefreshToken),
		}
		res = append(res, user)
	}
	return &protos.UserGetAllResponse{User: res}, nil
}

func (s *GRPCUserHandler) RefreshTokenPair(ctx context.Context, req *protos.RefreshTokenPairRequest) (*protos.RefreshTokenPairResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": req.Id}).Errorf("RefreshTokenPair: %v", err)
		return nil, err
	}

	accessToken, refreshToken, err := s.srv.RefreshTokenPair(ctx, req.AccessRefresh.AccessToken, req.AccessRefresh.RefreshToken, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reqBody.AccessToken": req.AccessRefresh.AccessToken, "reqBody.RefreshToken": req.AccessRefresh.RefreshToken, "id": id}).Errorf("RefreshTokenPair: %v", err)
		return nil, err
	}
	accessRefresh := protos.AccessRefresh{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return &protos.RefreshTokenPairResponse{AccessRefresh: &accessRefresh}, nil
}

func (s *GRPCUserHandler) Delete(ctx context.Context, req *protos.UserDeleteRequest) (*protos.UserDeleteResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		logrus.Errorf("Parse: %v", err)
		return nil, err
	}
	err = s.srv.Delete(ctx, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": id}).Errorf("Delete: %v", err)
		return nil, err
	}
	return &protos.UserDeleteResponse{Id: id.String()}, nil
}
