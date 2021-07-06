package lvsnetwork_test

import (
	"testing"

	"github.com/jeremmfr/terraform-provider-lvsnetwork/lvsnetwork"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestProvider(t *testing.T) {
	if err := lvsnetwork.Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = lvsnetwork.Provider()
}
