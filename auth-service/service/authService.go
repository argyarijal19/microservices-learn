package service

import (
	"auth-service/lib"
	"auth-service/model"
	"auth-service/repository"
	"log"
	"os"
)

type AuthService interface {
	GenerateToken(clientData model.User) (*model.TokenResponse, error)
	CheckToken(token string) (*model.User, error)
	RevokeToken(token string) error
}

type authService struct {
	repository repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) AuthService {
	return &authService{
		repository: r,
	}
}

func (s *authService) GenerateToken(clientData model.User) (*model.TokenResponse, error) {
	data, ok := s.repository.GetDataByClientID(clientData.ClientID)
	if !ok || data == nil {
		return nil, &lib.NotFoundError{
			Message:    "client not found",
			MessageDev: "client not found",
		}
	}

	if data["client_secret"] != clientData.ClientSecret {
		return nil, &lib.UnauthorizedError{
			Message:    "client secret is invalid",
			MessageDev: "client secret is invalid",
		}
	}

	token, exp, err := lib.GenerateToken(&clientData)
	if err != nil {
		return nil, &lib.InternalServerError{
			Message:    "internal server error",
			MessageDev: err.Error(),
		}
	}

	tokenResponse := &model.TokenResponse{
		Token:     token,
		ExpiresIn: int(exp),
	}

	return tokenResponse, nil
}

func (s *authService) CheckToken(token string) (*model.User, error) {
	log.Println("CheckToken token: ", token)
	claims, err := lib.VerifyJWTToken(token, os.Getenv("SECRET_KEY_JWT"))
	if err != nil {
		if err.Error() == "invalid access" {
			return nil, &lib.UnauthorizedError{
				Message:    "invalid access",
				MessageDev: "invalid access",
			}
		} else if err.Error() == "token has been revoked" {
			return nil, &lib.UnauthorizedError{
				Message:    "token has been revoked",
				MessageDev: "token has been revoked",
			}
		}
		return nil, &lib.InternalServerError{
			Message:    "internal server error",
			MessageDev: err.Error(),
		}
	}
	jwtClaims, ok := claims.(*lib.JWTCustomClaims)
	if !ok {
		return nil, &lib.InternalServerError{
			Message:    "internal server error",
			MessageDev: "token conversion error",
		}
	}

	userData := &model.User{
		ClientID:     jwtClaims.ClientID,
		ClientSecret: jwtClaims.ClientSecret,
	}

	return userData, nil

}

func (s *authService) RevokeToken(token string) error {
	if err := lib.RevokeToken(token); err != nil {
		return &lib.InternalServerError{
			Message:    "internal server error",
			MessageDev: err.Error(),
		}
	}

	return nil
}
