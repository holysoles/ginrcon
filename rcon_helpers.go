package main

import (
	"net"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorcon/rcon"
)

func openRcon(i openConnInfo) (*rcon.Conn, error) {
	// basic validation so we can return more verbose error messages
	err := i.validateConnectionInfo()
	if err != nil {
		return nil, ErrInvalidConnectionDetails
	}

	//actually open
	conn, err := rcon.Dial(i.Server, i.Password)
	if err != nil {
		switch err {
		case rcon.ErrAuthNotRCON:
			fallthrough
		case rcon.ErrInvalidAuthResponse:
			return nil, ErrInvalidResponseFromRcon
		case rcon.ErrAuthFailed:
			fallthrough
		default:
			return nil, ErrInvalidConnectionDetails
		}
	}
	return conn, nil
}

func (i *openConnInfo) validateConnectionInfo() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Var(i.Server, "required,hostname_port")
	if err != nil {
		return err
	}
	dialer := &net.Dialer{Timeout: time.Second * 1}
	_, err = dialer.Dial("tcp", i.Server)
	if err != nil {
		return err
	}
	return nil
}
