package vault

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/crypto/bcrypt"
	//context2 "golang.org/x/net/context"
)

// Service provide password hashing capabilities.
type Service interface {
	Hash(ctx context.Context, password string) (string, error)
	Validate(ctx context.Context, password, hash string) (bool, error)
}

type vaultService struct{}

// NewService constructs new service
func NewService() Service {
	return vaultService{}
}

func (vaultService) Hash(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func (vaultService) Validate(ctx context.Context, password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

type hashRequest struct {
	Password string `json:"password"`
}
type hashResponse struct {
	Hash string `json:"hash"`
	Err  string `json:"err,omitempty"`
}

func decodeHashRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req hashRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

type validateRequest struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}
type validateResponse struct {
	Valid bool   `json:"valid"`
	Err   string `json:"err,omitempty"`
}

func decodeValidateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req validateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// MakeHashEndpoint create go-kit endpoint for connection between the communication portocol
// and Vault service
func MakeHashEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hashRequest)
		h, err := srv.Hash(ctx, req.Password)
		if err != nil {
			return hashResponse{h, err.Error()}, nil
		}
		return hashResponse{h, ""}, nil
	}
}

// MakeValidateEndpoint create go-kit endpoint for connection between the communication portocol
// and Vault service
func MakeValidateEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validateRequest)
		valid, err := srv.Validate(ctx, req.Password, req.Hash)
		if err != nil {
			return validateResponse{false, err.Error()}, nil
		}
		return validateResponse{valid, ""}, nil
	}
}

// Endpoints struct wrap Hash and Validate endpoints while implementing the Service interface
type Endpoints struct {
	HashEndpoint     endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

// Hash hashes the given password. However instead of directly using the service, it calls HashEndpoint.
// Also implements Hash function from Service inteface.
func (e Endpoints) Hash(ctx context.Context, password string) (string, error) {
	req := hashRequest{Password: password}
	resp, err := e.HashEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	hashResponse := resp.(hashResponse)
	if hashResponse.Err != "" {
		return "", errors.New(hashResponse.Err)
	}
	return hashResponse.Hash, nil
}

// Validate validates the password with the hash. However instead of directly using the service,
// it calls ValidateEndpoint. Also implements Validate function from Service inteface.
func (e Endpoints) Validate(ctx context.Context, password, hash string) (bool, error) {
	req := validateRequest{
		Password: password,
		Hash:     hash,
	}
	resp, err := e.ValidateEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	validateResponse := resp.(validateResponse)
	if validateResponse.Err != "" {
		return false, errors.New(validateResponse.Err)
	}
	return validateResponse.Valid, nil
}
