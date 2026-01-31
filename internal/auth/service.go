package auth

import (
	"context"
	"vietio/config"
)

type Service struct {
	Config *config.Config
}

func NewService(config *config.Config) *Service {
	return &Service{
		Config: config,
	}
}

func (s *Service) GetJwtToken(ctx context.Context, payload AuthLoginRequestBody) (AuthLoginResponse, error) {
	var result AuthLoginResponse

	result.Token = "new_token"
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
