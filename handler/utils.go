package handler

import (
	"encoding/json"
	mux "github.com/dimfeld/httptreemux"
	"github.com/marove2000/hack-and-pay/errors"
	"net/http"
)

type superFunc func(r *http.Request, pathParams map[string]string) (interface{}, error)

func wrapError(fn superFunc) mux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := fn(r, params)
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
