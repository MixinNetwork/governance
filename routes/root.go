package routes

import (
	"net/http"
	"runtime"

	"github.com/MixinNetwork/safe/governance/config"
	"github.com/MixinNetwork/safe/governance/models"
	"github.com/MixinNetwork/safe/governance/session"
	"github.com/MixinNetwork/safe/governance/views"
	"github.com/dimfeld/httptreemux"
)

func RegisterRoutes(router *httptreemux.TreeMux) {
	RegisterHanders(router)
	router.GET("/_hc", health)
	router.GET("/template", template)

	registerNode(router)
}

func health(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderDataResponse(w, r, map[string]string{
		"build": config.BuildVersion + "-" + runtime.Version(),
	})
}

func template(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	nodes, err := models.ReadNodes(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	var list []string
	for _, n := range nodes {
		list = append(list, n.AppID.String)
	}
	views.RenderTemplate(w, r, list)
}

func RegisterHanders(router *httptreemux.TreeMux) {
	router.MethodNotAllowedHandler = func(w http.ResponseWriter, r *http.Request, _ map[string]httptreemux.HandlerFunc) {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	}
	router.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	}
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, rcv any) {
		err, _ := rcv.(error)
		views.RenderErrorResponse(w, r, session.ServerError(r.Context(), err))
	}
}
