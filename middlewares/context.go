package middlewares

import (
	"net/http"

	"github.com/MixinNetwork/safe/governance/session"
	"github.com/MixinNetwork/safe/governance/store"
	"github.com/unrolled/render"
)

func Context(handler http.Handler, db *store.Database, render *render.Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := session.WithRequest(r.Context(), r)
		ctx = session.WithDatabase(ctx, db)
		ctx = session.WithRender(ctx, render)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
