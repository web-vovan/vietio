package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"vietio/config"
	appUser "vietio/internal/user"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	Config    *config.Config
	Validator *Validator
	UserRepo  UserRepo
}

type UserRepo interface {
	GetUserByTelegramId(ctx context.Context, telegramId int64) (appUser.UserModel, error)
	UpdateUsername(context.Context, appUser.UserModel) error
	CreateUser(context.Context, appUser.UserModel) (id int64, err error)
}

func NewService(config *config.Config, validator *Validator, userRepo UserRepo) *Service {
	return &Service{
		Config:    config,
		Validator: validator,
		UserRepo:  userRepo,
	}
}

func (s *Service) GetJwtToken(ctx context.Context, payload AuthLoginRequestBody) (AuthLoginResponse, error) {
	var result AuthLoginResponse

	telegramUser, err := s.Validator.ValidateWebAppData(payload.InitData, s.Config.BotToken)
	if err != nil {
		return result, err
	}

	user, err := s.UserRepo.GetUserByTelegramId(ctx, telegramUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// создаем нового пользователя
			user := appUser.UserModel{
				TelegramId: telegramUser.ID,
				Username:   telegramUser.Username,
			}
			id, err := s.UserRepo.CreateUser(ctx, user)
			if err != nil {
				return result, err
			}
			user.Id = id
		} else {
			return result, err
		}
	}

	// обновляем username
	if user.Username != telegramUser.Username {
		user.Username = telegramUser.Username
		err = s.UserRepo.UpdateUsername(ctx, user)
		if err != nil {
			return result, err
		}
	}

	token, err := s.generateJwtToken(user)
	if err != nil {
		return result, err
	}

	result.Token = token
	return result, nil
}

func (s *Service) GenerateTestInitData(username string) (TestInitDataResponse, error) {
	var result TestInitDataResponse
	initData, err := generateTestInitData(s.Config.BotToken, username)
	if err != nil {
		return result, err
	}

	result.InitData = initData

	return result, nil
}

func (s *Service) generateJwtToken(user appUser.UserModel) (string, error) {
	claims := AccessTokenClaims{
		UserId: user.Id,
		TelegramId:  user.TelegramId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString([]byte(s.Config.JwtSecret))
}

func (s *Service) ParseAndValidateJWT(tokenString string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(s.Config.JwtSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
