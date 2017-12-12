package redfi

import "testing"

func TestSelectRule(t *testing.T) {
	p := &Plan{
		Rules: []Rule{},
	}
	p.Rules = append(p.Rules, Rule{
		Delay:      1e3,
		ClientAddr: "192.0.0.1:8001",
		Percentage: 20,
	})
}
