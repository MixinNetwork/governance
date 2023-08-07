package views

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"github.com/MixinNetwork/safe/governance/session"
)

//go:embed example.toml
var exampletoml string

type ResponseView struct {
	Data  any    `json:"data,omitempty"`
	Error error  `json:"error,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}

func RenderDataResponse(w http.ResponseWriter, r *http.Request, view any) {
	session.Render(r.Context()).JSON(w, http.StatusOK, ResponseView{Data: view})
}

func RenderErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	sessionError, ok := err.(*session.Error)
	if !ok {
		sessionError = session.ServerError(r.Context(), err)
	}
	if sessionError.Code == 10001 {
		sessionError.Code = 500
	}
	session.Render(r.Context()).JSON(w, sessionError.Status, ResponseView{Error: sessionError})
}

func RenderBlankResponse(w http.ResponseWriter, r *http.Request) {
	session.Render(r.Context()).JSON(w, http.StatusOK, ResponseView{})
}

func RenderOriginalResponse(w http.ResponseWriter, r *http.Request, view any) {
	session.Render(r.Context()).JSON(w, http.StatusOK, view)
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, list []string) {
	str := strings.Join(list, `",
  "`)
	l := len(list)
	session.Render(r.Context()).Text(w, http.StatusOK, fmt.Sprintf(exampletoml, l*2/3, str, l*2/3+1, str, l*2/3+1))
}
