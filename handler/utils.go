package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/marove2000/hack-and-pay/errors"

	mux "github.com/dimfeld/httptreemux"
	"github.com/marove2000/hack-and-pay/ctxutil"
	"github.com/pborman/uuid"
)

type superFunc func(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error)

func wrap(fn superFunc) mux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		ctx := context.Background()

		var corrId string
		if corrIdSlice := r.Header["correlation-id"]; len(corrIdSlice) != 0 {
			uid := uuid.Parse(corrIdSlice[0])
			if uid == nil {
				corrId = uid.String()
			}
		}
		if corrId == "" {
			corrId = uuid.NewRandom().String()
		}

		resp, err := fn(ctxutil.InjectCorrelationId(ctx, corrId), r, params)
		if err != nil {
			w.WriteHeader(err.(*errors.Error).Status)
			json.NewEncoder(w).Encode(err)
			return
		}

		if resp == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}
