package lvsnetwork_test

import (
	"testing"

	"github.com/jeremmfr/terraform-provider-lvsnetwork/lvsnetwork"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := lvsnetwork.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = lvsnetwork.Provider()
}
