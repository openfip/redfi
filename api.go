package redfi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type API struct {
	plan *Plan
}

var ErrMsg = `{"ok": false, "msg": %s}`

func NewAPI(p *Plan) *API {
	return &API{plan: p}
}

type Response struct {
	OK      bool   `json:"ok"`
	Message string `json:"msg,omitempty"`
	Rules   []Rule `json:"rules,omitempty"`
}

func (a *API) listRules(rw http.ResponseWriter, req *http.Request) {
	resp := Response{}

	rules := a.plan.ListRules()
	resp.OK = true
	resp.Rules = rules
	err := writeResponse(rw, resp, http.StatusOK)
	if err == nil {
		return
	}

	// error encountered while writing the rules
	resp.OK = false
	resp.Rules = []Rule{}
	resp.Message = err.Error()

	err = writeResponse(rw, resp, http.StatusOK)
	if err != nil {
		writeErr(rw, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}

}

func (a *API) createRule(rw http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		http.Error(rw, "Please send a request body", 400)
		return
	}

	r := Rule{}
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	err = a.plan.AddRule(r)
	if err != nil {
		writeErr(rw, err.Error(), http.StatusBadRequest)
		return
	}
	resp := Response{}
	resp.OK = true
	resp.Rules = []Rule{r}
	err = writeResponse(rw, resp, http.StatusOK)
	if err != nil {
		writeErr(rw, err.Error(), http.StatusBadRequest)
		return
	}

}

func (a *API) getRule(rw http.ResponseWriter, req *http.Request) {
	ruleName := chi.URLParam(req, "ruleName")
	rule, err := a.plan.GetRule(ruleName)
	if err != nil {
		writeErr(rw, err.Error(), http.StatusNotFound)
		return
	}

	resp := Response{}
	resp.OK = true
	resp.Rules = []Rule{rule}
	err = writeResponse(rw, resp, http.StatusOK)
	if err != nil {
		writeErr(rw, err.Error(), http.StatusBadRequest)
		return
	}

}

func (a *API) deleteRule(rw http.ResponseWriter, req *http.Request) {
	ruleName := chi.URLParam(req, "ruleName")
	err := a.plan.DeleteRule(ruleName)
	if err != nil {
		writeErr(rw, err.Error(), http.StatusNotFound)
		return
	}

	resp := Response{}
	resp.OK = true
	err = writeResponse(rw, resp, http.StatusOK)
	if err != nil {
		writeErr(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func writeErr(rw http.ResponseWriter, errMsg string, status int) {
	resp := Response{}
	resp.OK = false
	resp.Message = errMsg

	err := writeResponse(rw, resp, status)
	if err != nil {
		http.Error(rw, err.Error(), 400)
	}
}

func writeResponse(rw http.ResponseWriter, r Response, status int) error {
	out, err := json.Marshal(r)
	if err != nil {
		return err
	}

	rw.WriteHeader(status)
	_, err = rw.Write(out)
	return err
}
