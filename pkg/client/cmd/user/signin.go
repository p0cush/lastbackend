package user

import (
	"errors"
	"fmt"
	"github.com/howeyc/gopass"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func SignInCmd() {

	var (
		login    string
		password string
		ctx      = context.Get()
	)

	fmt.Print("Login: ")
	fmt.Scan(&login)

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		ctx.Log.Error(err)
		return
	}
	password = string(pass)

	err = SignIn(login, password)
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

type userLoginS struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func SignIn(login, password string) error {

	var (
		err error
		ctx = context.Get()
	)

	//fmt.Println(ctx)

	er := new(e.Http)
	res := struct {
		Token string `json:"token"`
	}{}

	_, _, err = ctx.HTTP.
		POST("/session").
		AddHeader("Content-Type", "application/json").
		BodyJSON(userLoginS{login, password}).
		Request(&res, er)
	if err != nil {
		return err
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	} else {
		fmt.Println("Login successful")
	}

	err = ctx.Storage.Set("session", res)
	if err != nil {
		return err
	}

	return nil
}
