package main

import (
	"authentication/grpc/server/auth"
	"authentication/internal/gmail"
	"authentication/models"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

func (a *Application) CreateUser(ctx context.Context, user *auth.UserDetails) (*auth.AuthResponse, error) {
	userD := models.User{
		Email:    user.Email,
		Password: user.Password,
		UserName: user.UserName,
	}

	ok := a.Store.Auth.CreateUser(ctx, userD)

	if !ok {
		return &auth.AuthResponse{Message: "Fail"}, errors.New("failed to create an user")
	}

	token := a.SetToken(ctx, user.Email)

	go gmail.SendMail(user.Email, token)

	return &auth.AuthResponse{Message: "Created User"}, nil
}

func (a *Application) SetToken(ctx context.Context, email string) string {
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	t := email + time.Now().String() + fmt.Sprintf("%x", uuid)

	h := sha256.New()
	h.Write([]byte(t))
	bs := h.Sum(nil)

	token := fmt.Sprintf("%x", bs)

	a.Store.Redis.SetEmailToken(ctx, email, token)

	return token
}

func (a *Application) VerifyUser(ctx context.Context, token *auth.Token) (*auth.VerifyResponse, error) {
	value := a.Store.Redis.GetEmailFromToken(ctx, token.Token)

	email := strings.Split(value, ":")[2]

	err := a.Store.Auth.VerifyUser(ctx, email)
	if err != nil {
		log.Println(err)
		return &auth.VerifyResponse{Message: "Fail"}, err
	}

	err = a.Store.Redis.DeleteEmailToken(ctx, token.Token)
	if err != nil {
		return &auth.VerifyResponse{Message: "Expired"}, err
	}

	return &auth.VerifyResponse{Message: "Success"}, nil
}

func (a *Application) LoginUser(ctx context.Context, user *auth.UserDetails) (*auth.LoginResponse, error) {
	userData, err := a.Store.Auth.LoginUser(ctx, user.UserName, user.Password)

	if err != nil {
		return &auth.LoginResponse{
			Id:       userData.Id,
			UserName: userData.UserName,
			Role:     userData.Role,
		}, err
	}

	return &auth.LoginResponse{Id: userData.Id, UserName: userData.UserName, Role: userData.Role}, nil
}
