package project

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func RemoveCmd(name string) {

	var ctx = context.Get()

	err := Remove(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func Remove(name string) error {

	var (
		err error
		ctx = context.Get()
	)

	token := struct {
		Token string `json:"token"`
	}{}

	err = ctx.Storage.Get("session", &token)
	if err != nil {
		return errors.New(err.Error())
	}
	if token.Token == "" {
		return errors.New(e.StatusAccessDenied)
	}

	er := new(e.Http)
	res := struct{}{}

	_, _, err = ctx.HTTP.
		DELETE("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+token.Token).
		Request(&res, er)
	if err != nil {
		return errors.New(err.Error())
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	return nil
}
