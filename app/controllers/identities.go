package controllers

import (
	"github.com/revel/revel"

	"identities/app/api"
)

type Identities struct {
	App
}

func (c *Identities) ID(network, username string) revel.Result {
	if res, err := api.GetIdentity(network, "", username); err == nil {
		return c.Data(res)
	} else {
		return c.Error(err)
	}
}

func (c *Identities) Username(network, id string) revel.Result {
	if res, err := api.GetIdentity(network, id, ""); err == nil {
		return c.Data(res)
	} else {
		return c.Error(err)
	}
}

func (c *Identities) Identity(user string) revel.Result {
	if res, err := api.Identify(user); err == nil {
		return c.Data(res)
	} else {
		return c.Error(err)
	}
}
