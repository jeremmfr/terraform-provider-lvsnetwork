package main

import (
	"github.com/jeremmfr/terraform-provider-lvsnetwork/lvsnetwork"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: lvsnetwork.Provider,
	})
}
