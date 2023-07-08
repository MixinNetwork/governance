package middlewares

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/MixinNetwork/mixin/config"
	"github.com/MixinNetwork/safe/governance/session"
	"github.com/MixinNetwork/safe/governance/views"
)

func Stats(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("INFO -- : Started %s '%s'\n", r.Method, r.URL)
		defer func() {
			log.Printf("INFO -- : Completed %s in %fms\n", r.Method, time.Now().Sub(start).Seconds())
		}()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
			return
		}
		if len(body) > 0 {
			log.Printf("INFO -- : Paremeters %s\n", string(body))
		}
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		r = r.WithContext(session.WithRequestBody(r.Context(), string(body)))
		w.Header().Set("X-Build-Info", config.BuildVersion+"-"+runtime.Version())
		handler.ServeHTTP(w, r)
	})
}
