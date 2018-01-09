package redfi

import (
	"fmt"
	"net/http"
)

type API struct {
	plan *Plan
}

func NewAPI(p *Plan) *API {
	return &API{plan: p}
}

func (a *API) listRules(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(a.plan.ListRules())
}

func (a *API) createRule(rw http.ResponseWriter, req *http.Request) {
	panic("todo")
}

func (a *API) getRule(rw http.ResponseWriter, req *http.Request) {
	panic("todo")
}

func (a *API) deleteRule(rw http.ResponseWriter, req *http.Request) {
	panic("todo")
}
