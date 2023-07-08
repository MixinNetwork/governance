package views

import (
	"net/http"
	"time"

	"github.com/MixinNetwork/safe/governance/models"
)

type NodeView struct {
	Custodian string    `json:"custodian"`
	Payee     string    `json:"payee"`
	KernelID  string    `json:"kernel_id"`
	AppID     string    `json:"app_id"`
	MixinHash string    `json:"mixin_hash"`
	Keystore  string    `json:"keystore"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func buildNodeView(n *models.Node) *NodeView {
	return &NodeView{
		Custodian: n.Custodian,
		Payee:     n.Payee,
		KernelID:  n.KernelID,
		AppID:     n.AppID.String,
		MixinHash: n.MixinHash.String,
		Keystore:  n.Keystore,
		PublicKey: n.PublicKey,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

func RenderNode(w http.ResponseWriter, r *http.Request, node *models.Node) {
	view := buildNodeView(node)
	RenderDataResponse(w, r, view)
}

func RenderNodes(w http.ResponseWriter, r *http.Request, nodes []*models.Node) {
	views := make([]*NodeView, len(nodes))
	for i, n := range nodes {
		views[i] = buildNodeView(n)
	}
	RenderDataResponse(w, r, views)
}
