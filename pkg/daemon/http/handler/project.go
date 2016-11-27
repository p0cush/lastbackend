package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/utils"
	"io"
	"io/ioutil"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"net/http"
)

func ProjectListH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		err     *e.Err
		session *model.Session
		ctx     = c.Get()
	)

	ctx.Log.Debug("List project handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.User.AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	projects, err := ctx.Storage.Project().GetByUser(session.Uid)
	if err != nil {
		ctx.Log.Error("Error: find projects by user", err)
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := projects.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func ProjectInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		err     *e.Err
		session *model.Session
		ctx     = c.Get()
		params  = mux.Vars(r)
		id      = params["id"]
	)

	ctx.Log.Debug("Get project handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.User.AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)
	var project *model.Project

	if !utils.IsUUID(id) {
		project, err = ctx.Storage.Project().GetByName(session.Uid, id)
	} else {
		project, err = ctx.Storage.Project().GetByID(session.Uid, id)
	}

	if err == nil && project == nil {
		e.Project.NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Err())
		err.Http(w)
		return
	}

	response, err := project.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

type projectCreate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *projectCreate) decodeAndValidate(reader io.Reader) *e.Err {

	var (
		err error
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return e.User.Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return e.Project.IncorrectJSON(err)
	}

	if s.Name == nil {
		return e.Project.BadParameter("name")
	}

	if s.Description == nil {
		s.Description = new(string)
	}

	return nil
}

func ProjectCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		err     *e.Err
		session *model.Session
		ctx     = c.Get()
	)

	ctx.Log.Debug("Create project handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.User.AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	// request body struct
	rq := new(projectCreate)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		err.Http(w)
		return
	}

	p := new(model.Project)
	p.User = session.Uid
	p.Name = *rq.Name
	p.Description = *rq.Description

	exists, er := ctx.Storage.Project().ExistByName(p.User, p.Name)
	if er != nil {
		ctx.Log.Error("Error: check exists by name", er.Error())
		e.HTTP.InternalServerError(w)
		return
	}
	if exists {
		e.Project.NameExists().Http(w)
		return
	}

	project, err := ctx.Storage.Project().Insert(p)
	if err != nil {
		ctx.Log.Error("Error: insert project to db", err)
		e.HTTP.InternalServerError(w)
		return
	}

	namespace := &v1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name:      project.ID,
			Namespace: project.ID,
			Labels: map[string]string{
				"user": session.Username,
			},
		},
	}

	_, er = ctx.K8S.Core().Namespaces().Create(namespace)
	if er != nil {
		ctx.Log.Error("Error: create namespace", er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := project.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

type projectReplace struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *projectReplace) decodeAndValidate(reader io.Reader) *e.Err {

	var (
		err error
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return e.User.Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return e.Project.IncorrectJSON(err)
	}

	if s.Name == nil {
		return e.Project.BadParameter("name")
	}

	if s.Description == nil {
		s.Description = new(string)
	}

	return nil
}

func ProjectUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		err     *e.Err
		session *model.Session
		ctx     = c.Get()
		params  = mux.Vars(r)
		id      = params["id"]
	)

	ctx.Log.Debug("Update project handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.User.AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	// request body struct
	rq := new(projectReplace)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		err.Http(w)
		return
	}

	if !utils.IsUUID(id) {
		project, err := ctx.Storage.Project().GetByName(session.Uid, id)
		if err == nil && project == nil {
			e.Project.NotFound().Http(w)
			return
		}
		if err != nil {
			ctx.Log.Error("Error: find project by id", err.Err())
			err.Http(w)
			return
		}

		id = project.ID
	}

	p := new(model.Project)
	p.ID = id
	p.User = session.Uid
	p.Name = *rq.Name
	p.Description = *rq.Description

	exists, er := ctx.Storage.Project().ExistByName(p.User, p.Name)
	if er != nil {
		ctx.Log.Error("Error: check exists by name", er.Error())
		e.HTTP.InternalServerError(w)
		return
	}
	if exists {
		e.Project.NameExists().Http(w)
		return
	}

	project, err := ctx.Storage.Project().Update(p)
	if err != nil {
		ctx.Log.Error("Error: insert project to db", err)
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := project.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func ProjectRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		ctx     = c.Get()
		session *model.Session
		params  = mux.Vars(r)
		id      = params["id"]
	)

	ctx.Log.Info("Remove project")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.User.AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	if !utils.IsUUID(id) {
		project, err := ctx.Storage.Project().GetByName(session.Uid, id)
		if err == nil && project == nil {
			e.Project.NotFound().Http(w)
			return
		}
		if err != nil {
			ctx.Log.Error("Error: find project by id", err.Err())
			err.Http(w)
			return
		}

		id = project.ID
	}

	var opts = new(api.DeleteOptions)

	er = ctx.K8S.Core().Namespaces().Delete(id, opts)
	if er != nil {
		ctx.Log.Error("Error: remove namespace", er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	err := ctx.Storage.Project().Remove(id)
	if err != nil {
		ctx.Log.Error("Error: remove project from db", err)
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
