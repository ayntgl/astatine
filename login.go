package astatine

import "encoding/json"

type LoginResponse struct {
	Mfa    bool   `json:"mfa"`
	Sms    bool   `json:"sms"`
	Ticket string `json:"ticket"`
	Token  string `json:"token"`
}

func (s *Session) Login(email, password string) (*LoginResponse, error) {
	data := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{email, password}
	resp, err := s.RequestWithBucketID("POST", EndpointLogin, data, EndpointLogin)

	var lr *LoginResponse
	err = json.Unmarshal(resp, &lr)
	if err != nil {
		return nil, err
	}

	return lr, nil
}

func (s *Session) Totp(code, ticket string) (*LoginResponse, error) {
	data := struct {
		Code   string `json:"code"`
		Ticket string `json:"ticket"`
	}{code, ticket}
	resp, err := s.RequestWithBucketID("POST", EndpointTotp, data, EndpointTotp)
	if err != nil {
		return nil, err
	}

	var lr *LoginResponse
	err = json.Unmarshal(resp, &lr)
	if err != nil {
		return nil, err
	}

	return lr, nil
}
