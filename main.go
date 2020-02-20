package main

import (
	"terraform-provider-lvsnetwork/lvsnetwork"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: lvsnetwork.Provider,
	})
}
