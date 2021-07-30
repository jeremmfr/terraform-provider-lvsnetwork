# terraform-provider-lvsnetwork

![GitHub release (latest by date)](https://img.shields.io/github/v/release/jeremmfr/terraform-provider-lvsnetwork)
[![Registry](https://img.shields.io/badge/registry-doc%40latest-lightgrey?logo=terraform)](https://registry.terraform.io/providers/jeremmfr/lvsnetwork/latest/docs)
[![Go Status](https://github.com/jeremmfr/terraform-provider-lvsnetwork/workflows/Go%20Tests/badge.svg)](https://github.com/jeremmfr/terraform-provider-lvsnetwork/actions)
[![Lint Status](https://github.com/jeremmfr/terraform-provider-lvsnetwork/workflows/GolangCI-Lint/badge.svg)](https://github.com/jeremmfr/terraform-provider-lvsnetwork/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeremmfr/terraform-provider-lvsnetwork)](https://goreportcard.com/report/github.com/jeremmfr/terraform-provider-lvsnetwork)

Terraform's provider for setup network interface and keepalived vrrp_instance on two server (master/slave) with [lvsnetwork-api](https://github.com/jeremmfr/lvsnetwork-api)

## Automatic install (Terraform 0.13 and later)

Add source information inside the Terraform configuration block for automatic provider installation:

```hcl
terraform {
  required_providers {
    lvsnetwork = {
      source = "jeremmfr/lvsnetwork"
    }
  }
}
```

## Documentation

[registry.terraform.io](https://registry.terraform.io/providers/jeremmfr/lvsnetwork/latest/docs)

or in docs :

[terraform-provider-lvsnetwork](docs/index.md)  

Resources:

* [lvsnetwork_ifacevrrp](docs/resources/ifacevrrp.md)
* [lvsnetwork_vrrp_script](docs/resources/vrrp_script.md)

## Compile

```shell
go build
```
