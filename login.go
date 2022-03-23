// Copyright 2020 diamondburned

// Permission to use, copy, modify, and/or distribute this software for any purpose
// with or without fee is hereby granted, provided that the above copyright notice
// and this permission notice appear in all copies.

// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
// REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
// FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
// INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS
// OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
// TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF
// THIS SOFTWARE.

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
