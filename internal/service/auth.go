package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/auths"
	"github.com/mikalai2006/kingwood-api/pkg/hasher"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	hasher       hasher.PasswordHasher
	tokenManager auths.TokenManager

	repository   repository.Authorization
	otpGenerator utils.Generator

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	verificationCodeLength int
	UserService            *UserService
	Hub                    *Hub
}

func NewAuthService(
	repo repository.Authorization,
	hasherPkg hasher.PasswordHasher,
	tokenManager auths.TokenManager,
	refreshTokenTTL time.Duration,
	accessTokenTTL time.Duration,
	otp utils.Generator,
	verificationCodeLength int,
	userService *UserService,
	Hub *Hub,
) *AuthService {
	return &AuthService{
		hasher:       hasherPkg,
		tokenManager: tokenManager,

		repository:   repo,
		otpGenerator: otp,

		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,

		verificationCodeLength: verificationCodeLength,
		UserService:            userService,
		Hub:                    Hub,
	}
}

func (s *AuthService) CreateAuth(auth *domain.AuthInput) (string, error) {
	passwordHash, err := s.hasher.Hash(auth.Password)
	if err != nil {
		return "", err
	}

	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)

	authData := &domain.AuthInputMongo{
		// VkID:      auth.VkID,
		// GoogleID:  auth.GoogleID,
		// GithubID:  auth.GithubID,
		// AppleID:   auth.AppleID,
		// Roles:     []string{"user"},
		Login:     auth.Login,
		Password:  passwordHash,
		Email:     auth.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Strategy:  auth.Strategy, //"local"
		Verification: domain.Verification{
			Code: verificationCode,
		},
		// MaxDistance: auth.MaxDistance,
		// RoleId: auth.RoleId,
		// PostId: auth.PostId,
	}

	id, err := s.repository.CreateAuth(authData)
	if err != nil {
		return "", err
	}

	// if created auth, send email with verification code

	return id, nil
}

func (s *AuthService) ExistAuth(auth *domain.AuthInput) (domain.Auth, error) {
	return s.repository.CheckExistAuth(auth)
}

func (s *AuthService) GetAuth(id string) (domain.Auth, error) {
	return s.repository.GetAuth(id)
}

func (s *AuthService) SignIn(auth *domain.AuthInput) (domain.ResponseTokens, error) {
	var result domain.ResponseTokens
	passwordHash, err := s.hasher.Hash(auth.Password)
	if err != nil {
		return result, err
	}
	auth.Password = passwordHash

	user, err := s.repository.GetByCredentials(auth)
	// fmt.Println("sign in ", user.Role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, err
		}
		return result, err
	}

	return s.CreateSession(&user)
}

func (s *AuthService) ResetPassword(authID string, userID string, input *domain.ResetPassword) (string, error) {
	result := ""

	newPassword := input.Password
	if newPassword == "" {
		newPassword = s.tokenManager.GetRandomPassword()
	}

	passwordHash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return result, err
	}

	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)

	authData := &domain.AuthInput{
		Password:  passwordHash,
		UpdatedAt: time.Now(),
		Verification: domain.Verification{
			Code: verificationCode,
		},
		// MaxDistance: auth.MaxDistance,
		// RoleId: auth.RoleId,
		// PostId: auth.PostId,
	}

	_, err = s.UpdateAuth(authID, authData)
	if err != nil {
		return result, err
	}

	result = newPassword

	return result, nil
}

func (s *AuthService) CreateSession(auth *domain.Auth) (domain.ResponseTokens, error) {
	var (
		res domain.ResponseTokens
		err error
	)

	claims := domain.DataForClaims{
		// Roles:  auth.Role.Value,
		UserID: auth.ID.Hex(),
		// Md:     auth.MaxDistance,
		UID: auth.User.ID.Hex(),
	}

	res.AccessToken, err = s.tokenManager.NewJWT(claims, s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	// expiresIn := time.Now().Add(s.refreshTokenTTL)

	timeExpires := time.Now().Local().Add(time.Duration(s.accessTokenTTL))
	// time.Hour*time.Duration(timeDuration.Hours()) +
	// 	time.Minute*time.Duration(timeDuration.Minutes()) +
	// fmt.Println("expiresIn: ", timeExpires, timeExpires.UnixMilli(), s.refreshTokenTTL.Minutes(), time.Now().Add(s.refreshTokenTTL).UnixMilli())

	res.ExpiresIn = timeExpires.UnixMilli()

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Local().Add(time.Duration(s.refreshTokenTTL)), // time.Now().Add(s.refreshTokenTTL),
	}

	res.ExpiresInR = session.ExpiresAt.UnixMilli()

	err = s.repository.SetSession(auth.ID, session)

	return res, err
}

func (s *AuthService) VerificationCode(userID, hash string) error {
	err := s.repository.VerificationCode(userID, hash)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (domain.ResponseTokens, error) {
	var result domain.ResponseTokens

	user, err := s.repository.RefreshToken(refreshToken)
	fmt.Println("RefreshTokens service", err, refreshToken)
	if err != nil {
		return result, err
	}

	return s.CreateSession(&user)
}

func (s *AuthService) RemoveRefreshTokens(refreshToken string) (string, error) {
	var result = ""

	_, err := s.repository.RemoveRefreshToken(refreshToken)
	if err != nil {
		return result, err
	}

	return result, err
}

func (s *AuthService) UpdateAuth(id string, data *domain.AuthInput) (domain.Auth, error) {
	result, err := s.repository.UpdateAuth(id, data)

	// if data.AppInfo.VersionApp != "" || data.AppInfo.VersionBuild != "" {
	// 	fmt.Println("Update user AppInfo from auth update", data.AppInfo)
	// 	_, err = s.UserService.UpdateUser(result.User.ID.Hex(), &domain.UserInput{
	// 		AppInfo: data.AppInfo,
	// 	})
	// }
	// s.Hub.HandleMessage(domain.Message{Type: "message", Method: "PATCH", Sender: id, Recipient: "user2", Content: result, ID: "room1", Service: "user"})

	return result, err
}

func (s *AuthService) DeleteAuth(id string) (domain.Auth, error) {
	return s.repository.DeleteAuth(id)
}
