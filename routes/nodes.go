package routes

import (
	"encoding/json"
	"net/http"

	"github.com/MixinNetwork/safe/governance/models"
	"github.com/MixinNetwork/safe/governance/session"
	"github.com/MixinNetwork/safe/governance/views"
	"github.com/dimfeld/httptreemux"
)

type nodeRequest struct {
	Extra string `json:"extra"`
}

type nodeImpl struct{}

func registerNode(router *httptreemux.TreeMux) {
	impl := &nodeImpl{}

	router.POST("/nodes", impl.create)
	router.GET("/nodes", impl.index)
}

func (impl *nodeImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body nodeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if node, err := models.CreateNodeByExtra(r.Context(), body.Extra); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderNode(w, r, node)
	}
}

func (impl *nodeImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	nodes, err := models.ReadNodes(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderNodes(w, r, nodes)
	}
}
