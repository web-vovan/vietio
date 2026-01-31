package auth

type AuthLoginRequestBody struct {
    InitData string
}

type AuthLoginResponse struct {
    Token string `json:"token"`
}

type TestInitDataResponse struct {
    InitData string `json:"init_data"`
}