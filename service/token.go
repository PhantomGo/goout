package service

type TokenService struct {
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (svc *TokenService) Generate(openId string) (token string, err error) {
	token = openId
	return
}
