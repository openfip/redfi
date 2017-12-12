package redfi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrNotFound is returned iff SelectRule can't find a Rule that applies
	ErrNotFound = errors.New("no matching rule found")
)

// Plan defines a set of rules to be applied by the proxy
type Plan struct {
	Rules []*Rule `json:"rules,omitempty"`
}

// Rule is what get's applied on every client message iff it matches it
type Rule struct {
	Delay      int  `json:"delay,omitempty"`
	Drop       bool `json:"drop,omitempty"`
	Percentage int  `json:"percentage,omitempty"`
	// SelectRule does prefix matching on this value
	ClientAddr string `json:"client_addr,omitempty"`
	Command    string `json:"command,omitempty"`
	// filled by marshalCommand
	marshaledCmd []byte
}

// Parse the plan.json file
func Parse(planPath string) (*Plan, error) {
	fullPath, err := filepath.Abs(planPath)
	if err != nil {
		return nil, err
	}

	fd, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	plan := &Plan{}
	err = json.Unmarshal(buf, plan)
	if err != nil {
		return nil, err
	}

	for _, rule := range plan.Rules {
		if rule == nil {
			continue
		}
		if len(rule.Command) > 0 {
			rule.marshaledCmd = marshalCommand(rule.Command)
		}
	}

	return plan, nil
}

func marshalCommand(cmd string) []byte {
	return []byte(fmt.Sprintf("\r\n%s\r\n", strings.ToUpper(cmd)))
}

// SelectRule finds the first rule that applies to the given variables
func (p *Plan) SelectRule(clientAddr string, buf []byte) *Rule {
	var chosenRule *Rule
	for _, rule := range p.Rules {
		if len(rule.ClientAddr) > 0 && strings.HasPrefix(clientAddr, rule.ClientAddr) {
			continue
		}

		if len(rule.Command) > 0 && !bytes.Contains(buf, rule.marshaledCmd) {
			continue
		}

		chosenRule = rule
		break

	}
	if chosenRule == nil {
		return nil
	}

	if chosenRule.Percentage > 0 && rand.Intn(100) > chosenRule.Percentage {
		return nil
	}
	return chosenRule
}
