package project_test

import (
	//h "github.com/lastbackend/lastbackend/libs/http"
	//"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	//"github.com/lastbackend/lastbackend/pkg/client/context"
	//"github.com/stretchr/testify/assert"
	//"net/http"
	//"net/http/httptest"
	"testing"
)

func TestGet_Success(t *testing.T) {

	//const (
	//	name        string = "project"
	//	description string = "project describe"
	//	token       string = "mocktoken"
	//)
	//
	//var ctx = context.Mock()
	//ctx.Session.Set(token)
	//
	//server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//
	//	tk := r.Header.Get("Authorization")
	//	assert.NotEmpty(t, tk, "token should be not empty")
	//	assert.Equal(t, tk, "Bearer "+token, "they should be equal")
	//
	//	w.WriteHeader(200)
	//	_, err := w.Write([]byte(`{"id":"mock", "name":"` + name + `", "description":"` + description + `"}`))
	//	if err != nil {
	//		t.Error(err)
	//		return
	//	}
	//}))
	//defer server.Close()
	//
	//ctx.HTTP = h.New(server.URL)
	//
	//err := project.Get(name)
	//if err != nil {
	//	t.Error(err)
	//}
}
