package redfi

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tidwall/redcon"
)

var configAddr = ":6380"

const (
	RULEADD      = "ruleadd"
	RULEDEL      = "ruledel"
	RULELIST     = "rulelist"
	RULEGET      = "ruleget"
	RULECOUNTERS = "rulecounters"
)

const (
	COMMAND    = "command"
	DELAY      = "delay"
	DROP       = "drop"
	RETEMPTY   = "return_empty"
	RETERR     = "return_err"
	CLIENTADDR = "client_addr"
	PERCENTAGE = "percentage"
)

type Controller struct {
	plan *Plan
}

func newController(p *Plan) (*Controller, error) {
	return &Controller{
		plan: p,
	}, nil
}

func (c *Controller) parseRule(rule *Rule, buf string) error {
	kv := strings.Split(buf, "=")
	if len(kv) < 2 {
		return fmt.Errorf("rule arguments are too low")
	}

	key := strings.ToLower(kv[0])
	switch key {
	case DELAY:
		delay, err := strconv.Atoi(kv[1])
		if err != nil {
			return fmt.Errorf("parse delay error: %s", err.Error())
		}
		rule.Delay = delay
	case DROP:
		rule.Drop = false
		if strings.ToLower(kv[1]) == "true" || kv[1] == "1" {
			rule.Drop = true
		}
	case RETEMPTY:
		rule.ReturnEmpty = false
		if strings.ToLower(kv[1]) == "true" || kv[1] == "1" {
			rule.ReturnEmpty = true
		}
	case CLIENTADDR:
		if len(kv[1]) <= 0 {
			return fmt.Errorf("client addr mustn't be empty")
		}
		rule.ClientAddr = kv[1]
	case RETERR:
		if len(kv[1]) <= 0 {
			return fmt.Errorf("error message mustn't be empty")
		}
		rule.ReturnErr = kv[1]
	case PERCENTAGE:
		perc, err := strconv.Atoi(kv[1])
		if err != nil {
			return fmt.Errorf("parse delay error: %s", err.Error())
		}
		rule.Percentage = perc
	}

	return nil
}

func (c *Controller) Start() error {
	err := redcon.ListenAndServe(configAddr,
		func(conn redcon.Conn, cmd redcon.Command) {
			switch strings.ToLower(string(cmd.Args[0])) {
			default:
				conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
			case RULEADD:
				rulename := string(cmd.Args[1])
				if len(rulename) <= 0 {
					conn.WriteError("name mustn't be empty")
					return
				}

				rule := Rule{
					Name: rulename,
				}
				for _, frag := range cmd.Args[2:] {
					err := c.parseRule(&rule, string(frag))
					if err != nil {
						conn.WriteError(err.Error())
						return
					}
				}

				err := c.plan.AddRule(rule)
				if err != nil {
					conn.WriteError(err.Error())
					return
				}

				conn.WriteString("OK")
			case RULEDEL:
				rulename := string(cmd.Args[1])
				if len(rulename) <= 0 {
					conn.WriteError("name mustn't be empty")
					return
				}

				err := c.plan.DeleteRule(rulename)
				if err != nil {
					conn.WriteError(err.Error())
					return
				}
				conn.WriteString("OK")
			case RULELIST:
				rules := c.plan.ListRules()
				conn.WriteArray(len(rules))
				for _, rule := range rules {
					conn.WriteString(rule.String())
				}
			}

		},
		func(conn redcon.Conn) bool {
			fmt.Println(conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)
	fmt.Println("Closed", err)
	return err
}
