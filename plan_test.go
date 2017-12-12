package redfi

import (
	"fmt"
	"testing"
)

func TestSelectRule(t *testing.T) {
	p := &Plan{
		Rules: []*Rule{},
	}

	// // test ip matching
	p.Rules = append(p.Rules, &Rule{
		Delay:      1e3,
		ClientAddr: "192.0.0.1:8001",
	})

	rule := p.SelectRule("192.0.0.1", []byte(""))
	if rule == nil {
		t.Fatal("rule must not be nil")
	}

	// test command matching
	p.Rules = []*Rule{}
	p.Rules = append(p.Rules, &Rule{
		Delay:   1e3,
		Command: "GET",
	})
	p.MarshalCommands()

	rule = p.SelectRule("192.0.0.1", []byte("\r\nGET\r\nfff"))
	if rule == nil {
		t.Fatal("rule must not be nil")
	}

	rule = p.SelectRule("172.0.0.1", []byte("\r\nKEYS\r\nfff"))
	if rule != nil {
		fmt.Println(rule)
		t.Fatal("rule must BE nil")
	}

}
